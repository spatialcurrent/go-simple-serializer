// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gssjs

import (
	"fmt"
	"reflect"
)

func toStringSlice(v interface{}) []string {
	if stringSlice, ok := v.([]string); ok {
		return stringSlice
	}
	slc := make([]string, 0)
	if k := reflect.TypeOf(v).Kind(); k == reflect.Array || k == reflect.Slice {
		vv := reflect.ValueOf(v)
		slc = make([]string, 0, vv.Len())
		for i := 0; i < vv.Len(); i++ {
			slc = append(slc, fmt.Sprint(vv.Index(i).Interface()))
		}
	}
	return slc
}
