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

// ConvertInput provides the input for the Convert function.
type ConvertInput struct {
	InputBytes              []byte
	InputFormat             string
	InputHeader             []interface{}
	InputComment            string
	InputLazyQuotes         bool
	InputSkipLines          int
	InputLimit              int
	InputLineSeparator      string
	InputDropCR             bool
	OutputFormat            string
	OutputHeader            []interface{}
	OutputLimit             int
	OutputPretty            bool
	OutputSorted            bool
	OutputValueSerializer   func(object interface{}) (string, error)
	OutputLineSeparator     string
	OutputKeyValueSeparator string
	OutputEscapePrefix      string
	OutputEscapeSpace       bool
	OutputEscapeNewLine     bool
	OutputEscapeEqual       bool
	Async                   bool
	Verbose                 bool
}

func NewConvertInput(bytes []byte, inputFormat string, outputFormat string) *ConvertInput {
	return &ConvertInput{
		InputBytes:              bytes,
		InputFormat:             inputFormat,
		InputHeader:             NoHeader,
		InputComment:            NoComment,
		InputLazyQuotes:         false,
		InputSkipLines:          NoSkip,
		InputLimit:              NoLimit,
		InputLineSeparator:      "\n",
		InputDropCR:             true,
		OutputFormat:            outputFormat,
		OutputHeader:            NoHeader,
		OutputLimit:             NoLimit,
		OutputPretty:            false,
		OutputSorted:            false,
		OutputValueSerializer:   nil,
		OutputLineSeparator:     "\n",
		OutputKeyValueSeparator: "=",
		OutputEscapePrefix:      "\\",
		OutputEscapeSpace:       false,
		OutputEscapeNewLine:     false,
		OutputEscapeEqual:       false,
		Async:                   false,
		Verbose:                 false,
	}
}

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
				Bytes:         input.InputBytes,
				Format:        input.InputFormat,
				Header:        input.InputHeader,
				Comment:       input.InputComment,
				LazyQuotes:    input.InputLazyQuotes,
				SkipLines:     input.InputSkipLines,
				Limit:         input.InputLimit,
				LineSeparator: input.InputLineSeparator,
				DropCR:        input.InputDropCR,
				Type:          inputType,
				Async:         input.Async,
				Verbose:       input.Verbose,
			})
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			outputString, err := SerializeString(&SerializeInput{
				Object:            object,
				Format:            input.OutputFormat,
				Header:            input.OutputHeader,
				Limit:             input.OutputLimit,
				Pretty:            input.OutputPretty,
				Sorted:            input.OutputSorted,
				LineSeparator:     input.OutputLineSeparator,
				KeyValueSeparator: input.OutputKeyValueSeparator,
				ValueSerializer:   input.OutputValueSerializer,
				EscapePrefix:      input.OutputEscapePrefix,
				EscapeSpace:       input.OutputEscapeSpace,
				EscapeNewLine:     input.OutputEscapeNewLine,
				EscapeEqual:       input.OutputEscapeEqual,
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
				Bytes:         input.InputBytes,
				Format:        input.InputFormat,
				Header:        input.InputHeader,
				Comment:       input.InputComment,
				LazyQuotes:    input.InputLazyQuotes,
				SkipLines:     input.InputSkipLines,
				Limit:         input.InputLimit,
				LineSeparator: input.InputLineSeparator,
				DropCR:        input.InputDropCR,
				Type:          inputType,
				Async:         input.Async,
				Verbose:       input.Verbose,
			})
			if err != nil {
				return "", errors.Wrap(err, "Error deserializing input")
			}
			outputString, err := SerializeString(&SerializeInput{
				Object:            object,
				Format:            input.OutputFormat,
				Header:            input.OutputHeader,
				Limit:             input.OutputLimit,
				Pretty:            input.OutputPretty,
				Sorted:            input.OutputSorted,
				LineSeparator:     input.OutputLineSeparator,
				KeyValueSeparator: input.OutputKeyValueSeparator,
				ValueSerializer:   input.OutputValueSerializer,
				EscapePrefix:      input.OutputEscapePrefix,
				EscapeSpace:       input.OutputEscapeSpace,
				EscapeNewLine:     input.OutputEscapeNewLine,
				EscapeEqual:       input.OutputEscapeEqual,
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
