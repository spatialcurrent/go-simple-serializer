// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"fmt"
	"reflect"
)

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// DeserializeYAML deserializes the YAML input bytes into a Go object
func DeserializeYAML(inputBytes []byte, outputType reflect.Type) (interface{}, error) {
	if outputType.Kind() == reflect.Map {
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeMap(outputType))
		err := yaml.Unmarshal(inputBytes, ptr.Interface())
		return ptr.Elem().Interface(), err
	} else if outputType.Kind() == reflect.Slice {
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeSlice(outputType, 0, 0))
		err := yaml.Unmarshal(inputBytes, ptr.Interface())
		return stringify.StringifyMapKeys(ptr.Elem().Interface()), err
	}
	return nil, errors.New("Invalid output type for yaml " + fmt.Sprint(outputType))
}
