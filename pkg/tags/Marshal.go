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
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/escaper"
	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
)

// Marshal formats an object into a slice of bytes of tags.
// The value serializer is used to render the key and value of each pair into strings.
// If sorted and not reversed, then the keys are sorted in alphabetical order.
// If sorted and reversed, then the keys are sorted in reverse alphabetical order.
func Marshal(object interface{}, keySerializer stringify.Stringer, valueSerializer stringify.Stringer, sorted bool, reversed bool) ([]byte, error) {

	if keySerializer == nil {
		return make([]byte, 0), ErrMissingKeySerializer
	}

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
		keys := inspector.GetKeysFromValue(objectValue, sorted, reversed)
		for i, key := range keys {
			keyString, errKeyString := keySerializer(key)
			if errKeyString != nil {
				return []byte(str), errors.Wrap(errKeyString, "error serializing tag key")
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
		fieldNames := inspector.GetFieldNamesFromValue(objectValue, sorted, reversed)
		for i, fieldName := range fieldNames {
			keyString, errKeyString := keySerializer(fieldName)
			if errKeyString != nil {
				return []byte(str), errors.Wrap(errKeyString, "error serializing tag key")
			}
			str += keyString + "="
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
