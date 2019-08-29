// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package cli

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

	DefaultSkipLines  int = 0
	DefaultInputLimit int = -1
)

var (
	DefaultInputHeader = []string{}
)

type ErrMissingInputFormat struct {
	Expected []string
}

func (e ErrMissingInputFormat) Error() string {
	return fmt.Sprintf("missing input format, expecting one of %q", e.Expected)
}

type ErrInvalidInputFormat struct {
	Value    string
	Expected []string
}

func (e ErrInvalidInputFormat) Error() string {
	return fmt.Sprintf("invalid input format %q, expecting one of %q", e.Value, e.Expected)
}

// Initialize input flags
func InitInputFlags(flag *pflag.FlagSet, formats []string) {
	flag.StringP(FlagInputFormat, "i", "", "The input format.  One of: "+strings.Join(formats, ", "))
	flag.StringSlice(FlagInputHeader, DefaultInputHeader, "The input header if the stdin input has no header.")
	flag.StringP(FlagInputComment, "c", "", "The input comment character, e.g., #.  Commented lines are not sent to output.")
	flag.Bool(FlagInputLazyQuotes, false, "allows lazy quotes for CSV and TSV")
	flag.Int(FlagInputReaderBufferSize, 4096, "the buffer size of the file reader")
	flag.Int(FlagInputScannerBufferSize, 0, "the initial buffer size for the scanner")
	flag.Int(FlagInputSkipLines, DefaultSkipLines, "The number of lines to skip before processing")
	flag.IntP(FlagInputLimit, "l", DefaultInputLimit, "The input limit")
	flag.BoolP(FlagInputTrim, "t", false, "trim input lines")
	flag.String(FlagInputLineSeparator, "\n", "override line separator.  Used with properties and JSONL formats.")
	flag.String(FlagInputKeyValueSeparator, "=", "override key-value separator.  not used.")
	flag.Bool(FlagInputDropCR, false, "drop carriage return characters that immediately precede new line characters")
	flag.String(FlagInputEscapePrefix, "", "override escape prefix.  Used with properties format.")
	flag.Bool(FlagInputUnescapeColon, false, "Unescape colon characters in input.  Used with properties format.")
	flag.Bool(FlagInputUnescapeEqual, false, "Unescape equal characters in input.  Used with properties format.")
	flag.Bool(FlagInputUnescapeSpace, false, "Unescape space characters in input.  Used with properties format.")
	flag.Bool(FlagInputUnescapeNewLine, false, "Unescape new line characters in input.  Used with properties format.")
}

// CheckInput checks the output configuration.
func CheckInput(v *viper.Viper, formats []string) error {
	inputFormat := v.GetString(FlagInputFormat)
	if len(inputFormat) == 0 {
		return &ErrMissingInputFormat{Expected: formats}
	}
	if !stringSliceContains(formats, inputFormat) {
		return &ErrInvalidInputFormat{Value: inputFormat, Expected: formats}
	}
	if ls := v.GetString(FlagInputLineSeparator); len(ls) != 1 {
		return ErrMissingLineSeparator
	}
	if ls := v.GetString(FlagInputKeyValueSeparator); len(ls) != 1 {
		return ErrMissingKeyValueSeparator
	}
	if len(v.GetString(FlagInputEscapePrefix)) == 0 {
		if v.GetBool(FlagInputUnescapeColon) {
			return errors.Wrap(ErrMissingEscapePrefix, "unescaping colon requires an escape prefix")
		}
		if v.GetBool(FlagInputUnescapeEqual) {
			return errors.Wrap(ErrMissingEscapePrefix, "unescaping equal requires an escape prefix")
		}
		if v.GetBool(FlagInputUnescapeSpace) {
			return errors.Wrap(ErrMissingEscapePrefix, "unescaping space requires an escape prefix")
		}
		if v.GetBool(FlagInputUnescapeNewLine) {
			return errors.Wrap(ErrMissingEscapePrefix, "unescaping new line requires an escape prefix")
		}
	}
	return nil
}
