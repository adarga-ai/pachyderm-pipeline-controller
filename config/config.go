// Copyright 2021 Adarga Limited
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	PachydermAddress string
	Namespace        string
}

func GenerateConfig() *Config {
	viper.AutomaticEnv()

	viper.SetDefault("PACHYDERM_ADDRESS", "localhost:30650")
	viper.SetDefault("NAMESPACE", "pach-feed-1-12")

	config := &Config{
		PachydermAddress: viper.GetString("PACHYDERM_ADDRESS"),
		Namespace:        viper.GetString("NAMESPACE"),
	}

	return config
}
