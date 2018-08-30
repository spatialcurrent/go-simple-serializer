// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// gss.js is the Javascript package for go-simple-serializer (GSS).
//
// Usage
//
// In you html document, the simplest workflow is to add GSS as a script and call gss.Convert(input_string, input_format, output_format);
//
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/spatialcurrent/go-simple-serializer/gss"
	"github.com/spatialcurrent/go-simple-serializer/gssjs"
)

func main() {
	js.Global.Set("gss", map[string]interface{}{
		"version":     gss.VERSION,
		"formats":     gss.Formats,
		"convert":     gssjs.Convert,
		"deserialize": gssjs.Deserialize,
		"serialize":   gssjs.Serialize,
	})
}
