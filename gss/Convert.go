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

// Convert converts an input_string from the inputFormat to the outputFormat.
// Returns the output string and error, if any.
func Convert(inputBytes []byte, inputFormat string, inputHeader []string, inputComment string, inputLazyQuotes bool, inputSkipLines int, inputLimit int, outputFormat string, outputHeader []string, outputLimit int, verbose bool) (string, error) {

	inputType, err := GetType(inputBytes, inputFormat)
	if err != nil {
		return "", errors.Wrap(err, "error creating new object for format "+inputFormat)
	}

	if verbose {
		fmt.Println("Input Format: " + inputFormat)
		fmt.Println("Output Format: " + outputFormat)
		fmt.Println("Input Type: " + fmt.Sprint(inputType))
	}

	if inputFormat == "bson" || inputFormat == "json" || inputFormat == "hcl" || inputFormat == "hcl2" || inputFormat == "properties" || inputFormat == "toml" || inputFormat == "yaml" {
		if outputFormat == "bson" || outputFormat == "json" || inputFormat == "hcl" || inputFormat == "hcl2" || outputFormat == "properties" || outputFormat == "toml" || outputFormat == "yaml" {
			object, err := DeserializeBytes(inputBytes, inputFormat, inputHeader, inputComment, inputLazyQuotes, inputSkipLines, inputLimit, inputType, verbose)
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			if verbose {
				fmt.Println("Object:", object)
			}
			output_string, err := SerializeString(object, outputFormat, outputHeader, outputLimit)
			if err != nil {
				return "", errors.Wrap(err, "Error serializing output")
			}
			return output_string, nil
		} else if outputFormat == "jsonl" {
			object, err := DeserializeBytes(inputBytes, inputFormat, inputHeader, inputComment, inputLazyQuotes, inputSkipLines, inputLimit, inputType, verbose)
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			output_string, err := SerializeString(object, outputFormat, outputHeader, outputLimit)
			if err != nil {
				return "", errors.Wrap(err, "Error serializing output")
			}
			return output_string, nil
		} else if outputFormat == "csv" {
			return "", errors.New("Error: incompatible output format \"" + outputFormat + "\"")
		} else if outputFormat == "tsv" {
			return "", errors.New("Error: incompatible output format \"" + outputFormat + "\"")
		} else {
			return "", errors.New("Error: unknown output format \"" + outputFormat + "\"")
		}
	} else if inputFormat == "jsonl" || inputFormat == "csv" || inputFormat == "tsv" {
		if outputFormat == "bson" || outputFormat == "json" || outputFormat == "hcl" || outputFormat == "hcl2" || outputFormat == "toml" || outputFormat == "yaml" || outputFormat == "jsonl" || outputFormat == "csv" || outputFormat == "tsv" {
			object, err := DeserializeBytes(inputBytes, inputFormat, inputHeader, inputComment, inputLazyQuotes, inputSkipLines, inputLimit, inputType, verbose)
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			output_string, err := SerializeString(object, outputFormat, outputHeader, outputLimit)
			if err != nil {
				return "", errors.Wrap(err, "Error serializing output")
			}
			return output_string, nil
		} else {
			return "", errors.New("Error: unknown output format \"" + inputFormat + "\"")
		}
	}
	return "", errors.New("Error: unknown input format \"" + inputFormat + "\"")
}
