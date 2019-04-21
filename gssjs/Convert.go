// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gssjs

import (
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
		convertInput.InputHeader = toStringSlice(v)
	}

	if v, ok := m["inputComment"]; ok {
		if vv, ok := v.(string); ok {
			convertInput.InputComment = vv
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
		if vv, ok := v.(int); ok {
			convertInput.InputLimit = vv
		}
	}

	if v, ok := m["outputHeader"]; ok {
		convertInput.OutputHeader = toStringSlice(v)
	}

	if v, ok := m["outputLimit"]; ok {
		if vv, ok := v.(int); ok {
			convertInput.OutputLimit = vv
		}
	}

	if v, ok := m["async"]; ok {
		switch vv := v.(type) {
		case bool:
			convertInput.Async = vv
		case int:
			convertInput.Async = vv > 0
		}
	}

	if v, ok := m["outputPretty"]; ok {
		switch vv := v.(type) {
		case bool:
			convertInput.OutputPretty = vv
		case int:
			convertInput.OutputPretty = vv > 0
		}
	}

	outputString, err := gss.Convert(convertInput)
	if err != nil {
		console.Error(errors.Wrap(err, "error converting input").Error())
		return ""
	}

	return outputString
}
