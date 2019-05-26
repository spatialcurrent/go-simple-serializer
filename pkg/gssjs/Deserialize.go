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
	"github.com/spatialcurrent/go-simple-serializer/pkg/gss"
	"honnef.co/go/js/console"
)

// Deserialize is a function provided to gss.js that wraps gss.Deserialize to support JavaScript.
func Deserialize(inputString string, inputFormat string, options *js.Object) interface{} {

	m := map[string]interface{}{}
	for _, key := range js.Keys(options) {
		m[key] = options.Get(key).Interface()
	}

	inputHeader := gss.NoHeader
	inputComment := gss.NoComment
	inputLazyQuotes := false
	inputSkipLines := gss.NoSkip
	inputLineSeparator := "\n"
	inputLimit := gss.NoLimit
	async := false

	if v, ok := m["header"]; ok {
		inputHeader = toStringSlice(v)
	}

	if v, ok := m["inputComment"]; ok {
		if vv, ok := v.(string); ok {
			inputComment = vv
		}
	}

	if v, ok := m["inputLineSeparator"]; ok {
		if vv, ok := v.(string); ok {
			inputLineSeparator = vv
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
		if vv, ok := v.(int); ok {
			inputLimit = vv
		}
	}

	if v, ok := m["async"]; ok {
		switch v := v.(type) {
		case bool:
			async = v
		}
	}

	inputType, err := gss.GetType([]byte(inputString), inputFormat)
	if err != nil {
		console.Error(errors.Wrap(err, "error creating new object for format "+inputFormat))
		return ""
	}

	outputObject, err := gss.DeserializeString(
		inputString,
		inputFormat,
		inputHeader,
		inputComment,
		inputLazyQuotes,
		inputSkipLines,
		inputLineSeparator,
		inputLimit,
		inputType,
		async,
		false)
	if err != nil {
		console.Error(errors.Wrap(err, "error deserializing input into object").Error())
		return ""
	}

	return outputObject
}
