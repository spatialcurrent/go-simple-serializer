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

// Serialize is a function provided to gss.js that wraps gss.Serialize to support JavaScript.
func Serialize(input_object interface{}, output_format string, options *js.Object) interface{} {

	m := map[string]interface{}{}
	for _, key := range js.Keys(options) {
		m[key] = options.Get(key).Interface()
	}

	output_header := []string{}
	output_limit := -1

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

	output_string, err := gss.SerializeString(input_object, output_format, output_header, output_limit)
	if err != nil {
		console.Error(errors.Wrap(err, "error serializing object").Error())
		return ""
	}

	return output_string
}
