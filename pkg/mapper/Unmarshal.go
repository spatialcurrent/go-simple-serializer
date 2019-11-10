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

// Unmarshal unmarshaling the source data into the target converting maps to structs as indicated by struct tags.
func Unmarshal(data interface{}, v interface{}) error {
	return UnmarshalValue(reflect.ValueOf(data), reflect.ValueOf(v))
}

func UnmarshalValue(sourceValue reflect.Value, targetValue reflect.Value) error {

	sourceType := sourceValue.Type()
	sourceKind := sourceType.Kind()

	if sourceKind == reflect.Interface {
		return errors.Errorf("source %v (%T) is of kind interface", sourceValue, sourceValue)
	}

	//for reflect.TypeOf(targetValue.Interface()).Kind() == reflect.Ptr {
	//	targetValue = targetValue.Elem()
	//}
	targetType := targetValue.Type()
	targetKind := targetType.Kind()

	// If target is of kind pointer, then dereference the target value.
	if targetKind == reflect.Ptr {
		return UnmarshalValue(sourceValue, targetValue.Elem())
	}

	if !targetValue.CanAddr() {
		return errors.Errorf("target %#v (%v) is not addressable", targetValue.Interface(), targetValue.Type())
	}

	if !targetValue.CanSet() {
		return errors.Errorf("target %#v (%v) cannot be set", targetValue.Interface(), targetValue.Type())
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

	if targetKind == reflect.Slice {
		return UnmarshalSliceValue(sourceValue, targetValue)
	}

	// If target is a map
	if targetKind == reflect.Map {
		if sourceKind == reflect.Map {
			if !sourceType.Key().AssignableTo(targetType.Key()) {
				return errors.Errorf("source map key %q is not assignable to target map key %q", sourceType.Key(), targetType.Key())
			}
			for it := sourceValue.MapRange(); it.Next(); {
				v := reflect.New(targetType.Elem())
				err := Unmarshal(it.Value().Interface(), v.Interface())
				if err != nil {
					return errors.Wrapf(err, "error unmarshaling %#v", it.Value().Interface())
				}
				if t := reflect.TypeOf(v); !t.AssignableTo(targetType.Elem()) {
					return errors.Errorf("source map value %v is not assignable to target map value %v", t, targetType.Elem())
				}
				targetValue.SetMapIndex(it.Key(), v)
			}
		}
		return nil
	}

	// If target is a struct
	if targetKind == reflect.Struct {

		if sourceKind == reflect.Map {

			if !reflect.TypeOf("").AssignableTo(sourceType.Key()) {
				return errors.Errorf("string is not assignable to source map key %q", sourceType.Key())
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
					return errors.Wrapf(err, "error unmarshaling struct tag value %q", f.Tag)
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
					return errors.Wrapf(err, "key %q found, but could not assign to field %q", key, f.Name)
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
