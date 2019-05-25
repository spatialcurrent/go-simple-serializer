// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

// ConvertInput provides the input for the Convert function.
type ConvertInput struct {
	InputBytes              []byte
	InputFormat             string
	InputHeader             []string
	InputComment            string
	InputLazyQuotes         bool
	InputSkipLines          int
	InputLimit              int
	InputLineSeparator      string
	InputDropCR             bool
	OutputFormat            string
	OutputHeader            []string
	OutputLimit             int
	OutputPretty            bool
	OutputSorted            bool
	OutputValueSerializer   func(object interface{}) (string, error)
	OutputLineSeparator     string
	OutputKeyValueSeparator string
	OutputEscapePrefix      string
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
		Async:                   false,
		Verbose:                 false,
	}
}
