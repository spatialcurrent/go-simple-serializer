// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gssjs

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

// Convert is a function provided to gss.js that wraps gss.Convert to support JavaScript.
func Convert(inputString string, inputFormat string, outputFormat string, deserializeOptions map[string]interface{}, serializeOptions map[string]interface{}) (string, error) {

	if len(inputString) == 0 {
		return "", ErrMissingInputString
	}

	if len(inputFormat) == 0 || inputFormat == "undefined" {
		return "", ErrMissingInputFormat
	}

	if len(outputFormat) == 0 || outputFormat == "undefined" {
		return "", ErrMissingOutputFormat
	}

	inputSerializer, err := serializer.NewWithOptions(inputFormat, DeserializeDefaults, deserializeOptions)
	if err != nil {
		return "", err
	}

	object, err := inputSerializer.Deserialize([]byte(inputString))
	if err != nil {
		return "", err
	}

	outputSerializer, err := serializer.NewWithOptions(outputFormat, SerializeDefaults, serializeOptions)
	if err != nil {
		return "", err
	}

	outputBytes, err := outputSerializer.Serialize(object)
	if err != nil {
		return "", err
	}

	return string(outputBytes), nil
}
