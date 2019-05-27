// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"reflect"
)

// ToRow converts an object into a row of strings and returns an error, if any.
func ToRow(object interface{}, columns []interface{}, valueSerializer func(object interface{}) (string, error)) ([]string, error) {
	return ToRowFromValue(reflect.ValueOf(object), columns, valueSerializer)
}
