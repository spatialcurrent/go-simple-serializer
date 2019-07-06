// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// GetType takes in the content of an object as a string and the serialization format.
// Returns the type using reflection.
// This type is fixed and can be passed through functions without losing type information (unlike an empty object).
func GetType(content []byte, format string) (reflect.Type, error) {

	if format == "json" || format == "yaml" {
		str := string(content)
		if format == "json" {
			return GetTypeJSON(strings.TrimLeftFunc(str, unicode.IsSpace)), nil
		} else if format == "yaml" {
			str = strings.TrimLeftFunc(str, unicode.IsSpace)
			if len(str) > 0 && str[0] == '-' {
				return reflect.TypeOf([]interface{}{}), nil
			}
			return reflect.TypeOf(map[string]interface{}{}), nil
		}
	} else if format == "bson" || format == "hcl" || format == "hcl2" || format == "properties" || format == "toml" {
		return reflect.TypeOf(map[string]interface{}{}), nil
	} else if format == "csv" || format == "tsv" {
		return reflect.TypeOf([]map[string]interface{}{}), nil
	} else if format == "jsonl" {
		return reflect.TypeOf([]interface{}{}), nil
	}

	return nil, errors.New("could not get type for format " + format)
}
