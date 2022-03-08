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
	"strconv"
	"strings"
	"unicode/utf8" // utf8 is used to decode the first rune in the string

	goyaml "gopkg.in/yaml.v2" // import the YAML library from https://github.com/go-yaml/yaml
)

// Unmarshal parses a slice of bytes into an object using a few simple type inference rules.
// This package is useful when your program needs to parse data,
// that you have no a priori awareness of its structure or type.
// If no input is given, then returns ErrEmptyInput.
// If the first rune is invalid, then returns ErrInvalidRune.
//
//  - true => true (bool)
//  - false => false (bool)
//  - null => nil
//  - [...] => []interface{}
//  - -... => []interface{}
//  - {...} => map[string]interface{}
//  - "..." => string
//  - otherwise trys to parse as float
func Unmarshal(b []byte) (interface{}, error) {

	if len(b) == 0 {
		return nil, ErrEmptyInput
	}

	if bytes.Equal(b, Y) || bytes.Equal(b, True) {
		return true, nil
	}
	if bytes.Equal(b, False) {
		return false, nil
	}
	if bytes.Equal(b, Null) {
		return nil, nil
	}

	if bytes.HasPrefix(b, BoundaryMarker) {
		s := NewDocumentScanner(bytes.NewReader(b), true)
		obj := make([]interface{}, 0)
		i := 0
		for s.Scan() {
			if d := s.Bytes(); len(d) > 0 {
				element, err := Unmarshal(d)
				if err != nil {
					return obj, fmt.Errorf("error scanning document %d: %w", i, err)
				}
				obj = append(obj, element)
				i++
			}
		}
		if err := s.Err(); err != nil {
			return obj, fmt.Errorf("error scanning YAML %q: %w", string(b), err)
		}
		return obj, nil
	}

	first, _ := utf8.DecodeRune(b)
	if first == utf8.RuneError {
		return nil, ErrInvalidRune
	}

	switch first {
	case '[', '-':
		obj := make([]interface{}, 0)
		err := goyaml.Unmarshal(b, &obj)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling YAML %q: %w", string(b), err)
		}
		return obj, nil
	case '{':
		obj := map[string]interface{}{}
		err := goyaml.Unmarshal(b, &obj)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling YAML %q: %w", string(b), err)
		}
		return obj, nil
	case '"':
		obj := ""
		err := goyaml.Unmarshal(b, &obj)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling YAML %q: %w", string(b), err)
		}
		return obj, nil
	}

	if _, _, ok := ParseKeyValue(b); ok {
		obj := map[string]interface{}{}
		err := goyaml.Unmarshal(b, &obj)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling YAML %q: %w", string(b), err)
		}
		return obj, nil
	}

	str := string(b)

	i, err := strconv.Atoi(str)
	if err == nil {
		return i, nil
	}

	f, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return f, nil
	}

	str = strings.TrimSpace(str)
	if len(str) > 0 {
		return str, nil
	}
	// if empty string, then return nil
	return nil, nil
}
