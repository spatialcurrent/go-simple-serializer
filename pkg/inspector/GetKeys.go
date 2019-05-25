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
func GetKeys(object interface{}, sorted bool) []interface{} {
	return GetKeysFromValue(reflect.ValueOf(object), sorted)
}
