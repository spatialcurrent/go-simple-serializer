// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package mapper

import (
	//"fmt"
	"github.com/pkg/errors"
	"reflect"

	"github.com/spatialcurrent/go-simple-serializer/pkg/fit"
	"github.com/spatialcurrent/go-simple-serializer/pkg/tagger"
)

func unmarshalFieldValue(mapValue reflect.Value, target reflect.Value) error {
	// if field implements Unmarshaler interface.
	if target.Type().Implements(reflect.TypeOf((*Unmarshaler)(nil)).Elem()) {
		value := reflect.New(target.Type().Elem())
		result := value.MethodByName("UnmarshalMap").Call([]reflect.Value{mapValue})
		err := result[0].Interface()
		if err != nil {
			return err.(error)
		}
		target.Set(value)
		return nil
	}

	// If target is a pointer to a struct, then recurse.
	if target.Type().Kind() == reflect.Ptr {
		if target.Type().Elem().Kind() == reflect.Struct {
			if !target.Elem().IsValid() {
				target.Set(reflect.New(target.Type().Elem()))
			}
			return UnmarshalValue(mapValue, target)
		}
	}

	if !mapValue.Type().AssignableTo(target.Type()) {
		// If the raw value is not assignable, then try with the fitted value.
		if fitted := fit.FitValue(mapValue); fitted.Type().AssignableTo(target.Type()) {
			target.Set(fitted)
			return nil
		}
		return errors.Errorf("value %#v (%q) not assignable to field type %q", mapValue.Interface(), mapValue.Type(), target.Type())
	}
	target.Set(mapValue)
	return nil
}

// Unmarshal unmarshaling the source data into the target converting maps to structs as indicated by struct tags.
func Unmarshal(data interface{}, v interface{}) error {
	if unmarshaler, ok := v.(Unmarshaler); ok {
		return unmarshaler.UnmarshalMap(data)
	}
	return UnmarshalValue(reflect.ValueOf(data), reflect.ValueOf(v))
}

func UnmarshalValue(sourceValue reflect.Value, targetValue reflect.Value) error {
	//fmt.Println(fmt.Sprintf("UnmarshalValue(%#v, %#v)", sourceValue, targetValue))

	sourceType := sourceValue.Type()
	sourceKind := sourceType.Kind()

	//for reflect.TypeOf(targetValue.Interface()).Kind() == reflect.Ptr {
	//	targetValue = targetValue.Elem()
	//}
	targetType := targetValue.Type()
	targetKind := targetType.Kind()

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

	// If target is a pointer
	if targetKind == reflect.Ptr {

		targetElemValue := targetValue.Elem()
		targetElemType := targetValue.Type().Elem()
		targetElemKind := targetElemType.Kind()

		// If target element is a struct
		if targetElemKind == reflect.Struct {
			if sourceKind == reflect.Map {
				if !reflect.TypeOf("").AssignableTo(sourceType.Key()) {
					return errors.Errorf("string is not assignable to source map key %q", sourceType.Key())
				}

				// Iterate throught the struct fields
				for i := 0; i < targetElemValue.NumField(); i++ {
					f := targetElemType.Field(i)   // field
					fv := targetElemValue.Field(i) // field value
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
					// unmarshal the concrete map value into the field
					err = unmarshalFieldValue(reflect.ValueOf(mv.Interface()), fv)
					if err != nil {
						return errors.Wrapf(err, "key %q found, but could not assign to field %q", key, f.Name)
					}
				}
				return nil
			}
		}

	}

	return nil
}
