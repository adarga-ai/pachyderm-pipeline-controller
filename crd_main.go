// Copyright 2021 Adarga Limited
// SPDX-License-Identifier: Apache-2.0

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
