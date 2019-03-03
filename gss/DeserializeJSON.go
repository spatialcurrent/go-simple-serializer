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
func DeserializeJSON(input_bytes []byte, output_type reflect.Type) (interface{}, error) {
	if output_type.Kind() == reflect.Map {
		ptr := reflect.New(output_type)
		ptr.Elem().Set(reflect.MakeMap(output_type))
		err := json.Unmarshal(input_bytes, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling bytes into JSON")
		}
		return ptr.Elem().Interface(), nil
	} else if output_type.Kind() == reflect.Slice {
		ptr := reflect.New(output_type)
		ptr.Elem().Set(reflect.MakeSlice(output_type, 0, 0))
		err := json.Unmarshal(input_bytes, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshalling bytes into JSON")
		}
		return ptr.Elem().Interface(), nil
	}
	return nil, errors.New("Invalid output type for json " + fmt.Sprint(output_type))
}
