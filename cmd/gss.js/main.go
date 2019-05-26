// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
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
	"github.com/spatialcurrent/go-simple-serializer/pkg/gss"
	"github.com/spatialcurrent/go-simple-serializer/pkg/gssjs"
)

var gitBranch string
var gitCommit string

func main() {
	version := map[string]interface{}{}
	if len(gitBranch) > 0 {
		version["branch"] = gitBranch
	}
	if len(gitCommit) > 0 {
		version["commit"] = gitCommit
	}
	js.Global.Set("gss", map[string]interface{}{
		"version":     version,
		"formats":     gss.Formats,
		"convert":     gssjs.Convert,
		"deserialize": gssjs.Deserialize,
		"serialize":   gssjs.Serialize,
	})
}
