// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

// DeserializeJSON deserializes the input bytes into a Go object.
//  - https://golang.org/pkg/encoding/json/
func DeserializeJSON(inputBytes []byte, outputType reflect.Type) (interface{}, error) {
	if outputType.Kind() == reflect.Map {
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeMap(outputType))
		err := json.Unmarshal(inputBytes, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling bytes into JSON")
		}
		return ptr.Elem().Interface(), nil
	} else if outputType.Kind() == reflect.Slice {
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeSlice(outputType, 0, 0))
		err := json.Unmarshal(inputBytes, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling bytes into JSON")
		}
		return ptr.Elem().Interface(), nil
	}
	return nil, errors.New("Invalid output type for json " + fmt.Sprint(outputType))
}
