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

// Convert is a function provided to gss.js that wraps gss.Convert to support JavaScript.
func Convert(inputString string, inputFormat string, outputFormat string, options *js.Object) string {

	convertInput := gss.NewConvertInput([]byte(inputString), inputFormat, outputFormat)

	m := map[string]interface{}{}
	for _, key := range js.Keys(options) {
		m[key] = options.Get(key).Interface()
	}

	if v, ok := m["inputHeader"]; ok {
		inputHeader := make([]string, 0)
		switch v.(type) {
		case []string:
			inputHeader = v.([]string)
		case []interface{}:
			inputHeader = make([]string, 0, len(v.([]interface{})))
			for _, h := range v.([]interface{}) {
				inputHeader = append(inputHeader, fmt.Sprint(h))
			}
		}
		convertInput.InputHeader = inputHeader
	}

	if v, ok := m["inputComment"]; ok {
		switch v := v.(type) {
		case string:
			convertInput.InputComment = v
		}
	}

	if v, ok := m["inputLazyQuotes"]; ok {
		switch v := v.(type) {
		case bool:
			convertInput.InputLazyQuotes = v
		case int:
			convertInput.InputLazyQuotes = v > 0
		}
	}

	if v, ok := m["inputLimit"]; ok {
		switch v := v.(type) {
		case int:
			convertInput.InputLimit = v
		}
	}

	if v, ok := m["outputHeader"]; ok {
		outputHeader := make([]string, 0)
		switch v.(type) {
		case []string:
			outputHeader = v.([]string)
		case []interface{}:
			outputHeader = make([]string, 0, len(v.([]interface{})))
			for _, h := range v.([]interface{}) {
				outputHeader = append(outputHeader, fmt.Sprint(h))
			}
		}
		convertInput.OutputHeader = outputHeader
	}

	if v, ok := m["outputLimit"]; ok {
		switch v := v.(type) {
		case int:
			convertInput.OutputLimit = v
		}
	}

	if v, ok := m["async"]; ok {
		switch v := v.(type) {
		case bool:
			convertInput.Async = v
		}
	}

	if v, ok := m["outputPretty"]; ok {
		switch v := v.(type) {
		case bool:
			convertInput.OutputPretty = v
		}
	}

	outputString, err := gss.Convert(convertInput)
	if err != nil {
		console.Error(errors.Wrap(err, "error converting input").Error())
		return ""
	}

	return outputString
}
