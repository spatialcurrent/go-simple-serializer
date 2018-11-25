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
func DeserializeJSONL(input string, inputComment string, inputSkipLines int, inputLimit int, outputType reflect.Type) (interface{}, error) {
	output := reflect.MakeSlice(outputType, 0, 0)
	if inputLimit == 0 {
		return output.Interface(), nil
	}
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanLines)
	for i := 0; i < inputSkipLines; i++ {
		if !scanner.Scan() {
			break
		}
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(inputComment) == 0 || !strings.HasPrefix(line, inputComment) {
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
			if inputLimit > 0 && output.Len() >= inputLimit {
				break
			}
		}
	}

	return output.Interface(), nil
}
