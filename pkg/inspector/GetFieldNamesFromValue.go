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

// GetFieldNamesFromValue returns the field names of a struct as []string.
// If you want the field names to be sorted in alphabetical order, pass sorted equal to true.
func GetFieldNamesFromValue(value reflect.Value, sorted bool) []string {
	fieldNames := make([]string, 0, value.NumField())
	t := value.Type()
	for i := 0; i < value.NumField(); i++ {
		fieldNames = append(fieldNames, t.Field(i).Name)
	}
	if sorted {
		sort.Strings(fieldNames)
	}
	return fieldNames
}
