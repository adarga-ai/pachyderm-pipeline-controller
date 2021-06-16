// Copyright 2021 Adarga Limited
// SPDX-License-Identifier: Apache-2.0

module github.com/adarga-ai/pachyderm-pipeline-controller

go 1.16

// for compatibility with pachyderm, need to stick
// to this specific (and somewhat old) version:
replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190718183610-8e956561bbf5

require (
	github.com/Azure/go-autorest/autorest v0.11.18 // indirect
	github.com/kylelemons/godebug v1.1.0
	github.com/pachyderm/pachyderm v1.13.2
	github.com/spf13/viper v1.7.1
	github.com/spotahome/kooper/v2 v2.0.0-rc.2
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.37.0 // indirect
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v12.0.0+incompatible
)
