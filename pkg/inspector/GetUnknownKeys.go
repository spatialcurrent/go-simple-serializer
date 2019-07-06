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

// GetUnknownKeys returns the unknown keys for a map as an []interface{} given a set of known keys.
// If you want the keys to be sorted in alphabetical order, pass sorted equal to true.
// If sorted and reversed, then sorts in reverse alphabetical order.
func GetUnknownKeys(object interface{}, knownKeys map[interface{}]struct{}, sorted bool, reversed bool) []interface{} {
	return GetUnknownKeysFromValue(reflect.ValueOf(object), knownKeys, sorted, reversed)
}
