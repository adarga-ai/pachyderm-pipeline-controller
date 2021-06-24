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
