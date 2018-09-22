// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"bufio"
	"encoding/json"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

// DeserializeJSONL deserializes the input JSON lines bytes into a Go object.
//  - https://golang.org/pkg/encoding/json/
func DeserializeJSONL(input string, input_comment string, input_limit int, output_type reflect.Type) (interface{}, error) {
	output := reflect.MakeSlice(output_type, 0, 0)
	if input_limit == 0 {
		return output.Interface(), nil
	}
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(input_comment) == 0 || !strings.HasPrefix(line, input_comment) {
			lineType, err := GetType([]byte(line), "json")
			if err != nil {
				return nil, errors.Wrap(err, "error getting type for input line")
			}
			var ptr reflect.Value
			if lineType.Kind() == reflect.Array || lineType.Kind() == reflect.Slice {
				ptr = reflect.New(lineType)
				ptr.Elem().Set(reflect.MakeSlice(lineType, 0, 0))
			} else if lineType.Kind() == reflect.Map {
				ptr = reflect.New(lineType)
				ptr.Elem().Set(reflect.MakeMap(lineType))
			} else {
				return nil, errors.Wrap(err, "error creating object for JSON line "+line)
			}
			err = json.Unmarshal([]byte(line), ptr.Interface())
			if err != nil {
				return nil, errors.Wrap(err, "Error reading object from JSON line")
			}
			output = reflect.Append(output, ptr.Elem())
			if input_limit > 0 && output.Len() >= input_limit {
				break
			}
		}
	}

	return output.Interface(), nil
}
