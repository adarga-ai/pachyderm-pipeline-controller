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
