// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package yaml

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8" // utf8 is used to decode the first rune in the string

	goyaml "gopkg.in/yaml.v2" // import the YAML library from https://github.com/go-yaml/yaml
)

import (
	"github.com/pkg/errors"
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

	switch string(b) {
	case "true":
		return true, nil
	case "false":
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
		obj := make([]interface{}, 0)
		err := goyaml.Unmarshal(b, &obj)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q", string(b)))
		}
		return obj, nil
	case '{':
		obj := map[string]interface{}{}
		err := goyaml.Unmarshal(b, &obj)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q", string(b)))
		}
		return obj, nil
	case '"':
		obj := ""
		err := goyaml.Unmarshal(b, &obj)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q", string(b)))
		}
		return obj, nil
	}

	if strings.Contains(string(b), "\n") {
		obj := map[string]interface{}{}
		err := goyaml.Unmarshal(b, &obj)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling YAML %q", string(b)))
		}
		return obj, nil
	}

	i, err := strconv.Atoi(string(b))
	if err == nil {
		return i, nil
	}
	f, err := strconv.ParseFloat(string(b), 64)
	if err == nil {
		return f, nil
	}
	return string(b), nil
}
