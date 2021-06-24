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

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"go.uber.org/zap"

	adargaclients "github.com/adarga-ai/pachyderm-pipeline-controller/client/clientset/versioned"
	"github.com/adarga-ai/pachyderm-pipeline-controller/config"
	"github.com/adarga-ai/pachyderm-pipeline-controller/controllers"
)

func run() error {
	logger, _ := config.GetLogger()
	config := config.GenerateConfig()

	logger.Info("Starting the Controllers.")
	// Get k8s client.
	k8scfg, err := rest.InClusterConfig()
	if err != nil {
		// No in cluster? letr's try locally
		kubehome := filepath.Join(homedir.HomeDir(), ".kube", "config")
		k8scfg, err = clientcmd.BuildConfigFromFlags("", kubehome)
		if err != nil {
			return fmt.Errorf("error loading kubernetes configuration: %w", err)
		}
	}
	adargaclient, err := adargaclients.NewForConfig(k8scfg)
	if err != nil {
		return fmt.Errorf("error creating kubernetes client: %w", err)
	}

	pipelineContronller, err := controllers.GeneratePipelineController(
		k8scfg,
		config.Namespace,
		config.PachydermAddress)
	if err != nil {
		logger.DPanic("Initialization of pachyderm pipeline controller failed", zap.String("Error", err.Error()))
	}

	replicationController, err := controllers.GenerateReplicationController(
		k8scfg,
		adargaclient,
		config.Namespace)
	if err != nil {
		logger.DPanic("Initialization of Replication Controller failed", zap.String("error", err.Error()))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error)

	go func() {
		logger.Info("Starting the Replication Controller.")
		errChan <- replicationController.Run(ctx)
	}()
	go func() {
		logger.Info("Starting the Pipeline Controller.")
		errChan <- pipelineContronller.Run(ctx)
	}()

	err = <-errChan
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
