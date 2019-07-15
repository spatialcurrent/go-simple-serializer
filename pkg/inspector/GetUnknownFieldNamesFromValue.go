// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package inspector

import (
	"reflect"
	"sort"
)

// GetUnknownFieldNamesFromValue returns the unknown field names for a struct as a []string{} given a set of known field names.
// If you want the field names to be sorted in alphabetical order, pass sorted equal to true.
// If sorted and reversed, then sorts in reverse alphabetical order.
func GetUnknownFieldNamesFromValue(value reflect.Value, knownKeys map[string]struct{}, sorted bool, reversed bool) []string {
	unknownFieldNames := make([]string, 0)
	t := value.Type()
	for i := 0; i < value.NumField(); i++ {
		fieldName := t.Field(i).Name
		if _, exists := knownKeys[fieldName]; !exists {
			unknownFieldNames = append(unknownFieldNames, fieldName)
		}
	}
	if sorted {
		sort.Slice(unknownFieldNames, func(i, j int) bool {
			if reversed {
				return unknownFieldNames[i] > unknownFieldNames[j]
			}
			return unknownFieldNames[i] < unknownFieldNames[j]
		})
	}
	return unknownFieldNames
}
