// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package mapper

import (
	"fmt"
	"reflect"

	"github.com/spatialcurrent/go-simple-serializer/pkg/tagger"
)

// Unmarshal unmarshaling the source data into the target converting maps to structs as indicated by struct tags.
func Unmarshal(data interface{}, v interface{}) error {
	return UnmarshalValue(reflect.ValueOf(data), reflect.ValueOf(v))
}

func UnmarshalValue(sourceValue reflect.Value, targetValue reflect.Value) error {

	targetType := targetValue.Type()
	targetKind := targetType.Kind()

	// If target is of kind pointer, then dereference the target value.
	if targetKind == reflect.Ptr {
		return UnmarshalValue(sourceValue, targetValue.Elem())
	}

	if !targetValue.CanAddr() {
		return fmt.Errorf("target %#v (%v) is not addressable", targetValue.Interface(), targetValue.Type())
	}

	if !targetValue.CanSet() {
		return fmt.Errorf("target %#v (%v) cannot be set", targetValue.Interface(), targetValue.Type())
	}

	// If source value is not valid, then set the targetValue to it's zero value.
	// This can occur, when given a "reflect.ValueOf(nil)"" sourceValue.
	if !sourceValue.IsValid() {
		targetValue.Set(reflect.New(targetValue.Type()).Elem())
		return nil
	}

	sourceType := sourceValue.Type()
	sourceKind := sourceType.Kind()

	// If source is of kind pointer, then dereference the source value.
	if sourceKind == reflect.Ptr {
		return UnmarshalValue(sourceValue.Elem(), targetValue)
	}

	if sourceKind == reflect.Interface {
		if !sourceValue.CanInterface() {
			return fmt.Errorf("source %v (%v) is of kind interface", sourceValue, sourceValue.Type())
		}
		// Re-value the object
		return UnmarshalValue(reflect.ValueOf(sourceValue.Interface()), targetValue)
	}

	// If target implements the unmarshaler interface
	if reflect.PtrTo(targetType).Implements(reflect.TypeOf((*Unmarshaler)(nil)).Elem()) {
		outValue := reflect.New(targetType)
		result := outValue.MethodByName("UnmarshalMap").Call([]reflect.Value{sourceValue})
		err := result[0].Interface()
		if err != nil {
			return err.(error)
		}
		targetValue.Set(outValue.Elem())
		return nil
	}

	// If target is a slice
	if targetKind == reflect.Slice {
		return UnmarshalSliceValue(sourceValue, targetValue)
	}

	// If target is a map
	if targetKind == reflect.Map {
		return UnmarshalMapValue(sourceValue, targetValue)
	}

	// If target is a struct
	if targetKind == reflect.Struct {

		if sourceKind == reflect.Map {

			if !reflect.TypeOf("").AssignableTo(sourceType.Key()) {
				return fmt.Errorf("string is not assignable to source map key %q", sourceType.Key())
			}

			// Iterate throught the struct fields
			for i := 0; i < targetValue.NumField(); i++ {
				f := targetType.Field(i)   // field
				fv := targetValue.Field(i) // field value

				if f.Anonymous {
					continue
				}

				if !fv.CanSet() {
					continue
				}

				tagValue, err := tagger.Lookup(f.Tag, "map")
				if err != nil {
					return fmt.Errorf("error unmarshaling struct tag value %q: %w", f.Tag, err)
				}

				key := f.Name
				if tagValue != nil {
					if tagValue.Ignore {
						continue
					}
					if len(tagValue.Name) > 0 {
						key = tagValue.Name
					}
				}

				mv := sourceValue.MapIndex(reflect.ValueOf(key))
				if !mv.IsValid() {
					// if key was not found
					continue
				}

				// If source value is nil, then set target to its zero value.
				if k := mv.Kind(); k == reflect.Map || k == reflect.Interface {
					if mv.IsNil() {
						fv.Set(reflect.Zero(fv.Type()))
						continue
					}
				}

				// unmarshal the concrete map value into the field
				err = UnmarshalFieldValue(reflect.ValueOf(mv.Interface()), fv)
				if err != nil {
					return fmt.Errorf("key %q found, but could not assign to field %q: %w", key, f.Name, err)
				}
			}
			return nil
		}
	}

	// If source is assignale to target, then simply set it.
	if sourceType.AssignableTo(targetType) {
		targetValue.Set(sourceValue)
	}

	return nil
}
