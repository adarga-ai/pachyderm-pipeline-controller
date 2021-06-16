// Copyright 2021 Adarga Limited
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	adargaclients "github.com/adarga-ai/pachyderm-pipeline-controller/client/clientset/versioned"
	"github.com/adarga-ai/pachyderm-pipeline-controller/config"
	"github.com/spotahome/kooper/v2/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

// matches a PachydermPipeline pipeline object to a replicationController.
// returns true if the controller has an owner. If it doesnt, it will find one and join, then return true.
// if none is found.. most probably a manually deployed pipe, so returns false but no error raised.
func MergeController(adargaClient *adargaclients.Clientset, replicationController *corev1.ReplicationController, namespace string, k8scli kubernetes.Clientset, logger *zap.Logger) (bool, error) {
	logger.Info("Attempting to merge.")

	if len(replicationController.OwnerReferences) > 0 {
		// already linked
		logger.Info("Replication Controller already has owner. Will skip.")
		return true, nil
	}

	name := replicationController.Labels["pipelineName"]
	if name == "" {

		// replicationController not part of Pachyderm process
		logger.Info("ReplicatioController does not have Pipeline Label. Ignoring.")
		return false, nil
	}

	// get the pachyderm pipeline, it should exist.
	pipe, err := adargaClient.AdargaV1alpha1().PachydermPipelines(namespace).Get(
		name,
		v1.GetOptions{ResourceVersion: "0"})
	if err != nil {
		// there are two different types of errors here, one for not found( which we need to ignore)
		// and then anything else

		if strings.HasSuffix(err.Error(), "found") {
			// assume that there is no pipe for this, manually deployed.
			logger.Warn("Pipeline Object is not found. Potential Error! Potentially manually deployed pipeline.")

			return false, nil
		} else {
			return false, err
		}
	}

	ref := metav1.OwnerReference{
		APIVersion:         "v1",
		Kind:               "PachydermPipeline",
		Name:               pipe.Spec.Name,
		UID:                pipe.UID,
		Controller:         new(bool),
		BlockOwnerDeletion: new(bool),
	}

	replicationController.OwnerReferences = append(replicationController.OwnerReferences, ref)

	_, err = k8scli.CoreV1().ReplicationControllers(namespace).Update(replicationController)
	if err != nil {
		logger.Error("Failed to update through API", zap.String("Error", err.Error()))
		return false, err
	}

	return true, nil
}

func GenerateReplicationController(k8sConfig *rest.Config, adargaclient *adargaclients.Clientset, namespace string) (controller.Controller, error) {
	logger, klogger := config.GetLogger()
	k8scli, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes client: %w", err)
	}

	// business logic for the attaching Pipelines to Replication Controllers.
	replicationHandler := controller.HandlerFunc(func(_ context.Context, obj runtime.Object) error {
		repController, ok := obj.(*corev1.ReplicationController)

		if ok {
			rlogger := logger.With(zap.String("name", repController.Name))

			rlogger.Info("Found Replication Controller")
			merged, err := MergeController(adargaclient, repController, namespace, *k8scli, rlogger)
			if err != nil {
				rlogger.Error("Failed. End of Handler.")
				return err
			}

			if !merged {
				rlogger.Warn("No merging occurred but no error. Potential Bug")
			} else {
				rlogger.Info("Merging process complete.")
			}

			return nil
		}

		return nil
	})

	repController, err := controller.New(

		&controller.Config{
			Handler: replicationHandler,
			Retriever: controller.MustRetrieverFromListerWatcher(&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return k8scli.CoreV1().ReplicationControllers(namespace).List(options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return k8scli.CoreV1().ReplicationControllers(namespace).Watch(options)
				},
			}),

			Logger:               klogger,
			Name:                 "ReplicationController Merger",
			ConcurrentWorkers:    1,
			ResyncInterval:       10 * time.Second,
			ProcessingJobRetries: 5,
			DisableResync:        false,
		})

	return repController, nil
}
