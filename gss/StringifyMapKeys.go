// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"fmt"
)

// StringifyMapKeys recursively converts map keys from interface{} to string.
// This functionality is inspired by work done in https://github.com/gohugoio/hugo, but
// support many more types, including:
//	- []interface{}
//	- [][]interface{}
//	- []map[interface{}]interface{}
//	- map[interface{}]interface{}
//	- map[string]interface{}
//	- map[string][]interface{}
//	- map[string]map[string][]interface{}
//	- map[interface{}]struct{}
// See https://github.com/gohugoio/hugo/pull/4138 for background.
func StringifyMapKeys(in interface{}) interface{} {

	switch in := in.(type) {
	case []interface{}:
		res := make([]interface{}, len(in))
		for i, v := range in {
			res[i] = StringifyMapKeys(v)
		}
		return res
	case [][]interface{}:
		res := make([][]interface{}, len(in))
		for i, v := range in {
			res[i] = StringifyMapKeys(v).([]interface{})
		}
		return res
	case []map[interface{}]interface{}:
		res := make([]map[string]interface{}, len(in))
		for i, v := range in {
			res[i] = StringifyMapKeys(v).(map[string]interface{})
		}
		return res
	case map[interface{}]interface{}:
		res := make(map[string]interface{})
		for k, v := range in {
			res[fmt.Sprintf("%v", k)] = StringifyMapKeys(v)
		}
		return res
	case map[string]interface{}:
		res := make(map[string]interface{})
		for k, v := range in {
			res[k] = StringifyMapKeys(v)
		}
		return res
	case map[string][]interface{}:
		res := make(map[string][]interface{})
		for k, v := range in {
			res[k] = StringifyMapKeys(v).([]interface{})
		}
		return res
	case map[string]map[string]interface{}:
		res := make(map[string]map[string]interface{})
		for k, v := range in {
			res[k] = StringifyMapKeys(v).(map[string]interface{})
		}
		return res
	case map[string]map[string][]interface{}:
		res := make(map[string]map[string][]interface{})
		for k, v := range in {
			res[k] = StringifyMapKeys(v).(map[string][]interface{})
		}
		return res
	case map[interface{}]struct{}:
		res := make(map[string]interface{})
		for k, v := range in {
			res[fmt.Sprintf("%v", k)] = StringifyMapKeys(v)
		}
		return res
	default:
		return in
	}
}
