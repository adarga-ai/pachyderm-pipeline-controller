// Copyright 2021 Adarga Limited
// SPDX-License-Identifier: Apache-2.0

package utils

/*
   Testing checklist:
   add the comparison checker.
   creating request.

*/

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestComparisonReturnsEqual(t *testing.T) {
	logger := zaptest.NewLogger(t)

	if !CheckPipelineEquality(&newPipeline, &existingPipeline, logger) {
		t.Fail()
		t.Logf("The two objects are not equal.")
	}
}

func TestComparisonReturnsFalse(t *testing.T) {
	logger := zaptest.NewLogger(t)

	existingPipelineCopy := existingPipeline
	newPipelineCopy := newPipeline

	existingPipelineCopy.Transform.Image = "Should yeild False"
	if CheckPipelineEquality(&newPipelineCopy, &existingPipelineCopy, logger) {

		t.Fail()
		t.Logf("The two objects are equal when they should not be.")
	}
}
