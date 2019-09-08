// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package cli manages loading and testing configuration for serializing and deserializing objects from the command line.
package cli

import (
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var (
	ErrMissingKeyValueSeparator = errors.New("missing key-value separator")
	ErrMissingLineSeparator     = errors.New("missing line separator")
	ErrMissingEscapePrefix      = errors.New("missing escape prefix")
)

const (
	FlagNoStream string = "no-stream"
	FlagVerbose  string = "verbose"
)

// Initialize cli flags
func InitCliFlags(flag *pflag.FlagSet, formats []string) {

	InitInputFlags(flag, formats)

	InitOutputFlags(flag, formats)

	flag.Bool(FlagNoStream, false, "disable streaming")
	flag.BoolP(FlagVerbose, "v", false, "verbose output")
}
