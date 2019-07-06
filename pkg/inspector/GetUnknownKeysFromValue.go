// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package inspector

import (
	"fmt"
	"reflect"
	"sort"
)

// GetUnknownKeysFromValue returns the unknown keys for a map as an []interface{} given a set of known keys.
// If you want the keys to be sorted in alphabetical order, pass sorted equal to true.
// If sorted and reversed, then sorts in reverse alphabetical order.
func GetUnknownKeysFromValue(m reflect.Value, knownKeys map[interface{}]struct{}, sorted bool, reversed bool) []interface{} {
	unknownKeys := make([]interface{}, 0)
	for _, key := range m.MapKeys() {
		keyInterface := key.Interface()
		if _, exists := knownKeys[keyInterface]; !exists {
			unknownKeys = append(unknownKeys, keyInterface)
		}
	}
	if sorted {
		sort.Slice(unknownKeys, func(i, j int) bool {
			if reversed {
				return fmt.Sprint(unknownKeys[i]) > fmt.Sprint(unknownKeys[j])
			}
			return fmt.Sprint(unknownKeys[i]) < fmt.Sprint(unknownKeys[j])
		})
	}
	return unknownKeys
}
