// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

func RemoveWildcard(header []interface{}) []interface{} {
	newHeader := make([]interface{}, 0)
	for _, x := range header {
		if str, ok := x.(string); !ok || str != Wildcard {
			newHeader = append(newHeader, x)
		}
	}
	return newHeader
}
