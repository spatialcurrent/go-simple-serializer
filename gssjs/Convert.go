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
func Convert(inputString string, inputFormat string, outputFormat string, options *js.Object) string {

	m := map[string]interface{}{}
	for _, key := range js.Keys(options) {
		m[key] = options.Get(key).Interface()
	}

	inputHeader := gss.NoHeader
	inputComment := gss.NoComment
	inputLazyQuotes := false
	inputSkipLines := gss.NoSkip
	inputLimit := gss.NoLimit
	outputHeader := gss.NoHeader
	outputLimit := gss.NoLimit

	if v, ok := m["inputHeader"]; ok {
		switch v.(type) {
		case []string:
			inputHeader = v.([]string)
		case []interface{}:
			inputHeader = make([]string, 0, len(v.([]interface{})))
			for _, h := range v.([]interface{}) {
				inputHeader = append(inputHeader, fmt.Sprint(h))
			}
		}
	}

	if v, ok := m["inputComment"]; ok {
		switch v := v.(type) {
		case string:
			inputComment = v
		}
	}

	if v, ok := m["inputLazyQuotes"]; ok {
		switch v := v.(type) {
		case bool:
			inputLazyQuotes = v
		case int:
			inputLazyQuotes = v > 0
		}
	}

	if v, ok := m["inputLimit"]; ok {
		switch v := v.(type) {
		case int:
			inputLimit = v
		}
	}

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

	output_string, err := gss.Convert(
		[]byte(inputString),
		inputFormat,
		inputHeader,
		inputComment,
		inputLazyQuotes,
		inputSkipLines,
		inputLimit,
		outputFormat,
		outputHeader,
		outputLimit,
		false)
	if err != nil {
		console.Error(errors.Wrap(err, "error converting input").Error())
		return ""
	}

	return output_string
}
