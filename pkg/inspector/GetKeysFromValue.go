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

// GetKeysFromValue returns the keys for a map as an []interface{}.
// If you want the keys to be sorted in alphabetical order, pass sorted equal to true.
// If sorted and reversed, then sorts in reverse alphabetical order.
func GetKeysFromValue(m reflect.Value, sorted bool, reversed bool) []interface{} {
	keys := make([]interface{}, 0, m.Len())
	for _, key := range m.MapKeys() {
		keys = append(keys, key.Interface())
	}
	if sorted {
		sort.Slice(keys, func(i, j int) bool {
			if reversed {
				return fmt.Sprint(keys[i]) > fmt.Sprint(keys[j])
			}
			return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j])
		})
	}
	return keys
}
