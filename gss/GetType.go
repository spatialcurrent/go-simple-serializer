// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"unicode"
)

// GetType takes in the content of an object as a string and the serialization format.
// Returns the type using reflection.
// This type is fixed and can be passed through functions without losing type information (unlike an empty object).
func GetType(content string, format string) (reflect.Type, error) {

	if format == "json" {
		content = strings.TrimLeftFunc(content, unicode.IsSpace)
		if len(content) > 0 && content[0] == '[' {
			return reflect.TypeOf([]map[string]interface{}{}), nil
		}
		return reflect.TypeOf(map[string]interface{}{}), nil
	} else if format == "yaml" {
		content = strings.TrimLeftFunc(content, unicode.IsSpace)
		if len(content) > 0 && content[0] == '-' {
			return reflect.TypeOf([]map[string]interface{}{}), nil
		}
		return reflect.TypeOf(map[string]interface{}{}), nil
	} else if format == "bson" || format == "hcl" || format == "hcl2" || format == "properties" || format == "toml" {
		return reflect.TypeOf(map[string]interface{}{}), nil
	} else if format == "jsonl" || format == "csv" || format == "tsv" {
		return reflect.TypeOf([]map[string]interface{}{}), nil
	}

	return nil, errors.New("could not get type for format " + format)
}
