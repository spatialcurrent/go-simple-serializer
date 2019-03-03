// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"reflect"
)

// DeserializeYAML deserializes the YAML input bytes into a Go object
func DeserializeYAML(input_bytes []byte, output_type reflect.Type) (interface{}, error) {
	if output_type.Kind() == reflect.Map {
		ptr := reflect.New(output_type)
		ptr.Elem().Set(reflect.MakeMap(output_type))
		err := yaml.Unmarshal(input_bytes, ptr.Interface())
		return ptr.Elem().Interface(), err
	} else if output_type.Kind() == reflect.Slice {
		ptr := reflect.New(output_type)
		ptr.Elem().Set(reflect.MakeSlice(output_type, 0, 0))
		err := yaml.Unmarshal(input_bytes, ptr.Interface())
		return StringifyMapKeys(ptr.Elem().Interface()), err
	}
	return nil, errors.New("Invalid output type for yaml " + fmt.Sprint(output_type))
}
