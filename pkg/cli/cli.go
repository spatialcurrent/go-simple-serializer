// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package cli manages loading and testing configuration for serializing and deserializing objects from the command line.
package cli

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/cli/input"
	"github.com/spatialcurrent/go-simple-serializer/pkg/cli/output"
)

const (
	FlagNoStream string = "no-stream"
	FlagVerbose  string = "verbose"
)

const (
	FlagInputURI               = input.FlagInputURI
	FlagInputCompression       = input.FlagInputCompression
	FlagInputFormat            = input.FlagInputFormat
	FlagInputHeader            = input.FlagInputHeader
	FlagInputLimit             = input.FlagInputLimit
	FlagInputComment           = input.FlagInputComment
	FlagInputLazyQuotes        = input.FlagInputLazyQuotes
	FlagInputTrim              = input.FlagInputTrim
	FlagInputReaderBufferSize  = input.FlagInputReaderBufferSize
	FlagInputScannerBufferSize = input.FlagInputScannerBufferSize
	FlagInputSkipLines         = input.FlagInputSkipLines
	FlagInputLineSeparator     = input.FlagInputLineSeparator
	FlagInputKeyValueSeparator = input.FlagInputKeyValueSeparator
	FlagInputDropCR            = input.FlagInputDropCR
	FlagInputEscapePrefix      = input.FlagInputEscapePrefix
	FlagInputUnescapeColon     = input.FlagInputUnescapeColon
	FlagInputUnescapeEqual     = input.FlagInputUnescapeEqual
	FlagInputUnescapeSpace     = input.FlagInputUnescapeSpace
	FlagInputUnescapeNewLine   = input.FlagInputUnescapeNewLine
	FlagInputType              = input.FlagInputType
)

const (
	FlagOutputURI               = output.FlagOutputURI
	FlagOutputCompression       = output.FlagOutputCompression
	FlagOutputFormat            = output.FlagOutputFormat
	FlagOutputFormatSpecifier   = output.FlagOutputFormatSpecifier
	FlagOutputFit               = output.FlagOutputFit
	FlagOutputPretty            = output.FlagOutputPretty
	FlagOutputHeader            = output.FlagOutputHeader
	FlagOutputLimit             = output.FlagOutputLimit
	FlagOutputAppend            = output.FlagOutputAppend
	FlagOutputOverwrite         = output.FlagOutputOverwrite
	FlagOutputBufferMemory      = output.FlagOutputBufferMemory
	FlagOutputMkdirs            = output.FlagOutputMkdirs
	FlagOutputPassphrase        = output.FlagOutputPassphrase
	FlagOutputSalt              = output.FlagOutputSalt
	FlagOutputDecimal           = output.FlagOutputDecimal
	FlagOutputKeyLower          = output.FlagOutputKeyLower
	FlagOutputKeyUpper          = output.FlagOutputKeyUpper
	FlagOutputValueLower        = output.FlagOutputValueLower
	FlagOutputValueUpper        = output.FlagOutputValueUpper
	FlagOutputNoDataValue       = output.FlagOutputNoDataValue
	FlagOutputLineSeparator     = output.FlagOutputLineSeparator
	FlagOutputKeyValueSeparator = output.FlagOutputKeyValueSeparator
	FlagOutputExpandHeader      = output.FlagOutputExpandHeader
	FlagOutputEscapePrefix      = output.FlagOutputEscapePrefix
	FlagOutputEscapeColon       = output.FlagOutputEscapeColon
	FlagOutputEscapeEqual       = output.FlagOutputEscapeEqual
	FlagOutputEscapeNewLine     = output.FlagOutputEscapeNewLine
	FlagOutputEscapeSpace       = output.FlagOutputEscapeSpace
	FlagOutputSorted            = output.FlagOutputSorted
	FlagOutputReversed          = output.FlagOutputReversed
	FlagOutputType              = output.FlagOutputType
)
