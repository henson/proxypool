//
// Copyright (C) 2017-2018 The Aiicy Team. - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by CountStarlight <countstarlight@gmail.com>, June 2018
//
package util

import "strings"

// IsSliceContainsStr returns true if the string exists in given slice, ignore case.
func IsSliceContainsStr(sl []string, str string) bool {
	str = strings.ToLower(str)
	for _, s := range sl {
		if strings.ToLower(s) == str {
			return true
		}
	}
	return false
}
