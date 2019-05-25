// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package toml

import (
	"fmt"
	// import the BurntSushi toml library as bstoml
	bstoml "github.com/BurntSushi/toml"
)

import (
	"github.com/pkg/errors"
)

// Unmarshal parses a slice of bytes into a map[string]interface{} object
// This package is useful when your program needs to parse data,
// that you have no a priori awareness of its structure or type.
// If no input is given, then returns ErrEmptyInput.
func Unmarshal(b []byte) (interface{}, error) {

	if len(b) == 0 {
		return nil, ErrEmptyInput
	}

	obj := map[string]interface{}{}
	_, err := bstoml.Decode(string(b), &obj)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling TOML %q", string(b)))
	}
	return obj, nil
}
