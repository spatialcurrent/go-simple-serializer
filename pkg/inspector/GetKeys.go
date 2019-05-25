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

// GetKeys returns the keys for a map as an []interface{}.
// If you want the keys to be sorted in alphabetical order, pass sorted equal to true.
// If sorted and reversed, then sorts in reverse alphabetical order.
func GetKeys(object interface{}, sorted bool, reversed bool) []interface{} {
	return GetKeysFromValue(reflect.ValueOf(object), sorted, reversed)
}
