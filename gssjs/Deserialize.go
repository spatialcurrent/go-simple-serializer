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

// Deserialize is a function provided to gss.js that wraps gss.Deserialize to support JavaScript.
func Deserialize(input_string string, input_format string, options *js.Object) interface{} {

	m := map[string]interface{}{}
	for _, key := range js.Keys(options) {
		m[key] = options.Get(key).Interface()
	}

	input_header := []string{}
	input_comment := ""
	input_limit := -1

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

	if v, ok := m["input_limit"]; ok {
		switch v := v.(type) {
		case int:
			input_limit = v
		}
	}

	input_type, err := gss.GetType([]byte(input_string), input_format)
	if err != nil {
		console.Error(errors.Wrap(err, "error creating new object for format "+input_format))
		return ""
	}

	output_object, err := gss.Deserialize(input_string, input_format, input_header, input_comment, input_limit, input_type, false)
	if err != nil {
		console.Error(errors.Wrap(err, "error deserializing input into object").Error())
		return ""
	}

	return output_object
}
