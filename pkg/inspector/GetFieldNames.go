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

// GetFieldNames returns the field names of a struct as []string.
// If you want the field names to be sorted in alphabetical order, pass sorted equal to true.
// If sorted and reversed, then sorts in reverse alphabetical order.
func GetFieldNames(object interface{}, sorted bool, reversed bool) []string {
	return GetFieldNamesFromValue(reflect.ValueOf(object), sorted, reversed)
}
