// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
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

	outputHeader := gss.NoHeader
	outputLimit := gss.NoLimit

	if v, ok := m["outputHeader"]; ok {
		switch v.(type) {
		case []string:
			outputHeader = v.([]string)
		case []interface{}:
			outputHeader = make([]string, 0, len(v.([]interface{})))
			for _, h := range v.([]interface{}) {
				outputHeader = append(outputHeader, fmt.Sprint(h))
			}
		}
	}

	if v, ok := m["outputLimit"]; ok {
		switch v := v.(type) {
		case int:
			outputLimit = v
		}
	}

	output_string, err := gss.SerializeString(input_object, output_format, outputHeader, outputLimit)
	if err != nil {
		console.Error(errors.Wrap(err, "error serializing object").Error())
		return ""
	}

	return output_string
}
