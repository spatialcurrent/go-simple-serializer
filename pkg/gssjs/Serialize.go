// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gssjs

import (
	"github.com/pkg/errors"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

// Serialize is a function provided to gss.js that wraps gss.Serialize to support JavaScript.
//func Serialize(inputObject interface{}, outputFormat string, options *js.Object) interface{} {
func Serialize(inputObject interface{}, outputFormat string, options map[string]interface{}) (string, error) {

	if len(outputFormat) == 0 || outputFormat == "undefined" {
		return "", ErrMissingOutputFormat
	}

	s, err := serializer.NewWithOptions(outputFormat, SerializeDefaults, options)
	if err != nil {
		return "", errors.Wrap(err, "error creating serializer")
	}

	outputBytes, err := s.Serialize(inputObject)
	if err != nil {
		return "", errors.Wrap(err, "error serializing")
	}
	return string(outputBytes), nil
}
