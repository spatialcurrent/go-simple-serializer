// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gssjs

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

// Deserialize is a function provided to gss.js that wraps gss.Deserialize to support JavaScript.
//func Deserialize(inputString string, inputFormat string, options *js.Object) interface{} {
func Deserialize(inputString string, inputFormat string, options map[string]interface{}) (interface{}, error) {

	if len(inputFormat) == 0 || inputFormat == "undefined" {
		return "", ErrMissingInputFormat
	}

	s, err := serializer.NewWithOptions(inputFormat, DeserializeDefaults, options)
	if err != nil {
		return nil, errors.Wrap(err, "error creating serializer")
	}

	out, err := s.Deserialize([]byte(inputString))
	if err != nil {
		return nil, errors.Wrap(err, "error deserializing")
	}
	return out, nil
}
