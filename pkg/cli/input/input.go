// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package input contains the code for processing the user provided configuration for the input.
package input

import (
	"github.com/pkg/errors"
)

const (
	FlagInputURI               string = "input-uri"
	FlagInputCompression       string = "input-compression"
	FlagInputFormat            string = "input-format"
	FlagInputHeader            string = "input-header"
	FlagInputLimit             string = "input-limit"
	FlagInputComment           string = "input-comment"
	FlagInputLazyQuotes        string = "input-lazy-quotes"
	FlagInputTrim              string = "input-trim"
	FlagInputReaderBufferSize  string = "input-reader-buffer-size"
	FlagInputScannerBufferSize string = "input-scanner-buffer-size"
	FlagInputSkipLines         string = "input-skip-lines"
	FlagInputLineSeparator     string = "input-line-separator"
	FlagInputKeyValueSeparator string = "input-key-value-separator"
	FlagInputDropCR            string = "input-drop-cr"
	FlagInputEscapePrefix      string = "input-escape-prefix"
	FlagInputUnescapeColon     string = "input-unescape-colon"
	FlagInputUnescapeEqual     string = "input-unescape-equal"
	FlagInputUnescapeSpace     string = "input-unescape-space"
	FlagInputUnescapeNewLine   string = "input-unescape-new-line"
	FlagInputType              string = "input-type"

	DefaultSkipLines  int = 0
	DefaultInputLimit int = -1
)

var (
	ErrMissingInputKeyValueSeparator = errors.New("missing input key-value separator")
	ErrMissingInputLineSeparator     = errors.New("missing input line separator")
	ErrMissingInputEscapePrefix      = errors.New("missing input escape prefix")
)

var (
	DefaultInputHeader = []string{}
)

func stringSliceContains(slc []string, str string) bool {
	for _, x := range slc {
		if x == str {
			return true
		}
	}
	return false
}
