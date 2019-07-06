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
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

func isArray(object *js.Object) bool {
	return js.Global.Get("Array").Call("isArray", object).Bool()
}

func toArray(object *js.Object) []interface{} {
	arr := make([]interface{}, 0, object.Length())
	for i := 0; i < object.Length(); i++ {
		arr = append(arr, parseObject(object.Index(i)))
	}
	return arr
}

func parseObject(object *js.Object) interface{} {
	if isArray(object) {
		return toArray(object)
	}
	return object.Interface()
}

var Exports = map[string]interface{}{
	"formats": serializer.Formats,
	"convert": func(inputString string, inputFormat string, outputFormat string, inputOptions map[string]interface{}, outputOptions map[string]interface{}) map[string]interface{} {
		str, err := Convert(inputString, inputFormat, outputFormat, inputOptions, outputOptions)
		if err != nil {
			return map[string]interface{}{"str": nil, "err": errors.Wrap(err, "error converting input string").Error()}
		}
		return map[string]interface{}{"str": str, "err": nil}
	},
	"deserialize": func(inputString string, inputFormat string, options map[string]interface{}) map[string]interface{} {
		obj, err := Deserialize(inputString, inputFormat, options)
		if err != nil {
			return map[string]interface{}{"obj": nil, "err": errors.Wrap(err, "error deserializing input string").Error()}
		}
		return map[string]interface{}{"obj": obj, "err": nil}
	},
	"serialize": func(inputObject *js.Object, outputFormat string, options map[string]interface{}) map[string]interface{} {
		str, err := Serialize(parseObject(inputObject), outputFormat, options)
		if err != nil {
			return map[string]interface{}{"str": nil, "err": errors.Wrap(err, "error serializing input object").Error()}
		}
		return map[string]interface{}{"str": str, "err": nil}
	},
}
