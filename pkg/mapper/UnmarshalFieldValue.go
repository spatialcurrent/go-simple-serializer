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

	"github.com/spatialcurrent/go-simple-serializer/pkg/fit"
)

func UnmarshalFieldValue(source reflect.Value, target reflect.Value) error {

	targetType := target.Type()
	targetKind := target.Kind()

	// if field implements Unmarshaler interface.
	if targetType.Implements(reflect.TypeOf((*Unmarshaler)(nil)).Elem()) {
		value := reflect.New(target.Type().Elem())
		result := value.MethodByName("UnmarshalMap").Call([]reflect.Value{source})
		err := result[0].Interface()
		if err != nil {
			return err.(error)
		}
		target.Set(value)
		return nil
	}

	// If target is a pointer
	if targetKind == reflect.Ptr {

		targetElemType := targetType.Elem()
		targetElemKind := targetElemType.Kind()

		// If target is a pointer to a struct
		if targetElemKind == reflect.Struct {
			if !target.Elem().IsValid() {
				target.Set(reflect.New(target.Type().Elem()))
			}
			return UnmarshalValue(source, target)
		}

	}

	// if target is a pointer to a slice
	if targetKind == reflect.Slice {
		return UnmarshalSliceValue(source, target)
	}

	if source.Type().AssignableTo(target.Type()) {
		target.Set(source)
		return nil
	}

	// If the raw value is not assignable, then try with the fitted value.
	if fitted := fit.FitValue(source); fitted.Type().AssignableTo(target.Type()) {
		target.Set(fitted)
		return nil
	}

	return errors.Errorf("value %#v (%q) not assignable to field type %q", source.Interface(), source.Type(), target.Type())
}
