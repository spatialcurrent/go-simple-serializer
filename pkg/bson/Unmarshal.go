// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package bson

import (
	"fmt"

	mgobson "gopkg.in/mgo.v2/bson"
)

// Unmarshal parses a slice of bytes into an object using a few simple type inference rules.
// This package is useful when your program needs to parse data,
// that you have no a priori awareness of its structure or type.
// If no input is given, then returns ErrEmptyInput.
func Unmarshal(b []byte) (interface{}, error) {

	if len(b) == 0 {
		return nil, ErrEmptyInput
	}

	obj := map[string]interface{}{}
	err := mgobson.Unmarshal(b, &obj)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling BSON %q: %w", string(b), err)
	}

	return obj, nil
}
