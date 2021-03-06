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
)

// Returns an array with the PipelineStatus Object.
func CreateStatus() v1alpha1.PachydermPipelineStatus {
	status := &v1alpha1.PachydermPipelineStatus{
		Conditions: []v1alpha1.PachydermPipelineCondition{},
	}

	status.Conditions = append(status.Conditions, v1alpha1.PachydermPipelineCondition{
		Type:   v1alpha1.PachydermPipelineCreationCondition,
		Status: v1alpha1.ConditionUnknown,
	})

	status.Conditions = append(status.Conditions, v1alpha1.PachydermPipelineCondition{
		Type:   v1alpha1.PachydermPipelineRunningCondition,
		Status: v1alpha1.ConditionUnknown,
	})

	return *status
}

// TODO: refactor these two functions into one function.
// ammend the running status inplace.
func AmendRunningStatus(status *v1alpha1.PachydermPipelineStatus, statusCondition v1alpha1.PipelineStatus, message string, reason string) {
	// you need to use i, and cannot edit condition, otherwise it does not work, why? is it because it is not a slice of pointers?
	for i, condition := range status.Conditions {
		if condition.Type == v1alpha1.PachydermPipelineRunningCondition {
			status.Conditions[i].Message = message
			status.Conditions[i].Status = statusCondition
			status.Conditions[i].Reason = reason

		}
	}
}

// ammend the Creation status inplace.
func AmendCreationStatus(status *v1alpha1.PachydermPipelineStatus, statusCondition v1alpha1.PipelineStatus, message string, reason string) {
	for i, condition := range status.Conditions {
		if condition.Type == v1alpha1.PachydermPipelineCreationCondition {

			status.Conditions[i].Message = message
			status.Conditions[i].Status = statusCondition
			status.Conditions[i].Reason = reason

		}
	}
}
