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

// ExpandHeaderWithWildcard expands a table header at the wildcard position with keys from the object.
func ExpandHeaderWithWildcard(header []interface{}, knownKeys map[interface{}]struct{}, object reflect.Value, sorted bool, reversed bool) ([]interface{}, map[interface{}]struct{}) {
	newHeader := make([]interface{}, 0, len(header))
	newKnownKeys := map[interface{}]struct{}{}
	for _, knownKey := range header {
		if str, ok := knownKey.(string); ok && str == Wildcard {
			for _, unknownKey := range inspector.GetUnknownKeysFromValue(object, knownKeys, sorted, reversed) {
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
