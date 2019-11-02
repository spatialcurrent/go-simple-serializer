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

	"github.com/spatialcurrent/go-simple-serializer/pkg/fit"
	"github.com/spatialcurrent/go-simple-serializer/pkg/tagger"
)

// Unmarshal unmarshaling the source data into the target converting maps to structs as indicated by struct tags.
func Unmarshal(data interface{}, v interface{}) error {

	// If inputs implements the Marshaler interface.
	if unmarshaler, ok := v.(Unmarshaler); ok {
		return unmarshaler.UnmarshalMap(data)
	}

	sourceValue := reflect.ValueOf(data)
	sourceType := sourceValue.Type()
	sourceKind := sourceType.Kind()

	// Chase pointers to concerete value
	targetValue := reflect.ValueOf(v)
	for reflect.TypeOf(targetValue.Interface()).Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}
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

	if targetKind == reflect.Struct {
		if sourceKind == reflect.Map {
			if !reflect.TypeOf("").AssignableTo(sourceType.Key()) {
				return errors.Errorf("string is not assignable to source map key %q", sourceType.Key())
			}
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
				v := reflect.ValueOf(mv.Interface())
				if !v.Type().AssignableTo(f.Type) {
					// If the raw value is not assignable, then try with the fitted value.
					if fitted := fit.FitValue(v); fitted.Type().AssignableTo(f.Type) {
						fv.Set(fitted)
						continue
					}
					return errors.Errorf("key %q found, but value %#v (%s) not assignable to field %q with type %q", key, v, v, f.Name, f.Type)
				}
				fv.Set(v)
			}
			return nil
		}
	}

	return nil
}
