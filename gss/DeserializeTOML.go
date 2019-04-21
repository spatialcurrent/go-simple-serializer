// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"reflect"
)

// DeserializeTOML deserializes an input TOML string into a Go object
//  - https://godoc.org/pkg/github.com/BurntSushi/toml
func DeserializeTOML(input string, outputType reflect.Type) (interface{}, error) {
	if outputType.Kind() == reflect.Map {
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeMap(outputType))
		_, err := toml.Decode(input, ptr.Interface())
		return ptr.Elem().Interface(), err
	}
	return nil, errors.New("Invalid output type for toml " + fmt.Sprint(outputType))
}
