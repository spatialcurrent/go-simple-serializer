// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package mapper

import (
	"github.com/pkg/errors"
	"reflect"

	"github.com/spatialcurrent/go-simple-serializer/pkg/tagger"
)

func Marshal(object interface{}) (interface{}, error) {

	v := reflect.ValueOf(object)

	// If value is not valid or nil, return nil.
	if !v.IsValid() {
		return nil, nil
	}

	// If value is a pointer and nil, then return nil
	if k := v.Kind(); (k == reflect.Ptr || k == reflect.Map) && v.IsNil() {
		return nil, nil
	}

	// Chase pointers
	for reflect.TypeOf(v.Interface()).Kind() == reflect.Ptr {
		v = v.Elem()
	}

	c := v.Interface()

	// If inputs implements the Marshaler interface.
	if marshaler, ok := c.(Marshaler); ok {
		return marshaler.MarshalMap()
	}

	in := reflect.ValueOf(c) // sets value to concerete type
	t := v.Type()
	k := t.Kind()

	// If input is of kind map
	if k == reflect.Map {
		out := reflect.MakeMapWithSize(t, v.Len())
		for it := in.MapRange(); it.Next(); {
			v, err := Marshal(it.Value().Interface())
			if err != nil {
				return nil, errors.Wrapf(err, "error marshaling %#v", it.Value().Interface())
			}
			out.SetMapIndex(it.Key(), reflect.ValueOf(v))
		}
		return out.Interface(), nil
	}

	// If input is of kind struct
	if k == reflect.Struct {

		out := make(map[string]interface{}, in.NumField())
		for i := 0; i < in.NumField(); i++ {
			f := t.Field(i)               // field
			fv := in.Field(i).Interface() // field value

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

			// Marshal the underlying value
			mfv, err := Marshal(fv)
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

	return object, nil
}
