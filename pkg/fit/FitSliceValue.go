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

// FitSliceValue iterates through and fits the given slice value and its underlying values.
func FitSliceValue(in reflect.Value) reflect.Value {

	if in.Len() == 0 {
		return in
	}

	types := map[string]reflect.Type{}

	for i := 0; i < in.Len(); i++ {
		t := reflect.ValueOf(in.Index(i).Interface()).Type()
		types[t.Name()] = t
	}

	if len(types) == 1 {
		var t reflect.Type
		for _, v := range types {
			t = v
		}
		out := reflect.MakeSlice(reflect.SliceOf(t), 0, in.Len())
		for i := 0; i < in.Len(); i++ {
			out = reflect.Append(out, FitValue(reflect.ValueOf(in.Index(i).Interface())))
		}
		return out
	}

	out := reflect.MakeSlice(in.Type(), 0, in.Len())
	for i := 0; i < in.Len(); i++ {
		out = reflect.Append(out, FitValue(reflect.ValueOf(in.Index(i).Interface())))
	}

	return out

}
