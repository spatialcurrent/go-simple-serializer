// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package json

import (
	"bytes"
	stdjson "encoding/json" // import the standard json library as stdjson
	"fmt"
	"unicode/utf8"

	"github.com/pkg/errors"
	// utf8 is used to decode the first rune in the string
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
//  - {...} => map[string]interface{}
//  - "..." => string
//  - otherwise trys to parse as float
func Unmarshal(b []byte) (interface{}, error) {

	if len(b) == 0 {
		return nil, ErrEmptyInput
	}

	if bytes.HasPrefix(b, BytesTrue) {
		return true, nil
	}

	if bytes.HasPrefix(b, BytesFalse) {
		return false, nil
	}

	if bytes.HasPrefix(b, BytesNull) {
		return nil, nil
	}

	first, _ := utf8.DecodeRune(b)
	if first == utf8.RuneError {
		return nil, ErrInvalidRune
	}

	switch first {
	case '[':
		obj := make([]interface{}, 0)
		err := stdjson.Unmarshal(b, &obj)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q into %T", string(b), obj))
		}
		return obj, nil
	case '{':
		obj := map[string]interface{}{}
		err := stdjson.Unmarshal(b, &obj)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q into %T", string(b), obj))
		}
		return obj, nil
	case '"':
		obj := ""
		err := stdjson.Unmarshal(b, &obj)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q into string", string(b)))
		}
		return obj, nil
	}

	obj := 0.0
	err := stdjson.Unmarshal(b, &obj)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling JSON %q into float", string(b)))
	}
	return obj, nil
}
