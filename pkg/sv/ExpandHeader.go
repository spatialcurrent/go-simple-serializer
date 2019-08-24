// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"reflect"

	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
)

// ExpandHeader expands a table header with the keys from the given object.
func ExpandHeader(header []interface{}, knownKeys map[interface{}]struct{}, object reflect.Value, sorted bool, reversed bool) ([]interface{}, map[interface{}]struct{}) {
	newHeader := make([]interface{}, 0, len(header))
	newKnownKeys := map[interface{}]struct{}{}
	for _, knownKey := range header {
		newHeader = append(newHeader, knownKey)
		newKnownKeys[knownKey] = struct{}{}
	}
	for _, unknownKey := range inspector.GetUnknownKeysFromValue(object, knownKeys, sorted, reversed) {
		newHeader = append(newHeader, unknownKey)
		newKnownKeys[unknownKey] = struct{}{}
	}
	return newHeader, newKnownKeys
}
