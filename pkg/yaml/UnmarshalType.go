// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package yaml

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8" // utf8 is used to decode the first rune in the string

	goyaml "gopkg.in/yaml.v2" // import the YAML library from https://github.com/go-yaml/yaml
)

// UnmarshalType parses a slice of bytes into an object of a given type.
// If no input is given, then returns ErrEmptyInput.
// If the first rune in the slice of bytes is invalid, then returns ErrInvalidRune.
func UnmarshalType(b []byte, outputType reflect.Type) (interface{}, error) {

	if len(b) == 0 {
		return nil, ErrEmptyInput
	}

	// If the kind of the output type is interface{}, then simply use Unmarshal.
	if outputType.Kind() == reflect.Interface {
		return Unmarshal(b)
	}

	if bytes.Equal(b, Y) || bytes.Equal(b, True) {
		if outputType.Kind() != reflect.Bool {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Bool}}
		}
		return true, nil
	}
	if bytes.Equal(b, False) {
		if outputType.Kind() != reflect.Bool {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Bool}}
		}
		return false, nil
	}
	if bytes.Equal(b, Null) {
		return nil, nil
	}

	if bytes.HasPrefix(b, BoundaryMarker) {
		if outputType.Kind() != reflect.Slice {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Slice}}
		}
		s := NewDocumentScanner(bytes.NewReader(b), true)
		out := reflect.MakeSlice(outputType, 0, 0)
		i := 0
		for s.Scan() {
			if d := s.Bytes(); len(d) > 0 {
				obj, err := UnmarshalType(d, outputType.Elem())
				if err != nil {
					return out.Interface(), fmt.Errorf("error scanning document %d: %w", i, err)
				}
				out = reflect.Append(out, reflect.ValueOf(obj))
				i++
			}
		}
		if err := s.Err(); err != nil {
			return out.Interface(), fmt.Errorf("error scanning YAML %q: %w", string(b), err)
		}
		return out.Interface(), nil
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
			return nil, fmt.Errorf("error unmarshaling YAML %q: %W", string(b), err)
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
			return nil, fmt.Errorf("error unmarshaling YAML %q: %w", string(b), err)
		}
		return ptr.Elem().Interface(), nil
	case '"':
		if k := outputType.Kind(); k != reflect.String {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.String}}
		}
		obj := ""
		err := goyaml.Unmarshal(b, &obj)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling YAML %q: %w", string(b), err)
		}
		return obj, nil
	}

	if _, _, ok := ParseKeyValue(b); ok {
		k := outputType.Kind()

		if k == reflect.Map {
			ptr := reflect.New(outputType)
			ptr.Elem().Set(reflect.MakeMap(outputType))
			err := goyaml.Unmarshal(b, ptr.Interface())
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling YAML %q into map: %w", string(b), err)
			}
			return ptr.Elem().Interface(), nil
		}

		if k == reflect.Struct {
			ptr := reflect.New(outputType)
			ptr.Elem().Set(reflect.Zero(outputType))
			err := goyaml.Unmarshal(b, ptr.Interface())
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling YAML %q into struct: %w", string(b), err)
			}
			return ptr.Elem().Interface(), nil
		}

		return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Map, reflect.Struct}}
	}

	switch outputType.Kind() {
	case reflect.Int:
		i, err := strconv.Atoi(string(b))
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling YAML %q: %w", string(b), err)
		}
		return i, nil
	case reflect.Float64:
		f, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling YAML %q: %w", string(b), err)
		}
		return f, nil
	case reflect.String:
		str := strings.TrimSpace(string(b))
		if len(str) > 0 {
			return str, nil
		}
		// if empty string, then return nil
		return nil, nil
	}

	return nil, fmt.Errorf("could not unmarshal YAML %q into type %v", string(b), outputType)
}
