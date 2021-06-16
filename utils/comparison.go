// Copyright 2021 Adarga Limited
// SPDX-License-Identifier: Apache-2.0

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
