// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"fmt"
	"reflect"
	"strings"
)

// ToRowFromValue converts an object into a row of strings and returns an error, if any.
func ToRowFromValue(objectValue reflect.Value, columns []interface{}, valueSerializer func(object interface{}) (string, error)) ([]string, error) {
	for reflect.TypeOf(objectValue.Interface()).Kind() == reflect.Ptr {
		objectValue = objectValue.Elem()
	}
	objectValue = reflect.ValueOf(objectValue.Interface()) // sets value to concerete type
	row := make([]string, len(columns))
	switch objectValue.Type().Kind() {
	case reflect.Map:
		for j, key := range columns {
			if v := objectValue.MapIndex(reflect.ValueOf(key)); v.IsValid() && (v.Type().Kind() == reflect.String || !v.IsNil()) {
				str, err := valueSerializer(v.Interface())
				if err != nil {
					return row, fmt.Errorf("error serializing value: %w", err)
				}
				row[j] = str
			} else {
				str, err := valueSerializer(nil)
				if err != nil {
					return row, fmt.Errorf("error serializing value: %w", err)
				}
				row[j] = str
			}
		}
	case reflect.Struct:
		for j, column := range columns {
			columnLowerCase := strings.ToLower(fmt.Sprint(column))
			if f := objectValue.FieldByNameFunc(func(match string) bool { return strings.ToLower(match) == columnLowerCase }); f.IsValid() && (f.Type().Kind() == reflect.String || !f.IsNil()) {
				str, err := valueSerializer(f.Interface())
				if err != nil {
					return row, fmt.Errorf("error serializing value: %w", err)
				}
				row[j] = str
			} else {
				str, err := valueSerializer(nil)
				if err != nil {
					return row, fmt.Errorf("error serializing value: %w", err)
				}
				row[j] = str
			}
		}
	}

	return row, nil
}
