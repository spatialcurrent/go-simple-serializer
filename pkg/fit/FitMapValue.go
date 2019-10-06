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

// FitMapValue iterates through and fits the given map value and its underlying values.
func FitMapValue(in reflect.Value) reflect.Value {

	if in.Len() == 0 {
		return in
	}

	out := reflect.MakeMap(in.Type())

	it := in.MapRange()

	for it.Next() {
		out.SetMapIndex(it.Key(), FitValue(reflect.ValueOf(it.Value().Interface())))
	}

	return out

}
