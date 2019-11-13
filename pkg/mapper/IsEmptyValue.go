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

// IsEmptyValue returns true if the value is "empty".
//  - If the type is nil, then returns true.
//  - If the kind is bool, and the value is false, then returns true.
//  - If the value is an integer, unsigned integer, or float, and the value if zero, then returns true.
//  - If the value is of kind array, map, slice, or string, and if length is zero, then returns true.
//  - If the value is of kind chan, func, interface, or pointer, and is not valid or nil, then returns true.
//  - In all other cases, returns false.
// IsEmptyValue is adapted from the isEmptyValue function in the `encoding/json` package, but adds support for additional types.
func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Ptr:
		return (!v.IsValid()) || v.IsNil()
	}
	return false
}
