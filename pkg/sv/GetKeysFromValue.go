// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"fmt"
	"reflect"
	"sort"
)

// GetKeysFromValue returns the keys for a map as an []interface{}.
// If you want the keys to be sorted in alphabetical order, pass sorted equal to true.
func GetKeysFromValue(m reflect.Value, sorted bool) []interface{} {
	keys := make([]interface{}, m.Len())
	for i, key := range m.MapKeys() {
		keys[i] = key.Interface()
	}
	if sorted {
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j])
		})
	}
	return keys
}
