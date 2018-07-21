// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// GSS.JS is the Javascript version of GSS.
//
// Usage
//
// In you html document, the simplest workflow is to add GSS as a script and call gss.Convert(input_string, input_format, output_format);
package main

import (
	"github.com/spatialcurrent/go-simple-serializer/gss"
)

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
)

var GO_GSS_VERSION = "0.0.1"

func main() {
	js.Global.Set("gss", map[string]interface{}{
		"version":        GO_GSS_VERSION,
		"convert": Convert,
	})
}

func Convert(input_string string, input_format string, output_format string) string {
	output_string, err := gss.Convert(input_string, input_format, output_format)
	if err != nil {
		console.Log(err.Error())
		return ""
	}
	return output_string
}
