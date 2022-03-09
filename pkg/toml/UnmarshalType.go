// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package toml

import (
	"fmt"
	"reflect"

	bstoml "github.com/BurntSushi/toml"
)

// import the BurntSushi toml library as bstoml

// UnmarshalType parses a slice of bytes into an object of a given type.
// If no input is given, then returns ErrEmptyInput.
func UnmarshalType(b []byte, outputType reflect.Type) (interface{}, error) {

	if len(b) == 0 {
		return nil, ErrEmptyInput
	}

	if k := outputType.Kind(); k != reflect.Map {
		return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Map}}
	}
	ptr := reflect.New(outputType)
	ptr.Elem().Set(reflect.MakeMap(outputType))
	_, err := bstoml.Decode(string(b), ptr.Interface())
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling TOML %q: %w", string(b), err)
	}
	return ptr.Elem().Interface(), nil
}
