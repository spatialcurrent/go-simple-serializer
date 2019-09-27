// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package input

import (
	"github.com/spf13/pflag"
)

// InitInputFlags initializes the flags for processing the input data from the gss command.
func InitInputFlags(flag *pflag.FlagSet) {
	flag.StringP(FlagInputFormat, "i", "", "The input format")
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
	flag.String(FlagInputType, "", "if using GOB format, input type, default map[string]interface {}")
}
