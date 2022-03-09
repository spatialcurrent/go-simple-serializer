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

	"github.com/spatialcurrent/go-object/pkg/object"
	"github.com/spatialcurrent/go-simple-serializer/pkg/escaper"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

func marshalTag(keyValueSeparator []byte, keySerializer stringify.Stringer, valueSerializer stringify.Stringer, e *escaper.Escaper, obj object.Object, key interface{}) ([]byte, error) {
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

	value, valueStringError := valueSerializer(obj.Index(key).Value())
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
func Marshal(obj interface{}, keys object.ObjectArray, expandKeys bool, keyValueSeparator string, keySerializer stringify.Stringer, valueSerializer stringify.Stringer, sorted bool, reversed bool) ([]byte, error) {

	concrete := object.NewObject(obj).Concrete()

	if keySerializer == nil {
		return make([]byte, 0), ErrMissingKeySerializer
	}

	if valueSerializer == nil {
		return make([]byte, 0), ErrMissingValueSerializer
	}

	e := escaper.New().Prefix("\\").Sub("\"", "\n")

	switch concrete.Kind() {
	case reflect.Map:
		out := &bytes.Buffer{}
		if keys.Length() > 0 {
			allKeys := keys
			if expandKeys {
				allKeys = keys.Append(concrete.Keys().Unique().Subtract(keys.Unique()).Array())
			}
			for i, key := range allKeys.Value() {
				b, err := marshalTag(
					[]byte(keyValueSeparator),
					keySerializer,
					valueSerializer,
					e,
					concrete,
					key)
				if err != nil {
					return make([]byte, 0), fmt.Errorf("error serializing tag: %w", err)
				}
				_, err = out.Write(b)
				if err != nil {
					return make([]byte, 0), fmt.Errorf("error writing tag: %w", err)
				}
				if i < allKeys.Length()-1 {
					_, err = out.WriteRune(space)
					if err != nil {
						return make([]byte, 0), fmt.Errorf("error writing space: %w", err)
					}
				}
			}
		} else {
			allKeys := concrete.Keys()
			if sorted {
				allKeys = allKeys.Sort(reversed)
			}
			for i, key := range allKeys.Value() {
				b, err := marshalTag(
					[]byte(keyValueSeparator),
					keySerializer,
					valueSerializer,
					e,
					concrete,
					key)
				if err != nil {
					return make([]byte, 0), fmt.Errorf("error serializing tag: %w", err)
				}
				_, err = out.Write(b)
				if err != nil {
					return make([]byte, 0), fmt.Errorf("error writing tag: %w", err)
				}
				if i < allKeys.Length()-1 {
					_, err = out.WriteRune(space)
					if err != nil {
						return make([]byte, 0), fmt.Errorf("error writing space: %w", err)
					}
				}
			}
		}
		return out.Bytes(), nil
	case reflect.Struct:
		if keys.Length() > 0 {
			allFieldNames := keys.StringArray()
			if expandKeys {
				allFieldNames = keys.StringArray().Append(concrete.FieldNames().Unique().Subtract(keys.Unique().StringArray().Value()...).Array())
			}
			out := &bytes.Buffer{}
			for i, fieldName := range allFieldNames.Value() {
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
				value, err := valueSerializer(concrete.FieldByName(fieldName).Value())
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
				if i < allFieldNames.Length()-1 {
					_, err = out.WriteRune(space)
					if err != nil {
						return make([]byte, 0), fmt.Errorf("error writing space: %w", err)
					}
				}
			}
			return out.Bytes(), nil
		} else {
			out := &bytes.Buffer{}
			fieldNames := concrete.FieldNames()
			if sorted {
				fieldNames = fieldNames.Sort(reversed)
			}
			for i, fieldName := range fieldNames.Value() {
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
				value, err := valueSerializer(concrete.FieldByName(fieldName).Value())
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
				if i < fieldNames.Length()-1 {
					_, err = out.WriteRune(space)
					if err != nil {
						return make([]byte, 0), fmt.Errorf("error writing space: %w", err)
					}
				}
			}
			return out.Bytes(), nil
		}
	}

	return []byte(""), &ErrInvalidKind{Value: concrete.Type(), Expected: []reflect.Kind{reflect.Map, reflect.Struct}}
}
