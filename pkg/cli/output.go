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
)

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	FlagOutputURI               string = "output-uri"
	FlagOutputCompression       string = "output-compression"
	FlagOutputFormat            string = "output-format"
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

	DefaultOutputLimit = -1
)

var (
	DefaultOutputHeader = []string{}
)

var (
	ErrMissingLineSeparator = errors.New("line separator cannot be blank")
	ErrMissingEscapePrefix  = errors.New("escape prefix is missing")
)

type ErrMissingOutputFormat struct {
	Expected []string
}

func (e ErrMissingOutputFormat) Error() string {
	return fmt.Sprintf("missing output format, expecting one of %q", e.Expected)
}

type ErrInvalidOutputFormat struct {
	Value    string
	Expected []string
}

func (e ErrInvalidOutputFormat) Error() string {
	return fmt.Sprintf("invalid output format %q, expecting one of %q", e.Value, e.Expected)
}

// Initialize output flags
func InitOutputFlags(flag *pflag.FlagSet, formats []string) {
	flag.StringP(FlagOutputFormat, "o", "", "The output format.  One of: "+strings.Join(formats, ", "))
	flag.StringSlice(FlagOutputHeader, DefaultOutputHeader, "The output header if the stdout output has no header.")
	flag.IntP(FlagOutputLimit, "n", DefaultOutputLimit, "the output limit")
	flag.BoolP(FlagOutputPretty, "p", false, "print pretty output")
	flag.BoolP(FlagOutputSorted, "s", false, "sort output")
	flag.BoolP(FlagOutputReversed, "r", false, "if output is sorted, sort in reverse alphabetical order.")
	flag.BoolP(FlagOutputDecimal, "d", false, "when converting floats to strings use decimals rather than scientific notation")
	flag.Bool(FlagOutputKeyLower, false, "lower case output keys, including CSV column headers, tag names, and property names")
	flag.Bool(FlagOutputKeyUpper, false, "upper case output keys, including CSV column headers, tag names, and property names")
	flag.Bool(FlagOutputValueLower, false, "lower case output values, including tag values, and property values")
	flag.Bool(FlagOutputValueUpper, false, "upper case output values, including tag values, and property values")
	flag.StringP(FlagOutputNoDataValue, "0", "", "no data value, e.g., used for missing values when converting JSON to CSV")
	flag.String(FlagOutputLineSeparator, "\n", "override line separator.  Used with properties and JSONL formats.")
	flag.String(FlagOutputKeyValueSeparator, "=", "override key value separator.  Used with properties format.")
	flag.Bool(FlagOutputExpandHeader, false, "expand output header.  Used with CSV and TSV formats.")
	flag.String(FlagOutputEscapePrefix, "", "override escape prefix.  Used with properties format.")
	flag.Bool(FlagOutputEscapeColon, false, "Escape colon characters in output.  Used with properties format.")
	flag.Bool(FlagOutputEscapeEqual, false, "Escape equal characters in output.  Used with properties format.")
	flag.Bool(FlagOutputEscapeSpace, false, "Escape space characters in output.  Used with properties format.")
	flag.Bool(FlagOutputEscapeNewLine, false, "Escape new line characters in output.  Used with properties format.")
}

// CheckOutput checks the output configuration.
func CheckOutput(v *viper.Viper, formats []string) error {
	outputFormat := v.GetString(FlagOutputFormat)
	if len(outputFormat) == 0 {
		return &ErrMissingOutputFormat{Expected: formats}
	}
	if !stringSliceContains(formats, outputFormat) {
		return &ErrInvalidOutputFormat{Value: outputFormat, Expected: formats}
	}
	if ls := v.GetString(FlagOutputLineSeparator); len(ls) != 1 {
		return ErrMissingLineSeparator
	}
	if len(v.GetString(FlagOutputEscapePrefix)) == 0 {
		if v.GetBool(FlagOutputEscapeColon) {
			return errors.Wrap(ErrMissingEscapePrefix, "escaping colon requires an escape prefix")
		}
		if v.GetBool(FlagOutputEscapeEqual) {
			return errors.Wrap(ErrMissingEscapePrefix, "escaping equal requires an escape prefix")
		}
		if v.GetBool(FlagOutputEscapeSpace) {
			return errors.Wrap(ErrMissingEscapePrefix, "escaping space requires an escape prefix")
		}
		if v.GetBool(FlagOutputEscapeNewLine) {
			return errors.Wrap(ErrMissingEscapePrefix, "escaping new line requires an escape prefix")
		}
	}
	if v.GetBool(FlagOutputKeyLower) && v.GetBool(FlagOutputKeyUpper) {
		return errors.New("cannot lower case and upper case keys at the same time")
	}
	return nil
}
