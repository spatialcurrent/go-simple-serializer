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

	"github.com/pkg/errors"
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

	if bytes.Equal(b, True) {
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
					return out.Interface(), errors.Wrapf(err, "error scanning document %d", i)
				}
				out = reflect.Append(out, reflect.ValueOf(obj))
				i++
			}
		}
		if err := s.Err(); err != nil {
			return out.Interface(), errors.Wrap(err, fmt.Sprintf("error scanning YAML %q", string(b)))
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
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling YAML %q", string(b)))
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
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling YAML %q", string(b)))
		}
		return ptr.Elem().Interface(), nil
	case '"':
		if k := outputType.Kind(); k != reflect.String {
			return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.String}}
		}
		obj := ""
		err := goyaml.Unmarshal(b, &obj)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling YAML %q", string(b)))
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
