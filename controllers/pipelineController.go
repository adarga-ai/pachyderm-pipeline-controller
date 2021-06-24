/* Copyright 2021 Adarga Limited
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License"). 
 * You may not use this file except in compliance with the License. 
 * You may obtain a copy of the License at:
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
 * implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package controllers

import (
	"context"
	"time"

	"github.com/pachyderm/pachyderm/src/client/pps"
	"github.com/spotahome/kooper/v2/controller"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"go.uber.org/zap"

	"github.com/adarga-ai/pachyderm-pipeline-controller/apis/pachyderm/v1alpha1"
	adargaclients "github.com/adarga-ai/pachyderm-pipeline-controller/client/clientset/versioned"
	"github.com/adarga-ai/pachyderm-pipeline-controller/config"
	"github.com/adarga-ai/pachyderm-pipeline-controller/utils"
	pachydermclient "github.com/pachyderm/pachyderm/src/client"
)

func GeneratePipelineController(k8sConfig *rest.Config, namespace string, pachydermAddress string) (controller.Controller, error) {
	logger, klogger := config.GetLogger()

	pachyClient, err := pachydermclient.NewFromAddress(pachydermAddress)
	if err != nil {
		logger.Error("Connecting to PachD failed. Controller cannot function.", zap.String("Error", err.Error()))
		return nil, err
	}

	adargaclient, err := adargaclients.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}

	retr := controller.MustRetrieverFromListerWatcher(&cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return adargaclient.AdargaV1alpha1().PachydermPipelines(namespace).List(options)
		},

		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.Watch = true
			return adargaclient.AdargaV1alpha1().PachydermPipelines(namespace).Watch(options)
		},
	})

	hand := controller.HandlerFunc(func(_ context.Context, obj runtime.Object) error {
		pipeline := obj.(*v1alpha1.PachydermPipeline)

		plog := logger.With(zap.String("PachydermPipeline", pipeline.Name))

		plog.Info("Pipeline Object Recieved.")

		const finalizer = "adarga.ai/pachyderm-pipeline-finalizer"
		hasDeleteTime := !pipeline.DeletionTimestamp.IsZero()
		hasFinalizer := utils.StringPresentInSlice(pipeline.Finalizers, finalizer)

		switch {

		case hasDeleteTime && hasFinalizer:
			// Handle the Deletion process
			plog.Info("Object has deletion timestamp and finalizer. Running deletion.")
			ctx := context.Background()
			deleteRequest := &pps.DeletePipelineRequest{
				Pipeline: pachydermclient.NewPipeline(pipeline.Spec.Name),
				Force:    true,
				KeepRepo: true,
			}

			_, err := pachyClient.PpsAPIClient.DeletePipeline(ctx, deleteRequest)
			if err != nil {

				plog.Error("Could not delete Pipeline. Exiting Handler.", zap.String("Error", err.Error()))

				return err
			}

			// remove finalizer
			pipeline.Finalizers = utils.RemoveStringFromSlice(pipeline.Finalizers, finalizer)

			if _, err = adargaclient.AdargaV1alpha1().PachydermPipelines(namespace).Update(pipeline); err != nil {
				plog.Error("Error trying to update pipeline object after pachyderm deletion.", zap.String("Error", err.Error()))
			}
			plog.Info("End of deletion of PachydermPipeline object.")

		case !hasDeleteTime && !hasFinalizer:

			// Newly created
			plog.Info("Pipeline due for creation.")

			if len(pipeline.Status.Conditions) == 0 {
				plog.Info("Status not found, attaching Status struct.")
				pipeline.Status = utils.CreateStatus()
			}

			if pipeline, err = adargaclient.AdargaV1alpha1().PachydermPipelines(namespace).UpdateStatus(pipeline); err != nil {
				plog.Error("Error trying to update pipeline status.", zap.String("Error", err.Error()))

				return err
			}

			ctx := context.Background()
			pipelinerequest := utils.CreateRequest(*pipeline)
			_, err := pachyClient.PpsAPIClient.CreatePipeline(
				ctx,
				pipelinerequest,
			)
			if err != nil {

				utils.AmendCreationStatus(&pipeline.Status, v1alpha1.ConditionCreationError, "CreationError", "Error Applying Config. ")

				plog.Error("Creation of pipeline failed.", zap.String("Error", err.Error()))

			} else {
				utils.AmendCreationStatus(&pipeline.Status, v1alpha1.ConditionCreated, "Succeeded", "Succeeded.")
				utils.AmendRunningStatus(&pipeline.Status, v1alpha1.ConditionRunning, "Running", "Running.")

			}

			plog.Info("Now updating the pipeline object.")
			if pipeline, err = adargaclient.AdargaV1alpha1().PachydermPipelines(namespace).UpdateStatus(pipeline); err != nil {
				plog.Error("Error trying to update pipeline status after creation.", zap.String("Error", err.Error()))
			}
			// NOTE: appending finalizer behaviour is strange after using updateStatus -> has to be in between?
			pipeline.Finalizers = append(pipeline.Finalizers, finalizer)

			if pipeline, err = adargaclient.AdargaV1alpha1().PachydermPipelines(namespace).Update(pipeline); err != nil {
				plog.Error("Error trying to update pipeline finalizer after creation.", zap.String("Error", err.Error()))
			}

			plog.Info("Successful Deployment of Pipeline.")

		case hasDeleteTime && !hasFinalizer:
			// Deletion has already been handled.
			plog.Info("Pipeline has no finalizer but delete time. Pipeline deletion has been handled. No more action.")

		case !hasDeleteTime && hasFinalizer:
			plog.Info("No delete time, but finalizer. Pipeline Update in progress")

			ctx := context.Background()

			// Get all current pipelines.
			listInfo, err := pachyClient.ListPipeline()
			if err != nil {
				plog.Error("Client ListPipeline error!", zap.String("error", err.Error()))
				plog.Debug("Breaking rest of update.")
				break
			}

			// Create index
			mapping := utils.PipelineInfoListToMap(listInfo)

			if existingPipeline, exists := mapping[pipeline.Spec.Name]; exists {
				plog.Info("Existing pipeline found, checking for equality.")

				// otherwise,
				if existingPipeline.Stopped {
					utils.AmendRunningStatus(&pipeline.Status, v1alpha1.ConditionStopped, "", "")
					if _, err = adargaclient.AdargaV1alpha1().PachydermPipelines(namespace).UpdateStatus(pipeline); err != nil {
						plog.Error("Error trying to update pipeline about failing status.", zap.String("Error", err.Error()))
					}

				} else if existingPipeline.State == pps.PipelineState_PIPELINE_RUNNING {
					utils.AmendRunningStatus(&pipeline.Status, v1alpha1.ConditionRunning, "", "")
					if _, err = adargaclient.AdargaV1alpha1().PachydermPipelines(namespace).UpdateStatus(pipeline); err != nil {
						plog.Error("Error trying to update pipeline about failing status.", zap.String("Error", err.Error()))
					}
				}

				// Check equality of found pipeline.
				if !utils.CheckPipelineEquality(pipeline, existingPipeline, plog) {

					// Need to update
					plog.Info("Differences found, updating.")
					pipelinerequest := utils.CreateRequest(*pipeline)
					pipelinerequest.Update = true

					_, err = pachyClient.PpsAPIClient.CreatePipeline(
						ctx,
						pipelinerequest,
					)
					if err != nil {
						plog.Error("Update of pipeline failed", zap.String("Error", err.Error()))
					} else {
						plog.Info("Successful update of pipeline.")
					}
				} else {
					plog.Info("Objects are equal. No update.")
				}
			} else {
				// Should never occur
				utils.AmendRunningStatus(&pipeline.Status, v1alpha1.ConditionMissing, "PipelineNotExist", "ManualManipulation.")
				if _, err = adargaclient.AdargaV1alpha1().PachydermPipelines(namespace).UpdateStatus(pipeline); err != nil {
					plog.Error("Error trying to update pipeline about failing status.", zap.String("Error", err.Error()))
				}

				plog.Warn("Existing pipeline not found, but new object has finalizer. Indicates error in creation.")

			}

		}

		plog.Debug("End of Handler loop.")

		return nil
	})

	cfg := &controller.Config{
		Name:      "qd-pipeline-controller",
		Handler:   hand,
		Retriever: retr,
		Logger:    klogger,

		ProcessingJobRetries: 5,
		ResyncInterval:       15 * time.Second,
		ConcurrentWorkers:    1,
	}
	pipelineController, err := controller.New(cfg)
	if err != nil {
		return nil, err
	}

	return pipelineController, nil
}
