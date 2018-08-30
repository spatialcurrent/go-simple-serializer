// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gssjs

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/pkg/errors"
	"github.com/spatialcurrent/go-simple-serializer/gss"
	"honnef.co/go/js/console"
)

// Convert is a function provided to gss.js that wraps gss.Convert to support JavaScript.
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

	output_string, err := gss.Convert(input_string, input_format, input_header, input_comment, output_format, false)
	if err != nil {
		console.Error(errors.Wrap(err, "error converting input").Error())
		return ""
	}
	
	return output_string
}