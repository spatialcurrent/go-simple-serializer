// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package json

import (
	stdjson "encoding/json" // import the standard json library as stdjson
	"fmt"
	"reflect"

	"unicode/utf8"

	"github.com/pkg/errors" // utf8 is used to decode the first rune in the string
)

// UnmarshalType parses a slice of bytes into an object of a given type.
// If no input is given, then returns ErrEmptyInput.
// If the first rune is invalid, then returns ErrInvalidRune.
func UnmarshalType(b []byte, outputType reflect.Type) (interface{}, error) {

	if len(b) == 0 {
		return nil, ErrEmptyInput
	}

	switch string(b) {
	case "true":
		if outputType.Kind() != reflect.Bool {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Bool}}
		}
		return true, nil
	case "false":
		if outputType.Kind() != reflect.Bool {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Bool}}
		}
		return false, nil
	case "null":
		return nil, nil
	}

	first, _ := utf8.DecodeRune(b)
	if first == utf8.RuneError {
		return nil, ErrInvalidRune
	}

	switch first {
	case '[':
		if outputType.Kind() != reflect.Slice {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Slice}}
		}
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeSlice(outputType, 0, 0))
		err := stdjson.Unmarshal(b, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q into %T", string(b), ptr.Interface()))
		}
		return ptr.Elem().Interface(), nil
	case '{':
		switch outputType.Kind() {
		case reflect.Map:
			ptr := reflect.New(outputType)
			ptr.Elem().Set(reflect.MakeMap(outputType))
			err := stdjson.Unmarshal(b, ptr.Interface())
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q into %T", string(b), ptr.Interface()))
			}
			return ptr.Elem().Interface(), nil
		case reflect.Struct:
			ptr := reflect.New(outputType)
			err := stdjson.Unmarshal(b, ptr.Interface())
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q into %T", string(b), ptr.Interface()))
			}
			return ptr.Elem().Interface(), nil
		}
		return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Map}}
	case '"':
		if k := outputType.Kind(); k != reflect.String {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.String}}
		}
		obj := ""
		err := stdjson.Unmarshal(b, &obj)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q into string", string(b)))
		}
		return obj, nil
	}

	if k := outputType.Kind(); k != reflect.Float64 {
		return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Float64}}
	}
	obj := 0.0
	err := stdjson.Unmarshal(b, &obj)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q into float", string(b)))
	}
	return obj, nil
}
