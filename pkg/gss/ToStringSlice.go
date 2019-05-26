// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"fmt"
)

// ToStringSlice converts a slice of interface{} to a slice of strings using fmt.Sprint.
func ToStringSlice(keys []interface{}) []string {
	stringSlice := make([]string, 0, len(keys))
	for _, v := range keys {
		stringSlice = append(stringSlice, fmt.Sprint(v))
	}
	return stringSlice
}
