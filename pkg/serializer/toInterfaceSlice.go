// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package serializer

import (
	"fmt"
	"reflect"
)

func toInterfaceSlice(v interface{}) []interface{} {
	if interfaceSlice, ok := v.([]interface{}); ok {
		return interfaceSlice
	}
	slc := make([]interface{}, 0)
	if k := reflect.TypeOf(v).Kind(); k == reflect.Array || k == reflect.Slice {
		vv := reflect.ValueOf(v)
		slc = make([]interface{}, 0, vv.Len())
		for i := 0; i < vv.Len(); i++ {
			slc = append(slc, fmt.Sprint(vv.Index(i).Interface()))
		}
	}
	return slc
}
