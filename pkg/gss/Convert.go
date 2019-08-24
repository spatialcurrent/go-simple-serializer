// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"github.com/pkg/errors"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
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
