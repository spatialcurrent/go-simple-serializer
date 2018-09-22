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
	input_lazy_quotes := false
	input_limit := -1
	output_header := []string{}
	output_limit := -1

	if v, ok := m["input_header"]; ok {
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

	if v, ok := m["input_lazy_quotes"]; ok {
		switch v := v.(type) {
		case bool:
			input_lazy_quotes = v
		case int:
			input_lazy_quotes = v > 0
		}
	}

	if v, ok := m["input_limit"]; ok {
		switch v := v.(type) {
		case int:
			input_limit = v
		}
	}

	if v, ok := m["output_header"]; ok {
		switch v.(type) {
		case []string:
			output_header = v.([]string)
		case []interface{}:
			output_header = make([]string, 0, len(v.([]interface{})))
			for _, h := range v.([]interface{}) {
				output_header = append(output_header, fmt.Sprint(h))
			}
		}
	}

	if v, ok := m["output_limit"]; ok {
		switch v := v.(type) {
		case int:
			output_limit = v
		}
	}

	output_string, err := gss.Convert([]byte(input_string), input_format, input_header, input_comment, input_lazy_quotes, input_limit, output_format, output_header, output_limit, false)
	if err != nil {
		console.Error(errors.Wrap(err, "error converting input").Error())
		return ""
	}

	return output_string
}
