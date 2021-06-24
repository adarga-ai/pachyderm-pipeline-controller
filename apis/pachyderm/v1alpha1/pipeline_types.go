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

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	GroupName string = "adarga.ai"
	Version   string = "valpha1"
	Plural    string = "pachydermpipelines"
	Singluar  string = "pachydermpipeline"
	ShortName string = "pipe"
	Name      string = Plural + "." + GroupName
)

// +crd
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:singular=pachydermpipeline,path=pachydermpipelines,shortName=pipe,scope=Namespaced
// +groupName=adarga.ai
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name=Created,type="string",JSONPath=`.status.conditions[?(@.type=="Creation")].status`
// +kubebuilder:printcolumn:name=Condition,type="string",JSONPath=`.status.conditions[?(@.type=="Running")].status`
type PachydermPipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PachydermPipelineSpec `json:"spec"`

	// +kubebuilder:validation:Optional
	Status PachydermPipelineStatus `json:"status,omitempty"`
}

type PachydermPipelineStatus struct {
	// Conditions represents the latest available observations of a replication controllers current state.
	// +kubebuilder:validation:Optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []PachydermPipelineCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// The valid conditions of a pachyderm pipeline failure
type PachydermPipelineConditionType string

const (
	// When the pipeline object is not being instantiated properly (bad repo)
	PachydermPipelineCreationCondition PachydermPipelineConditionType = "Creation"
	// When the pipeline object has a finaliser, but no running pipeline in pachd.
	PachydermPipelineRunningCondition PachydermPipelineConditionType = "Running"
)

type PipelineStatus string

// potentially breaking API guidelines by going for this format
const (
	ConditionRunning       PipelineStatus = "Running"
	ConditionMissing       PipelineStatus = "Missing"
	ConditionStopped       PipelineStatus = "Stopped"
	ConditionCreated       PipelineStatus = "Success"
	ConditionCreationError PipelineStatus = "CreationError"
	ConditionUnknown       PipelineStatus = "Unknown"
)

// Describes the state of a pipeline condition at a certain point in time.
type PachydermPipelineCondition struct {
	Type PachydermPipelineConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ReplicationControllerConditionType"`
	// Status of the condition, one of True, False, Unknown.
	Status PipelineStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`

	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`

	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

type PachydermPipelineSpec struct {

	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	Transform Transform `json:"transform"`
	// +kubebuilder:validation:Required
	Input Input `json:"input"`
	// +kubebuilder:validation:Required
	EnableStats bool `json:"enable_stats"`

	// +kubebuilder:validation:Optional
	Egress *string `json:"egress,omitempty"`
	// +kubebuilder:validation:Optional
	ResourceLimits *ResourceLimits `json:"resource_limits,omitempty"`
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`
}

type Transform struct {
	// +kubebuilder:validation:Required
	CMD []string `json:"cmd"`
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// +kubebuilder:validation:Optional
	Env []v1.EnvVar `json:"env,omitempty"`
}

type Input struct {
	// +kubebuilder:validation:Optional
	PFS *PFS `json:"pfs,omitempty"`
	// +kubebuilder:validation:Optional
	Join []*PFS `json:"join,omitempty"`
	// +kubebuilder:validation:Optional
	Cross []*PFS `json:"cross,omitempty"`
}

type PFS struct {
	// +kubebuilder:validation:Required
	Repo string `json:"repo"`
	// +kubebuilder:validation:Required
	Glob string `json:"glob"`
	// +kubebuilder:validation:Required
	Branch string `json:"branch"`
}

type Egress struct {
	// +kubebuilder:validation:Required
	URL string `json:"url"`
}

// add CPU?
type ResourceLimits struct {
	// +kubebuilder:validation:Optional
	Memory *string `json:"memory,omitempty"`
	// +kubebuilder:validation:Optional
	GPU *GPU `json:"gpu,omitempty"`
}

type GPU struct {
	GPUType string `json:"type"`
	Number  int64  `json:"number"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type PachydermPipelineList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []PachydermPipeline
}
