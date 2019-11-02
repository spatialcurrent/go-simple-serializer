// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package mapper

import (
	"reflect"
)

// IsEmpty returns true if the value is "empty".
// If the type is nil, then returns true.
// If the value is of kind array, map, slice, or string and if length is zero, then returns true.
// If the value is of kind chan, func, interface, or pointer, and is not valid or nil, then returns true.
// In all other cases, returns false.
func IsEmpty(v interface{}) bool {
	if t := reflect.TypeOf(v); t != nil {
		switch t.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Ptr:
			return (!reflect.ValueOf(v).IsValid()) || reflect.ValueOf(v).IsNil()
		case reflect.Array, reflect.String, reflect.Map, reflect.Slice:
			return reflect.ValueOf(v).Len() == 0
		default:
			return false
		}
	}
	// if the type of the value is nil
	return true
}
