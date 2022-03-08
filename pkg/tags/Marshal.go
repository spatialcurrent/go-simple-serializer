// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/spatialcurrent/go-simple-serializer/pkg/escaper"
	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

func marshalTag(keyValueSeparator []byte, keySerializer stringify.Stringer, valueSerializer stringify.Stringer, e *escaper.Escaper, objectValue reflect.Value, key interface{}) ([]byte, error) {
	keyString, errKeyString := keySerializer(key)
	if errKeyString != nil {
		return make([]byte, 0), fmt.Errorf("error serializing tag key: %w", errKeyString)
	}

	out := &bytes.Buffer{}

	if strings.Contains(keyString, " ") {
		out.WriteString("\"" + e.Escape(keyString) + "\"")
	} else {
		out.WriteString(e.Escape(keyString))
	}

	_, err := out.Write(keyValueSeparator)
	if err != nil {
		return make([]byte, 0), fmt.Errorf("error writing key-value separator: %w", err)
	}

	value, valueStringError := valueSerializer(objectValue.MapIndex(reflect.ValueOf(key)).Interface())
	if valueStringError != nil {
		return make([]byte, 0), fmt.Errorf("error serializing tag value: %w", valueStringError)
	}

	if len(value) > 0 {
		if strings.Contains(value, " ") {
			out.WriteString("\"" + e.Escape(value) + "\"")
		} else {
			out.WriteString(e.Escape(value))
		}
	}

	return out.Bytes(), nil
}

// Marshal formats an object into a slice of bytes of tags (aka key=value pairs)
// The value serializer is used to render the key and value of each pair into strings.
// If keys is not empty, then prints the tags in the order specifed by keys.
// If expandKeys is true, then adds unknown keys to the end of the list of tags.
// If sorted and not reversed, then the keys are sorted in alphabetical order.
// If sorted and reversed, then the keys are sorted in reverse alphabetical order.
func Marshal(object interface{}, keys []interface{}, expandKeys bool, keyValueSeparator string, keySerializer stringify.Stringer, valueSerializer stringify.Stringer, sorted bool, reversed bool) ([]byte, error) {

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
		out := &bytes.Buffer{}
		if len(keys) > 0 {
			allKeys := keys
			if expandKeys {
				knownKeys := map[interface{}]struct{}{}
				for _, k := range keys {
					knownKeys[k] = struct{}{}
				}
				unknownKeys := inspector.GetUnknownKeysFromValue(objectValue, knownKeys, sorted, reversed)
				allKeys = append(keys, unknownKeys...)
			}
			for i, key := range allKeys {
				b, err := marshalTag(
					[]byte(keyValueSeparator),
					keySerializer,
					valueSerializer,
					e,
					objectValue,
					key)
				if err != nil {
					return make([]byte, 0), fmt.Errorf("error serializing tag: %w", err)
				}
				_, err = out.Write(b)
				if err != nil {
					return make([]byte, 0), fmt.Errorf("error writing tag: %w", err)
				}
				if i < len(allKeys)-1 {
					_, err = out.WriteRune(space)
					if err != nil {
						return make([]byte, 0), fmt.Errorf("error writing space: %w", err)
					}
				}
			}
		} else {
			allKeys := inspector.GetKeysFromValue(objectValue, sorted, reversed)
			for i, key := range allKeys {
				b, err := marshalTag(
					[]byte(keyValueSeparator),
					keySerializer,
					valueSerializer,
					e,
					objectValue,
					key)
				if err != nil {
					return make([]byte, 0), fmt.Errorf("error serializing tag: %w", err)
				}
				_, err = out.Write(b)
				if err != nil {
					return make([]byte, 0), fmt.Errorf("error writing tag: %w", err)
				}
				if i < len(allKeys)-1 {
					_, err = out.WriteRune(space)
					if err != nil {
						return make([]byte, 0), fmt.Errorf("error writing space: %w", err)
					}
				}
			}
		}
		return out.Bytes(), nil
	case reflect.Struct:
		if len(keys) > 0 {
			allFieldNames := make([]string, 0)
			if expandKeys {
				knownFieldNames := map[string]struct{}{}
				for _, k := range keys {
					if str, ok := k.(string); ok {
						allFieldNames = append(allFieldNames, str)
						knownFieldNames[str] = struct{}{}
					}
				}
				unknownFieldNames := inspector.GetUnknownFieldNamesFromValue(objectValue, knownFieldNames, sorted, reversed)
				allFieldNames = append(allFieldNames, unknownFieldNames...)
			} else {
				for _, k := range keys {
					if str, ok := k.(string); ok {
						allFieldNames = append(allFieldNames, str)
					}
				}
			}
			out := &bytes.Buffer{}
			for i, fieldName := range allFieldNames {
				keyString, errKeyString := keySerializer(fieldName)
				if errKeyString != nil {
					return out.Bytes(), fmt.Errorf("error serializing tag key: %w", errKeyString)
				}
				_, err := out.WriteString(keyString)
				if err != nil {
					return out.Bytes(), fmt.Errorf("error writing tag key: %w", err)
				}
				_, err = out.WriteString(keyValueSeparator)
				if err != nil {
					return out.Bytes(), fmt.Errorf("error writing key-value separator: %w", err)
				}
				value, err := valueSerializer(objectValue.FieldByName(fieldName).Interface())
				if err != nil {
					return out.Bytes(), fmt.Errorf("error serializing tag value: %w", err)
				}
				if len(value) > 0 {
					if strings.Contains(value, " ") {
						out.WriteString("\"" + e.Escape(value) + "\"")
					} else {
						out.WriteString(e.Escape(value))
					}
				}
				if i < len(allFieldNames)-1 {
					_, err = out.WriteRune(space)
					if err != nil {
						return make([]byte, 0), fmt.Errorf("error writing space: %w", err)
					}
				}
			}
			return out.Bytes(), nil
		} else {
			out := &bytes.Buffer{}
			fieldNames := inspector.GetFieldNamesFromValue(objectValue, sorted, reversed)
			for i, fieldName := range fieldNames {
				keyString, errKeyString := keySerializer(fieldName)
				if errKeyString != nil {
					return out.Bytes(), fmt.Errorf("error serializing tag key: %w", errKeyString)
				}
				_, err := out.WriteString(keyString)
				if err != nil {
					return out.Bytes(), fmt.Errorf("error writing tag key: %w", err)
				}
				_, err = out.WriteString(keyValueSeparator)
				if err != nil {
					return out.Bytes(), fmt.Errorf("error writing key-value separator: %w", err)
				}
				value, err := valueSerializer(objectValue.FieldByName(fieldName).Interface())
				if err != nil {
					return out.Bytes(), fmt.Errorf("error serializing tag value: %w", err)
				}
				if len(value) > 0 {
					if strings.Contains(value, " ") {
						out.WriteString("\"" + e.Escape(value) + "\"")
					} else {
						out.WriteString(e.Escape(value))
					}
				}
				if i < len(fieldNames)-1 {
					_, err = out.WriteRune(space)
					if err != nil {
						return make([]byte, 0), fmt.Errorf("error writing space: %w", err)
					}
				}
			}
			return out.Bytes(), nil
		}
	}

	return []byte(""), &ErrInvalidKind{Value: objectType, Expected: []reflect.Kind{reflect.Map, reflect.Struct}}
}
