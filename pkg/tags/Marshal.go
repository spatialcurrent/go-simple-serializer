// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"reflect"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/escaper"
	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
)

// Marshal
func Marshal(object interface{}, valueSerializer func(object interface{}) (string, error), sorted bool) ([]byte, error) {

	if valueSerializer == nil {
		return make([]byte, 0), ErrMissingValueSerializer
	}

	objectValue := reflect.ValueOf(object)
	for objectValue.Type().Kind() == reflect.Ptr {
		objectValue = objectValue.Elem()
	}

	objectType := objectValue.Type()

	e := escaper.New().Prefix("\\").Sub("\"", "\n")

	switch objectType.Kind() {
	case reflect.Map:
		str := ""
		keys := inspector.GetKeysFromValue(objectValue, sorted)
		for i, key := range keys {
			keyString, keyStringError := valueSerializer(key)
			if keyStringError != nil {
				return []byte(str), errors.Wrap(keyStringError, "error serializing tag key")
			}
			if strings.Contains(keyString, " ") {
				str += "\"" + e.Escape(keyString) + "\""
			} else {
				str += e.Escape(keyString)
			}
			str += "="
			value, valueStringError := valueSerializer(objectValue.MapIndex(reflect.ValueOf(key)).Interface())
			if valueStringError != nil {
				return []byte(str), errors.Wrap(valueStringError, "error serializing tag value")
			}
			if len(value) > 0 {
				if strings.Contains(value, " ") {
					str += "\"" + e.Escape(value) + "\""
				} else {
					str += e.Escape(value)
				}
			}
			if i < len(keys)-1 {
				str += " "
			}
		}
		return []byte(str), nil
	case reflect.Struct:
		str := ""
		fieldNames := inspector.GetFieldNamesFromValue(objectValue, sorted)
		for i, fieldName := range fieldNames {
			str += fieldName + "=" // we already know field name won't need to be esaped.
			value, err := valueSerializer(objectValue.FieldByName(fieldName).Interface())
			if err != nil {
				return []byte(str), errors.Wrap(err, "error serializing tag value")
			}
			if len(value) > 0 {
				if strings.Contains(value, " ") {
					str += "\"" + e.Escape(value) + "\""
				} else {
					str += e.Escape(value)
				}
			}
			if i < len(fieldNames)-1 {
				str += " "
			}
		}
		return []byte(str), nil
	}

	return []byte(""), &ErrInvalidKind{Value: objectType, Expected: []reflect.Kind{reflect.Map, reflect.Struct}}
}
