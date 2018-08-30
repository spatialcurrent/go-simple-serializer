// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gssjs

import (
	"github.com/pkg/errors"
	"github.com/spatialcurrent/go-simple-serializer/gss"
	"honnef.co/go/js/console"
)

// Serialize is a function provided to gss.js that wraps gss.Serialize to support JavaScript.
func Serialize(input_object interface{}, output_format string) interface{} {

	output_string, err := gss.Serialize(input_object, output_format)
	if err != nil {
		console.Error(errors.Wrap(err, "error serializing object").Error())
		return ""
	}

	return output_string
}
