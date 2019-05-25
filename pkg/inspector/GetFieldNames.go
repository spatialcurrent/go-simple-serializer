// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package inspector

import (
	"reflect"
)

// GetKeys returns the field names of a struct as []string.
// If you want the field names to be sorted in alphabetical order, pass sorted equal to true.
func GetFieldNames(object interface{}, sorted bool) []string {
	return GetFieldNamesFromValue(reflect.ValueOf(object), sorted)
}
