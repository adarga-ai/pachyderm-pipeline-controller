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

package utils

import (
	"github.com/adarga-ai/pachyderm-pipeline-controller/apis/pachyderm/v1alpha1"
	"github.com/pachyderm/pachyderm/src/client/pps"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var existingPipeline pps.PipelineInfo = pps.PipelineInfo{
	ID:       "test",
	Pipeline: &pps.Pipeline{Name: "test"},
	Version:  0,
	Transform: &pps.Transform{
		Image: "test",
		Cmd:   []string{"test"},
		Env:   map[string]string{},
	},
	Egress: &pps.Egress{URL: "test"},

	ResourceLimits: &pps.ResourceSpec{

		Memory: "test",
		Gpu: &pps.GPUSpec{
			Type:   "test",
			Number: 100,
		},
	},
	Input: &pps.Input{

		Join: []*pps.Input{
			{
				Pfs: &pps.PFSInput{
					Name:   "test",
					Repo:   "test",
					Branch: "test",
					Glob:   "test",
					JoinOn: "$1",
				},
			},
			{
				Pfs: &pps.PFSInput{
					Name:   "test",
					Repo:   "test",
					Branch: "test",
					Glob:   "test",
					JoinOn: "$1",
				},
			},
		},
	},
	Description: "test",
	EnableStats: false,
}
var testStringPointer = "test"

var newPipeline v1alpha1.PachydermPipeline = v1alpha1.PachydermPipeline{
	TypeMeta:   v1.TypeMeta{},
	ObjectMeta: v1.ObjectMeta{},
	Spec: v1alpha1.PachydermPipelineSpec{
		Name: "test",
		Transform: v1alpha1.Transform{
			CMD:   []string{"test"},
			Image: "test",
			Env:   []corev1.EnvVar{},
		},
		Input: v1alpha1.Input{

			Join: []*v1alpha1.PFS{
				{
					Repo:   "test",
					Glob:   "test",
					Branch: "test",
				},
				{
					Repo:   "test",
					Glob:   "test",
					Branch: "test",
				},
			},
		},
		EnableStats: false,
		Egress:      &testStringPointer,
		ResourceLimits: &v1alpha1.ResourceLimits{
			Memory: &testStringPointer,
			GPU: &v1alpha1.GPU{
				GPUType: "test",
				Number:  100,
			},
		},
		Description: "test",
	},
	Status: v1alpha1.PachydermPipelineStatus{},
}
