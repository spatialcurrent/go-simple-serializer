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
	"fmt"
	"github.com/spatialcurrent/go-simple-serializer/gss"
)

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
)

func main() {
	js.Global.Set("gss", map[string]interface{}{
		"version": gss.VERSION,
		"convert": Convert,
	})
}

func Convert(input_string string, input_format string, output_format string, options *js.Object) string {

	m := map[string]interface{}{}
	for _, key := range js.Keys(options) {
		m[key] = options.Get(key).Interface()
	}

	input_header := []string{}
	input_comment := ""

	if v, ok := m["header"]; ok {
		switch v.(type) {
		case []string:
			input_header = v.([]string)
		case []interface{}:
			input_header = make([]string, 0, len(v.([]interface{})))
			for _, h := range v.([]interface{}) {
				input_header = append(input_header, fmt.Sprint(h))
			}
		}
	}

	if v, ok := m["input_comment"]; ok {
		switch v := v.(type) {
		case string:
			input_comment = v
		}
	}

	output_string, err := gss.Convert(input_string, input_format, input_header, input_comment, output_format)
	if err != nil {
		console.Log(err.Error())
		return ""
	}
	return output_string
}
