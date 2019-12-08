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

// UnmarshalMap unmarshals the given map into the value, and returns an error, if any.
func UnmarshalMap(data interface{}, v interface{}) error {
	if unmarshaler, ok := v.(Unmarshaler); ok {
		return unmarshaler.UnmarshalMap(data)
	}
	return UnmarshalMapValue(reflect.ValueOf(data), reflect.ValueOf(v))
}

// UnmarshalMapValue unmarshals the given map value into the target map, and returns an error, if any.
func UnmarshalMapValue(sourceValue reflect.Value, targetValue reflect.Value) error {

	targetType := targetValue.Type()
	targetKind := targetValue.Kind()

	// If target is of kind pointer, then dereference the target value.
	if targetKind == reflect.Ptr {
		return UnmarshalMapValue(sourceValue, targetValue.Elem())
	}

	if !targetValue.CanAddr() {
		return errors.Errorf("target %#v (%T) is not addressable", targetValue, targetValue)
	}

	if !targetValue.CanSet() {
		return errors.Errorf("target %#v (%T) cannot be set", targetValue, targetValue)
	}

	if targetKind != reflect.Map {
		return errors.Errorf("target element is of type %v, expecting kind of map", targetType)
	}

	if !sourceValue.IsValid() {
		targetValue.Set(reflect.MakeMap(targetType))
		return nil
	}

	sourceType := sourceValue.Type()
	sourceKind := sourceValue.Kind()

	// If source is of kind pointer, then dereference the source value.
	if sourceKind == reflect.Ptr {
		return UnmarshalMapValue(sourceValue.Elem(), targetValue)
	}

	// Only accept map input
	if sourceKind != reflect.Map {
		return errors.Errorf("source is of type %v, expecting kind of map", sourceValue.Type())
	}

	if sourceValue.Len() == 0 {
		targetValue.Set(reflect.MakeMap(targetType))
		return nil
	}

	if !sourceType.Key().AssignableTo(targetType.Key()) {
		return errors.Errorf("source map key %q is not assignable to target map key %q", sourceType.Key(), targetType.Key())
	}

	// create the output slice
	out := reflect.MakeMap(targetType)

	for it := sourceValue.MapRange(); it.Next(); {
		v := reflect.New(targetType.Elem())
		err := UnmarshalValue(it.Value(), v.Elem())
		if err != nil {
			return errors.Wrapf(err, "error unmarshaling %#v", it.Value().Interface())
		}
		out.SetMapIndex(it.Key(), v.Elem())
	}

	// Set target to the new slice
	targetValue.Set(out)

	return nil
}
