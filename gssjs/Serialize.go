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
func Serialize(inputObject interface{}, outputFormat string, options *js.Object) interface{} {

	m := map[string]interface{}{}
	for _, key := range js.Keys(options) {
		m[key] = options.Get(key).Interface()
	}

	outputHeader := gss.NoHeader
	outputLimit := gss.NoLimit
	outputPretty := false

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

	if v, ok := m["outputPretty"]; ok {
		switch v := v.(type) {
		case bool:
			outputPretty = v
		}
	}

	outputString, err := gss.SerializeString(&gss.SerializeInput{
		Object: inputObject,
		Format: outputFormat,
		Header: outputHeader,
		Limit:  outputLimit,
		Pretty: outputPretty,
	})
	if err != nil {
		console.Error(errors.Wrap(err, "error serializing object").Error())
		return ""
	}

	return outputString
}
