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
)

// Convert converts an input_string from the inputFormat to the outputFormat.
// Returns the output string and error, if any.
//func Convert(inputBytes []byte, inputFormat string, inputHeader []string, inputComment string, inputLazyQuotes bool, inputSkipLines int, inputLimit int, outputFormat string, outputHeader []string, outputLimit int, async bool, verbose bool) (string, error) {
func Convert(input *ConvertInput) (string, error) {

	inputType, err := GetType(input.InputBytes, input.InputFormat)
	if err != nil {
		return "", errors.Wrap(err, "error creating new object for format "+input.InputFormat)
	}

	if input.Verbose {
		fmt.Println("Input Format: " + input.InputFormat)
		fmt.Println("Output Format: " + input.OutputFormat)
		fmt.Println("Input Type: " + fmt.Sprint(inputType))
	}

	switch input.InputFormat {
	case "bson", "json", "hcl", "hcl2", "properties", "toml", "yaml":
		switch input.OutputFormat {
		case "bson", "json", "jsonl", "hcl", "hcl2", "properties", "toml", "yaml":
			object, err := DeserializeBytes(&DeserializeInput{
				Bytes:      input.InputBytes,
				Format:     input.InputFormat,
				Header:     input.InputHeader,
				Comment:    input.InputComment,
				LazyQuotes: input.InputLazyQuotes,
				SkipLines:  input.InputSkipLines,
				Limit:      input.InputLimit,
				Type:       inputType,
				Async:      input.Async,
				Verbose:    input.Verbose,
			})
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			outputString, err := SerializeString(&SerializeInput{
				Object: object,
				Format: input.OutputFormat,
				Header: input.OutputHeader,
				Limit:  input.OutputLimit,
				Pretty: input.OutputPretty,
				Sorted: input.OutputSorted,
			})
			if err != nil {
				return "", errors.Wrap(err, "Error serializing output")
			}
			return outputString, nil
		case "csv", "tsv":
			return "", &ErrIncompatibleFormats{Input: input.InputFormat, Output: input.OutputFormat}
		}
		return "", errors.Wrap(&ErrUnknownFormat{Name: input.OutputFormat}, "unknown output format")
	case "jsonl", "csv", "tsv":
		switch input.OutputFormat {
		case "bson", "json", "hcl", "hcl2", "toml", "yaml", "jsonl", "csv", "tsv":
			object, err := DeserializeBytes(&DeserializeInput{
				Bytes:      input.InputBytes,
				Format:     input.InputFormat,
				Header:     input.InputHeader,
				Comment:    input.InputComment,
				LazyQuotes: input.InputLazyQuotes,
				SkipLines:  input.InputSkipLines,
				Limit:      input.InputLimit,
				Type:       inputType,
				Async:      input.Async,
				Verbose:    input.Verbose,
			})
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			outputString, err := SerializeString(&SerializeInput{
				Object:          object,
				Format:          input.OutputFormat,
				Header:          input.OutputHeader,
				Limit:           input.OutputLimit,
				Pretty:          input.OutputPretty,
				Sorted:          input.OutputSorted,
				ValueSerializer: input.OutputValueSerializer,
			})
			if err != nil {
				return "", errors.Wrap(err, "Error serializing output")
			}
			return outputString, nil
		}
		return "", errors.Wrap(&ErrUnknownFormat{Name: input.OutputFormat}, "unknown output format")
	}
	return "", errors.Wrap(&ErrUnknownFormat{Name: input.InputFormat}, "unknown output format")
}
