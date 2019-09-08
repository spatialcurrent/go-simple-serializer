// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gob

import (
	"encoding/gob"
	//"fmt"
	"io"
	"reflect"
)

type Encoder struct {
	*gob.Encoder
	Fit bool
}

// Encoder returns a new gob encoder given the underlying writer.
// If fit is true, then the types of slices are fit around their underlying values.
// For example, an []interface{} slice with only float64 values is replaced with []float64 during writing.
func NewEncoder(w io.Writer, fit bool) *Encoder {
	return &Encoder{
		Encoder: gob.NewEncoder(w),
		Fit:     fit,
	}
}

func (enc *Encoder) fitValue(in reflect.Value) reflect.Value {
	switch in.Type().Kind() {
	case reflect.Array, reflect.Slice:
		return enc.fitSliceValue(in)
	case reflect.Map:
		return enc.fitMapValue(in)
	}
	return in
}

func (enc *Encoder) fitMapValue(in reflect.Value) reflect.Value {
	out := reflect.MakeMap(in.Type())
	it := in.MapRange()
	for it.Next() {
		out.SetMapIndex(it.Key(), enc.fitValue(reflect.ValueOf(it.Value().Interface())))
	}
	return out
}

func (enc *Encoder) fitSliceValue(in reflect.Value) reflect.Value {
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
			out = reflect.Append(out, enc.fitValue(reflect.ValueOf(in.Index(i).Interface())))
		}
		return out
	}

	out := reflect.MakeSlice(in.Type(), 0, in.Len())
	for i := 0; i < in.Len(); i++ {
		out = reflect.Append(out, enc.fitValue(reflect.ValueOf(in.Index(i).Interface())))
	}

	return out
}

func (enc *Encoder) Encode(e interface{}) error {
	return enc.EncodeValue(reflect.ValueOf(e))
}

func (enc *Encoder) EncodeValue(value reflect.Value) error {
	if enc.Fit {
		return enc.Encoder.EncodeValue(enc.fitValue(value))
	}
	return enc.Encoder.EncodeValue(value)
}
