// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
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
	InputTrim               bool
	InputEscapePrefix       string
	InputUnescapeSpace      bool
	InputUnescapeNewLine    bool
	InputUnescapeEqual      bool
	OutputFormat            string
	OutputHeader            []interface{}
	OutputLimit             int
	OutputPretty            bool
	OutputSorted            bool
	OutputReversed          bool
	OutputKeySerializer     stringify.Stringer
	OutputValueSerializer   stringify.Stringer
	OutputLineSeparator     string
	OutputKeyValueSeparator string
	OutputEscapePrefix      string
	OutputEscapeSpace       bool
	OutputEscapeNewLine     bool
	OutputEscapeEqual       bool
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
		InputEscapePrefix:       "\\",
		InputUnescapeSpace:      false,
		InputUnescapeNewLine:    false,
		InputUnescapeEqual:      false,
		OutputFormat:            outputFormat,
		OutputHeader:            NoHeader,
		OutputLimit:             NoLimit,
		OutputPretty:            false,
		OutputSorted:            false,
		OutputReversed:          false,
		OutputKeySerializer:     nil,
		OutputValueSerializer:   nil,
		OutputLineSeparator:     "\n",
		OutputKeyValueSeparator: "=",
		OutputEscapePrefix:      "\\",
		OutputEscapeSpace:       false,
		OutputEscapeNewLine:     false,
		OutputEscapeEqual:       false,
	}
}

func Convert(input *ConvertInput) ([]byte, error) {
	in := serializer.New(input.InputFormat).
		Limit(input.InputLimit).
		Header(input.InputHeader).
		Comment(input.InputComment).
		LazyQuotes(input.InputLazyQuotes).
		SkipLines(input.InputSkipLines).
		LineSeparator(input.InputLineSeparator).
		DropCR(input.InputDropCR).
		Trim(input.InputTrim).
		EscapePrefix(input.InputEscapePrefix).
		UnescapeEqual(input.InputUnescapeEqual).
		UnescapeSpace(input.InputUnescapeSpace).
		UnescapeNewLine(input.InputUnescapeNewLine)

	obj, err := in.Deserialize(input.InputBytes)
	if err != nil {
		return make([]byte, 0), errors.Wrap(err, "error deserializing input")
	}

	out := serializer.New(input.OutputFormat).
		Limit(input.OutputLimit).
		Header(input.OutputHeader).
		Pretty(input.OutputPretty).
		Sorted(input.OutputSorted).
		Reversed(input.OutputReversed).
		KeySerializer(input.OutputKeySerializer).
		ValueSerializer(input.OutputValueSerializer).
		LineSeparator(input.OutputLineSeparator).
		KeyValueSeparator(input.OutputKeyValueSeparator).
		EscapePrefix(input.OutputEscapePrefix).
		EscapeEqual(input.OutputEscapeEqual).
		EscapeSpace(input.OutputEscapeSpace).
		EscapeNewLine(input.OutputEscapeNewLine)

	b, err := out.Serialize(obj)
	if err != nil {
		return make([]byte, 0), errors.Wrap(err, "error serializing output")
	}

	return b, nil
}

/*

// Convert converts an input_string from the inputFormat to the outputFormat.
// Returns the output string and error, if any.
//func Convert(inputBytes []byte, inputFormat string, inputHeader []string, inputComment string, inputLazyQuotes bool, inputSkipLines int, inputLimit int, outputFormat string, outputHeader []string, outputLimit int, async bool, verbose bool) (string, error) {
func OldConvert(input *ConvertInput) ([]byte, error) {

	inputType, err := GetType(input.InputBytes, input.InputFormat)
	if err != nil {
		return make([]byte, 0), errors.Wrap(err, "error creating new object for format "+input.InputFormat)
	}

	if input.Verbose {
		fmt.Println("Input Format: " + input.InputFormat)
		fmt.Println("Output Format: " + input.OutputFormat)
		fmt.Println("Input Type: " + fmt.Sprint(inputType))
	}

	switch input.InputFormat {
	case "bson", "json", "hcl", "hcl2", "properties", "tags", "toml", "yaml":
		switch input.OutputFormat {
		case "bson", "json", "jsonl", "hcl", "hcl2", "properties", "tags", "toml", "yaml":
			object, err := DeserializeBytes(&DeserializeBytesInput{
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
				return make([]byte, 0), errors.Wrap(err, "Error deserializing input")
			}
			outputBytes, err := SerializeBytes(&SerializeBytesInput{
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
				return make([]byte, 0), errors.Wrap(err, "Error serializing output")
			}
			return outputBytes, nil
		case "csv", "tsv":
			return make([]byte, 0), &ErrIncompatibleFormats{Input: input.InputFormat, Output: input.OutputFormat}
		}
		return make([]byte, 0), errors.Wrap(&ErrUnknownFormat{Name: input.OutputFormat}, "unknown output format")
	case "jsonl", "csv", "tsv":
		switch input.OutputFormat {
		case "bson", "json", "hcl", "hcl2", "toml", "yaml", "jsonl", "csv", "tags", "tsv":
			object, err := DeserializeBytes(&DeserializeBytesInput{
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
				return make([]byte, 0), errors.Wrap(err, "Error deserializing input")
			}
			outputBytes, err := SerializeBytes(&SerializeBytesInput{
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
				return make([]byte, 0), errors.Wrap(err, "Error serializing output")
			}
			return outputBytes, nil
		}
		return make([]byte, 0), errors.Wrap(&ErrUnknownFormat{Name: input.OutputFormat}, "unknown output format")
	}
	return make([]byte, 0), errors.Wrap(&ErrUnknownFormat{Name: input.InputFormat}, "unknown output format")
}
*/
