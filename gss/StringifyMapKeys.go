// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"fmt"
)

// StringifyMapKeys recursively converts map keys from interface{} to string.
// This functionality is inspired by https://github.com/gohugoio/hugo.
// See https://github.com/gohugoio/hugo/pull/4138 for more info.
func StringifyMapKeys(in interface{}) interface{} {
	switch in := in.(type) {
	case []interface{}:
		res := make([]interface{}, len(in))
		for i, v := range in {
			res[i] = StringifyMapKeys(v)
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
