// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"reflect"
)

import (
	"github.com/pkg/errors"
)

// ToRowI converts a map to a row of strings using a slice of interface{} keys.
func ToRowI(keys []interface{}, m reflect.Value, valueSerializer func(object interface{}) (string, error)) ([]string, error) {
	row := make([]string, len(keys))
	for j, key := range keys {
		if v := m.MapIndex(reflect.ValueOf(key)); v.IsValid() && (v.Type().Kind() == reflect.String || !v.IsNil()) {
			str, err := valueSerializer(v.Interface())
			if err != nil {
				return row, errors.Wrap(err, "error serializing value")
			}
			row[j] = str
		} else {
			str, err := valueSerializer(nil)
			if err != nil {
				return row, errors.Wrap(err, "error serializing value")
			}
			row[j] = str
		}
	}
	return row, nil
}
