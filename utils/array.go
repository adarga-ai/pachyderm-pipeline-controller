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

func StringPresentInSlice(ss []string, s string) bool {
	for _, f := range ss {
		if f == s {
			return true
		}
	}
	return false
}

func RemoveStringFromSlice(ss []string, s string) []string {
	for i, f := range ss {
		if f == s {
			return append(ss[:i], ss[i+1:]...)
		}
	}
	return ss
}
