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
func Convert(inputBytes []byte, inputFormat string, inputHeader []string, inputComment string, inputLazyQuotes bool, inputSkipLines int, inputLimit int, outputFormat string, outputHeader []string, outputLimit int, async bool, verbose bool) (string, error) {

	inputType, err := GetType(inputBytes, inputFormat)
	if err != nil {
		return "", errors.Wrap(err, "error creating new object for format "+inputFormat)
	}

	if verbose {
		fmt.Println("Input Format: " + inputFormat)
		fmt.Println("Output Format: " + outputFormat)
		fmt.Println("Input Type: " + fmt.Sprint(inputType))
	}

	switch inputFormat {
	case "bson", "json", "hcl", "hcl2", "properties", "toml", "yaml":
		switch outputFormat {
		case "bson", "json", "hcl", "hcl2", "properties", "toml", "yaml":
			object, err := DeserializeBytes(inputBytes, inputFormat, inputHeader, inputComment, inputLazyQuotes, inputSkipLines, inputLimit, inputType, async, verbose)
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
		case "jsonl":
			object, err := DeserializeBytes(inputBytes, inputFormat, inputHeader, inputComment, inputLazyQuotes, inputSkipLines, inputLimit, inputType, async, verbose)
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			output_string, err := SerializeString(object, outputFormat, outputHeader, outputLimit)
			if err != nil {
				return "", errors.Wrap(err, "Error serializing output")
			}
			return output_string, nil
		case "csv", "tsv":
			return "", errors.New("Error: incompatible output format \"" + outputFormat + "\"")
		}
		return "", errors.New("Error: unknown output format \"" + outputFormat + "\"")
	case "jsonl", "csv", "tsv":
		switch outputFormat {
		case "bson", "json", "hcl", "hcl2", "toml", "yaml", "jsonl", "csv", "tsv":
			object, err := DeserializeBytes(inputBytes, inputFormat, inputHeader, inputComment, inputLazyQuotes, inputSkipLines, inputLimit, inputType, async, verbose)
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			output_string, err := SerializeString(object, outputFormat, outputHeader, outputLimit)
			if err != nil {
				return "", errors.Wrap(err, "Error serializing output")
			}
			return output_string, nil
		}
		return "", errors.New("Error: unknown output format \"" + inputFormat + "\"")
	}
	return "", errors.New("Error: unknown input format \"" + inputFormat + "\"")
}
