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
//  - If the type is nil, then returns true.
//  - If the kind is bool, and the value is false, then returns true.
//  - If the value is an integer, unsigned integer, or float, and the value if zero, then returns true.
//  - If the value is of kind array, map, slice, or string, and if length is zero, then returns true.
//  - If the value is of kind chan, func, interface, or pointer, and is not valid or nil, then returns true.
//  - In all other cases, returns false.
// IsEmptyValue is adapted from the isEmptyValue function in the `encoding/json` package, but adds support for additional types.
func IsEmpty(v interface{}) bool {
	return IsEmptyValue(reflect.ValueOf(v))
}
