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

import "github.com/pachyderm/pachyderm/src/client/pps"

// copies over only the relevant pipeline information that is required for comparison
func copyRelevantPipelineInfo(existingPipeline *pps.PipelineInfo) (pipelineCopy *pps.PipelineInfo) {
	pipelineCopy = &pps.PipelineInfo{
		Pipeline:       existingPipeline.Pipeline,
		Egress:         existingPipeline.Egress,
		Input:          existingPipeline.Input,
		Description:    existingPipeline.Description,
		Transform:      existingPipeline.Transform,
		EnableStats:    existingPipeline.EnableStats,
		ResourceLimits: existingPipeline.ResourceLimits,
	}

	return
}
