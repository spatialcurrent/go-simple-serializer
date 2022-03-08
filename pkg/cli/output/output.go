// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package output contains the code for processing the user provided configuration for the output.
package output

import (
	"errors"
)

const (
	FlagOutputURI               string = "output-uri"
	FlagOutputCompression       string = "output-compression"
	FlagOutputFormat            string = "output-format"
	FlagOutputFormatSpecifier   string = "output-format-specifier"
	FlagOutputFit               string = "output-fit"
	FlagOutputPretty            string = "output-pretty"
	FlagOutputHeader            string = "output-header"
	FlagOutputLimit             string = "output-limit"
	FlagOutputAppend            string = "output-append"
	FlagOutputOverwrite         string = "output-overwrite"
	FlagOutputBufferMemory      string = "output-buffer-memory"
	FlagOutputMkdirs            string = "output-mkdirs"
	FlagOutputPassphrase        string = "output-passphrase"
	FlagOutputSalt              string = "output-salt"
	FlagOutputDecimal           string = "output-decimal"
	FlagOutputKeyLower          string = "output-key-lower"
	FlagOutputKeyUpper          string = "output-key-upper"
	FlagOutputValueLower        string = "output-value-lower"
	FlagOutputValueUpper        string = "output-value-upper"
	FlagOutputNoDataValue       string = "output-no-data-value"
	FlagOutputLineSeparator     string = "output-line-separator"
	FlagOutputKeyValueSeparator string = "output-key-value-separator"
	FlagOutputExpandHeader      string = "output-expand-header"
	FlagOutputEscapePrefix      string = "output-escape-prefix"
	FlagOutputEscapeColon       string = "output-escape-colon"
	FlagOutputEscapeEqual       string = "output-escape-equal"
	FlagOutputEscapeNewLine     string = "output-escape-new-line"
	FlagOutputEscapeSpace       string = "output-escape-space"
	FlagOutputSorted            string = "output-sorted"
	FlagOutputReversed          string = "output-reversed"
	FlagOutputType              string = "output-type"

	DefaultOutputLimit = -1
)

var (
	ErrMissingOutputKeyValueSeparator = errors.New("missing output key-value separator")
	ErrMissingOutputLineSeparator     = errors.New("missing output line separator")
	ErrMissingOutputEscapePrefix      = errors.New("missing output escape prefix")
)

var (
	DefaultOutputHeader = []string{}
)

func stringSliceContains(slc []string, str string) bool {
	for _, x := range slc {
		if x == str {
			return true
		}
	}
	return false
}
