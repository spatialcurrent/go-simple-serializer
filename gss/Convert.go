// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"fmt"
	"github.com/pkg/errors"
)

// Convert converts an input_string from the input_format to the output_format.
// Returns the output string and error, if any.
func Convert(input_bytes []byte, input_format string, input_header []string, input_comment string, input_limit int, output_format string, output_limit int, verbose bool) (string, error) {

	input_type, err := GetType(input_bytes, input_format)
	if err != nil {
		return "", errors.Wrap(err, "error creating new object for format "+input_format)
	}

	if verbose {
		fmt.Println("Input Format: " + input_format)
		fmt.Println("Output Format: " + output_format)
		fmt.Println("Input Type: " + fmt.Sprint(input_type))
	}

	if input_format == "bson" || input_format == "json" || input_format == "hcl" || input_format == "hcl2" || input_format == "properties" || input_format == "toml" || input_format == "yaml" {
		if output_format == "bson" || output_format == "json" || input_format == "hcl" || input_format == "hcl2" || output_format == "properties" || output_format == "toml" || output_format == "yaml" {
			object, err := Deserialize(string(input_bytes), input_format, input_header, input_comment, input_limit, input_type, verbose)
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			if verbose {
				fmt.Println("Object:", object)
			}
			output_string, err := Serialize(object, output_format)
			if err != nil {
				return "", errors.Wrap(err, "Error serializing output")
			}
			return output_string, nil
		} else if output_format == "jsonl" {
			object, err := Deserialize(string(input_bytes), input_format, input_header, input_comment, input_limit, input_type, verbose)
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
			object, err := Deserialize(string(input_bytes), input_format, input_header, input_comment, input_limit, input_type, verbose)
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
