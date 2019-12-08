// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package fit

import (
	"reflect"
)

// FitMapValue iterates through and fits the given map value and its underlying values.
func FitMapValue(in reflect.Value) reflect.Value {

	if in.Len() == 0 {
		return in
	}

	keyTypes := map[string]reflect.Type{}
	valueTypes := map[string]reflect.Type{}

	m := reflect.MakeMap(in.Type())

	for it := in.MapRange(); it.Next(); {
		k := reflect.ValueOf(it.Key().Interface())
		v := FitValue(reflect.ValueOf(it.Value().Interface()))
		keyTypes[k.Type().PkgPath()+"."+k.Type().String()] = k.Type()
		valueTypes[v.Type().PkgPath()+"."+v.Type().String()] = v.Type()
		m.SetMapIndex(k, v)
	}

	if len(keyTypes) == 1 {
		var keyType reflect.Type
		for _, v := range keyTypes {
			keyType = v
		}
		if len(valueTypes) == 1 {
			// If 1 key type and 1 value type.
			var valueType reflect.Type
			for _, v := range valueTypes {
				valueType = v
			}
			if in.Type().Key().AssignableTo(keyType) && in.Type().Elem().AssignableTo(valueType) {
				return m
			}
			out := reflect.MakeMap(reflect.MapOf(keyType, valueType))
			for it := m.MapRange(); it.Next(); {
				out.SetMapIndex(reflect.ValueOf(it.Key().Interface()), reflect.ValueOf(it.Value().Interface()))
			}
			return out
		}
		// If 1 key type, but multiple value types.
		out := reflect.MakeMap(reflect.MapOf(keyType, in.Type().Elem()))
		for it := m.MapRange(); it.Next(); {
			out.SetMapIndex(reflect.ValueOf(it.Key().Interface()), reflect.ValueOf(it.Value().Interface()))
		}
		return out
	}

	// If more than 1 key type, but only 1 value type.
	if len(valueTypes) == 1 {
		var valueType reflect.Type
		for _, v := range valueTypes {
			valueType = v
		}
		out := reflect.MakeMap(reflect.MapOf(in.Type().Key(), valueType))
		for it := m.MapRange(); it.Next(); {
			out.SetMapIndex(reflect.ValueOf(it.Key().Interface()), reflect.ValueOf(it.Value().Interface()))
		}
		return out
	}

	// If more than 1 key type and more than 1 value type.
	return m
}
