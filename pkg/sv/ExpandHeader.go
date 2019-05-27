// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"reflect"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
)

// ExpandHeader expands a table header given keys found in a new object
func ExpandHeader(header []interface{}, knownKeys map[interface{}]struct{}, object reflect.Value, sorted bool) ([]interface{}, map[interface{}]struct{}) {
	newHeader := make([]interface{}, 0, len(header))
	newKnownKeys := map[interface{}]struct{}{}
	for _, knownKey := range header {
		if str, ok := knownKey.(string); ok && str == Wildcard {
			for _, unknownKey := range inspector.GetUnknownKeysFromValue(object, knownKeys, sorted) {
				newHeader = append(newHeader, unknownKey)
				newKnownKeys[unknownKey] = struct{}{}
			}
			// Keep the wildcard in place for future expanding of the heade
			newHeader = append(newHeader, knownKey)
		} else {
			newHeader = append(newHeader, knownKey)
			newKnownKeys[knownKey] = struct{}{}
		}
	}
	return newHeader, newKnownKeys
}
