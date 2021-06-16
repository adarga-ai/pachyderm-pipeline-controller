// Copyright 2021 Adarga Limited
// SPDX-License-Identifier: Apache-2.0

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
