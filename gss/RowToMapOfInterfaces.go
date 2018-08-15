// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"strings"
)

// RowToMapOfInterfaces takes in a header and corresponding row and returns a map from the values.
func RowToMapOfInterfaces(header []string, row []string) map[string]interface{} {
	m := map[string]interface{}{}
	for i, h := range header {
		m[strings.ToLower(h)] = row[i]
	}
	return m
}
