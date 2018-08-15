// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"strings"
)

import (
	"github.com/pkg/errors"
)

// Convert converts an input_string from the input_format to the output_format and returns an error, if any.
func Convert(input_string string, input_format string, input_header []string, input_comment string, output_format string) (string, error) {

	var object interface{}

	if input_format == "json" {
		input_string = strings.TrimSpace(input_string)
		if len(input_string) > 0 && input_string[0] == '[' {
			object = []map[string]interface{}{}
		} else {
			object = map[string]interface{}{}
		}
	} else if input_format == "yaml" {
		input_string = strings.TrimSpace(input_string)
		if len(input_string) > 0 && input_string[0] == '-' {
			object = []map[string]interface{}{}
		} else {
			object = map[string]interface{}{}
		}
	} else if input_format == "bson" || input_format == "hcl" || input_format == "hcl2" || input_format == "properties" || input_format == "toml" {
		object = map[string]interface{}{}
	}

	if input_format == "bson" || input_format == "json" || input_format == "hcl" || input_format == "hcl2" || input_format == "properties" || input_format == "toml" || input_format == "yaml" {
		if output_format == "bson" || output_format == "json" || input_format == "hcl" || input_format == "hcl2" || output_format == "properties" || output_format == "toml" || output_format == "yaml" {
			err := Deserialize(input_string, input_format, input_header, input_comment, &object)
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			output_string, err := Serialize(object, output_format)
			if err != nil {
				return "", errors.Wrap(err, "Error serializing output")
			}
			return output_string, nil
		} else if output_format == "jsonl" {
			object := []map[string]interface{}{}
			err := Deserialize(input_string, input_format, input_header, input_comment, &object)
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			output_string, err := Serialize(object, output_format)
			if err != nil {
				return "", errors.Wrap(err, "Error serializing output")
			}
			return output_string, nil
		} else if output_format == "csv" {
			return "", errors.New("Error: incompatible output format \"" + output_format + "\"")
		} else if output_format == "tsv" {
			return "", errors.New("Error: incompatible output format \"" + output_format + "\"")
		} else {
			return "", errors.New("Error: unknown output format \"" + output_format + "\"")
		}
	} else if input_format == "jsonl" || input_format == "csv" || input_format == "tsv" {
		if output_format == "bson" || output_format == "json" || output_format == "hcl" || output_format == "hcl2" || output_format == "toml" || output_format == "yaml" || output_format == "jsonl" || output_format == "csv" || output_format == "tsv" {
			object := []map[string]interface{}{}
			err := Deserialize(input_string, input_format, input_header, input_comment, &object)
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			output_string, err := Serialize(object, output_format)
			if err != nil {
				return "", errors.Wrap(err, "Error serializing output")
			}
			return output_string, nil
		} else {
			return "", errors.New("Error: unknown output format \"" + input_format + "\"")
		}
	}
	return "", errors.New("Error: unknown input format \"" + input_format + "\"")
}
