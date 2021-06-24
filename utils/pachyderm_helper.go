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
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/adarga-ai/pachyderm-pipeline-controller/apis/pachyderm/v1alpha1"
	"go.uber.org/zap"

	"github.com/kylelemons/godebug/pretty"
	"github.com/pachyderm/pachyderm/src/client/pps"
)

func CreateRequest(pipeline v1alpha1.PachydermPipeline) *pps.CreatePipelineRequest {
	pipelinerequest := &pps.CreatePipelineRequest{
		Pipeline:    &pps.Pipeline{Name: pipeline.Spec.Name},
		Description: pipeline.Spec.Description,
		Transform: &pps.Transform{
			Image: pipeline.Spec.Transform.Image,
			Cmd:   pipeline.Spec.Transform.CMD,
			Stdin: []string{},
		},
		Input:       &pps.Input{},
		Update:      false,
		EnableStats: pipeline.Spec.EnableStats,
	}

	if pipeline.Spec.Input.PFS != nil {
		pipelinerequest.Input.Pfs = &pps.PFSInput{
			Name:   pipeline.Spec.Input.PFS.Repo,
			Repo:   pipeline.Spec.Input.PFS.Repo,
			Glob:   pipeline.Spec.Input.PFS.Glob,
			Branch: pipeline.Spec.Input.PFS.Branch,
		}
	}

	// add cross and join.
	if len(pipeline.Spec.Input.Join) > 0 {

		joins := make([]*pps.Input, len(pipeline.Spec.Input.Join))
		for i, join := range pipeline.Spec.Input.Join {

			joinRequest := &pps.Input{
				Pfs: &pps.PFSInput{
					Name:   join.Repo,
					Repo:   join.Repo,
					Glob:   join.Glob,
					Branch: join.Branch,
					JoinOn: "$1",
				},
			}

			joins[i] = joinRequest
		}

		sort.Slice(joins[:], func(i, j int) bool {
			return joins[i].Pfs.Name < joins[j].Pfs.Name
		})

		pipelinerequest.Input.Join = joins
	}

	if len(pipeline.Spec.Input.Cross) > 0 {

		crosses := make([]*pps.Input, len(pipeline.Spec.Input.Cross))
		for i, cross := range pipeline.Spec.Input.Cross {

			crossRequest := &pps.Input{
				Pfs: &pps.PFSInput{
					Name:   cross.Repo,
					Repo:   cross.Repo,
					Glob:   cross.Glob,
					Branch: cross.Branch,
				},
			}

			crosses[i] = crossRequest
		}

		sort.Slice(crosses[:], func(i, j int) bool {
			return crosses[i].Pfs.Name < crosses[j].Pfs.Name
		})
		pipelinerequest.Input.Cross = crosses
	}

	if pipeline.Spec.Egress != nil {
		pipelinerequest.Egress = &pps.Egress{URL: *pipeline.Spec.Egress}
	}

	// not quite right? Memory could be different
	if pipeline.Spec.ResourceLimits != nil && pipeline.Spec.ResourceLimits.GPU != nil {
		pipelinerequest.ResourceLimits = &pps.ResourceSpec{
			Gpu: &pps.GPUSpec{
				Type:   pipeline.Spec.ResourceLimits.GPU.GPUType,
				Number: int64(pipeline.Spec.ResourceLimits.GPU.Number),
			},
		}
	}

	if pipeline.Spec.ResourceLimits != nil {
		pipelinerequest.ResourceLimits = &pps.ResourceSpec{}
		if pipeline.Spec.ResourceLimits.GPU != nil {
			pipelinerequest.ResourceLimits.Gpu = &pps.GPUSpec{
				Type:   pipeline.Spec.ResourceLimits.GPU.GPUType,
				Number: int64(pipeline.Spec.ResourceLimits.GPU.Number),
			}
		}

		if pipeline.Spec.ResourceLimits.Memory != nil {
			pipelinerequest.ResourceLimits.Memory = *pipeline.Spec.ResourceLimits.Memory
		}

	}

	// fill envs
	if len(pipeline.Spec.Transform.Env) > 0 {
		vars := make(map[string]string)
		for _, value := range pipeline.Spec.Transform.Env {
			vars[value.Name] = value.Value
		}

		pipelinerequest.Transform.Env = vars
	}

	return pipelinerequest
}

func CheckPipelineEquality(newPipeline *v1alpha1.PachydermPipeline, existingPipeline *pps.PipelineInfo, logger *zap.Logger) bool {
	existingPipelineAttributes := copyRelevantPipelineInfo(existingPipeline)
	newPipelineRequest := CreateRequest(*newPipeline)

	// serialise them.
	serialisedExisting, _ := json.Marshal(existingPipelineAttributes)
	serialisedRequest, _ := json.Marshal(newPipelineRequest)

	existingMap := make(map[string]interface{})
	requestMap := make(map[string]interface{})
	// deserialise into interface to make easier to compare, otherwise []byte

	// no need to check for errors I think.
	json.Unmarshal(serialisedExisting, &existingMap)
	json.Unmarshal(serialisedRequest, &requestMap)

	if !reflect.DeepEqual(existingMap, requestMap) {

		logger.Debug("Difference Found.", zap.String("Difference", pretty.Compare(existingMap, requestMap)))

		return false
	}

	return true
}

func PipelineInfoListToMap(list []*pps.PipelineInfo) (mapping map[string]*pps.PipelineInfo) {
	mapping = make(map[string]*pps.PipelineInfo)
	for _, pipeline := range list {
		mapping[pipeline.Pipeline.Name] = pipeline
	}

	return
}

func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
