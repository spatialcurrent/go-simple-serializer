// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package mapper

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/spatialcurrent/go-simple-serializer/pkg/tagger"
)

// Marshal recursively marshals the given object using the following rules:
//  - if implements Marshaler intereface, then uses the MarshalMap method,
//  - deferences pointers
//  - converts structs to maps
func Marshal(object interface{}) (interface{}, error) {

	v := reflect.ValueOf(object)

	// If value is not valid or nil, return nil.
	if !v.IsValid() {
		return nil, nil
	}

	// Chase pointers
	for v.IsValid() && reflect.ValueOf(v.Interface()).Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// If value is nil and a a pointer, slice or map, then return nil
	if k := v.Kind(); (k == reflect.Ptr || k == reflect.Slice || k == reflect.Map) && v.IsNil() {
		return nil, nil
	}

	c := v.Interface()

	// If inputs implements the Marshaler interface.
	if marshaler, ok := c.(Marshaler); ok {
		return marshaler.MarshalMap()
	}

	// If input is a slice of literals or pointers to literals.
	switch slc := object.(type) {
	case []string:
		return slc, nil
	case []*string:
		out := make([]string, 0, len(slc))
		for i := 0; i < len(slc); i++ {
			out = append(out, *slc[i])
		}
		return out, nil
	case []int:
		return slc, nil
	case []*int:
		out := make([]int, 0, len(slc))
		for i := 0; i < len(slc); i++ {
			out = append(out, *slc[i])
		}
		return out, nil
	case []float64:
		return slc, nil
	case []*float64:
		out := make([]float64, 0, len(slc))
		for i := 0; i < len(slc); i++ {
			out = append(out, *slc[i])
		}
		return out, nil
	}

	// If input is a map of literals to literals.
	switch m := object.(type) {
	case map[int]float64:
		return m, nil
	case map[int]*float64:
		return m, nil
	case map[int]int:
		return m, nil
	case map[int]*int:
		return m, nil
	case map[int]string:
		return m, nil
	case map[int]*string:
		return m, nil

	case map[string]float64:
		return m, nil
	case map[string]*float64:
		return m, nil
	case map[string]int:
		return m, nil
	case map[string]*int:
		return m, nil
	case map[string]string:
		return m, nil
	case map[string]*string:
		return m, nil

	case map[interface{}]float64:
		return m, nil
	case map[interface{}]*float64:
		return m, nil
	case map[interface{}]int:
		return m, nil
	case map[interface{}]*int:
		return m, nil
	case map[interface{}]string:
		return m, nil
	case map[interface{}]*string:
		return m, nil
	}

	in := reflect.ValueOf(c) // sets value to concerete type
	t := v.Type()
	k := t.Kind()

	// If input is of kind slice.
	if k == reflect.Slice {
		out := make([]interface{}, 0)
		for i := 0; i < v.Len(); i++ {
			element, err := Marshal(v.Index(i).Interface())
			if err != nil {
				return nil, errors.Wrapf(err, "error marshaling %#v", v.Index(i).Interface())
			}
			out = append(out, element)
		}
		return out, nil
	}

	// If input is of kind map.
	if k == reflect.Map {
		out := reflect.MakeMapWithSize(reflect.MapOf(t.Key(), interfaceType), v.Len())
		for it := in.MapRange(); it.Next(); {
			v, err := Marshal(it.Value().Interface())
			if err != nil {
				return nil, errors.Wrapf(err, "error marshaling %#v", it.Value().Interface())
			}
			out.SetMapIndex(it.Key(), reflect.ValueOf(v))
		}
		return out.Interface(), nil
	}

	// If input is of kind struct.
	if k == reflect.Struct {

		out := make(map[string]interface{}, in.NumField())
		for i := 0; i < in.NumField(); i++ {
			f := t.Field(i)   // field
			fv := in.Field(i) // field value

			tagValue, err := tagger.Lookup(f.Tag, "map")
			if err != nil {
				return nil, errors.Wrapf(err, "error unmarshaling struct tag value %q", f.Tag)
			}

			key := f.Name
			omitEmpty := false
			if tagValue != nil {
				if tagValue.Ignore {
					continue
				}
				if len(tagValue.Name) > 0 {
					key = tagValue.Name
				}
				omitEmpty = tagValue.OmitEmpty
			}

			// If value is not valid or nil, return nil.
			if !fv.IsValid() {
				if omitEmpty {
					continue
				}
				out[key] = fv
				continue
			}

			// If value is nil
			if k := fv.Kind(); (k == reflect.Ptr || k == reflect.Map || k == reflect.Slice) && fv.IsNil() {
				// If omitempty struct tag attribute was present, then skip.
				if omitEmpty {
					continue
				}
				out[key] = fv
				continue
			}

			// Marshal the underlying value
			mfv, err := Marshal(fv.Interface())
			if err != nil {
				return nil, errors.Wrapf(err, "error marshaling value for field %v", f.Name)
			}

			// If marshaled field value is empty
			if IsEmpty(mfv) && omitEmpty {
				continue
			}

			out[key] = mfv
		}
		return out, nil
	}

	// return concerete value
	return c, nil
}
