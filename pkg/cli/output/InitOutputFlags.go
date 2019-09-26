// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package output

import (
	"github.com/spf13/pflag"
)

// Initialize output flags
func InitOutputFlags(flag *pflag.FlagSet) {
	flag.StringP(FlagOutputFormat, "o", "", "The output format")
	flag.String(FlagOutputFormatSpecifier, "", "The output format specifier")
	flag.Bool(FlagOutputFit, false, "Fit output")
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
	flag.String(FlagOutputType, "", "if using GOB format, the output type, default map[string]interface {}")
}
