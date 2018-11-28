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

	inputHeader := gss.NoHeader
	inputComment := gss.NoComment
	inputLazyQuotes := false
	inputSkipLines := gss.NoSkip
	inputLimit := gss.NoLimit
	async := false

	if v, ok := m["header"]; ok {
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

	if v, ok := m["async"]; ok {
		switch v := v.(type) {
		case bool:
			async = v
		}
	}

	input_type, err := gss.GetType([]byte(input_string), input_format)
	if err != nil {
		console.Error(errors.Wrap(err, "error creating new object for format "+input_format))
		return ""
	}

	outputObject, err := gss.DeserializeString(input_string, input_format, inputHeader, inputComment, inputLazyQuotes, inputSkipLines, inputLimit, input_type, async, false)
	if err != nil {
		console.Error(errors.Wrap(err, "error deserializing input into object").Error())
		return ""
	}

	return outputObject
}
