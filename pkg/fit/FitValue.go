// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package fit

import (
	"reflect"
)

// FitValue fits the value and any underlying values.
func FitValue(in reflect.Value) reflect.Value {
	switch in.Type().Kind() {
	case reflect.Array, reflect.Slice:
		return FitSliceValue(in)
	case reflect.Map:
		return FitMapValue(in)
	}
	return in
}
