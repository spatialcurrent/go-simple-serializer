// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package yaml

import (
	"fmt"
	// import the standard json library as stdjson
	goyaml "gopkg.in/yaml.v2"
	"reflect"
	"strconv"
	"strings"
	// utf8 is used to decode the first rune in the string
	"unicode/utf8"
)

import (
	"github.com/pkg/errors"
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
	case '[', '-':
		if outputType.Kind() != reflect.Slice {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Slice}}
		}
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeSlice(outputType, 0, 0))
		err := goyaml.Unmarshal(b, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q", string(b)))
		}
		return ptr.Elem().Interface(), nil
	case '{':
		if k := outputType.Kind(); k != reflect.Map {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Map}}
		}
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeMap(outputType))
		err := goyaml.Unmarshal(b, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q", string(b)))
		}
		return ptr.Elem().Interface(), nil
	case '"':
		if k := outputType.Kind(); k != reflect.String {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.String}}
		}
		obj := ""
		err := goyaml.Unmarshal(b, &obj)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q", string(b)))
		}
		return obj, nil
	}

	if strings.Contains(string(b), "\n") {
		if k := outputType.Kind(); k != reflect.Map {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Map}}
		}
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeMap(outputType))
		err := goyaml.Unmarshal(b, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling YAML %q", string(b)))
		}
		return ptr.Elem().Interface(), nil
	}

	switch outputType.Kind() {
	case reflect.Int:
		i, err := strconv.Atoi(string(b))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling YAML %q", string(b)))
		}
		return i, nil
	case reflect.Float64:
		f, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling YAML %q", string(b)))
		}
		return f, nil
	}
	return string(b), nil
}
