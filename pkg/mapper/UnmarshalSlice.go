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
)

// UnmarshalSlice unmarshals the given array or slice into the value, and returns an error, if any.
func UnmarshalSlice(data interface{}, v interface{}) error {
	if unmarshaler, ok := v.(Unmarshaler); ok {
		return unmarshaler.UnmarshalMap(data)
	}
	return UnmarshalSliceValue(reflect.ValueOf(data), reflect.ValueOf(v))
}

// UnmarshalSliceValue unmarshals the given array or slice value into the pointer to slice value, and returns an error, if any.
func UnmarshalSliceValue(source reflect.Value, target reflect.Value) error {

	sourceKind := source.Kind()

	// Only accept array or slice input
	if sourceKind != reflect.Array && sourceKind != reflect.Slice {
		return errors.Errorf("source is of type %v, expecting kind of array or slice", source.Type())
	}

	targetType := target.Type()
	targetKind := target.Kind()

	// If target is of kind pointer, then dereference the target value.
	if targetKind == reflect.Ptr {
		return UnmarshalSliceValue(source, target.Elem())
	}

	if !target.CanAddr() {
		return errors.Errorf("target %#v (%T) is not addressable", target, target)
	}

	if !target.CanSet() {
		return errors.Errorf("target %#v (%T) cannot be set", target, target)
	}

	if targetKind != reflect.Slice {
		return errors.Errorf("target element is of type %v, expecting kind of slice", targetType)
	}

	if source.Len() == 0 {
		target.Set(reflect.MakeSlice(targetType, 0, 0))
		return nil
	}

	// create the output slice
	out := reflect.MakeSlice(targetType, 0, source.Len())

	if targetType.Elem().Kind() == reflect.Ptr {
		for i := 0; i < source.Len(); i++ {
			v := reflect.New(targetType.Elem().Elem())
			v.Elem().Set(reflect.Zero(targetType.Elem().Elem()))
			err := UnmarshalValue(reflect.ValueOf(source.Index(i).Interface()), v.Elem())
			if err != nil {
				return errors.Wrapf(err, "error unmarshaling slice value %d", i)
			}
			out = reflect.Append(out, v)
		}
	} else {
		for i := 0; i < source.Len(); i++ {
			v := reflect.New(targetType.Elem())
			if targetType.Elem().Kind() == reflect.Map {
				v.Elem().Set(reflect.MakeMap(targetType.Elem()))
			} else {
				v.Elem().Set(reflect.Zero(targetType.Elem()))
			}
			err := UnmarshalValue(reflect.ValueOf(source.Index(i).Interface()), v.Elem())
			if err != nil {
				return errors.Wrapf(err, "error unmarshaling slice value %d", i)
			}
			out = reflect.Append(out, v.Elem())
		}
	}

	// Set target to the new slice
	target.Set(out)

	return nil
}
