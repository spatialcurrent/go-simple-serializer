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

// RowToMapOfStrings takes in a header and corresponding row and returns a map from the values.
func RowToMapOfStrings(header []string, row []string) map[string]string {
	m := map[string]string{}
	for i, h := range header {
		m[strings.ToLower(h)] = row[i]
	}
	return m
}
