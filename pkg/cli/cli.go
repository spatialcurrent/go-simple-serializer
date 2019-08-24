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
)

var (
	ErrMissingKeyValueSeparator = errors.New("missing key-value separator")
	ErrMissingLineSeparator     = errors.New("missing line separator")
	ErrMissingEscapePrefix      = errors.New("missing escape prefix")
)
