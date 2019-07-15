// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package gssjs includes functions for the JavaScript build of GSS.
//
package gssjs

import (
	"github.com/pkg/errors"
)

const (
	NoLimit = -1
)

var (
	ErrMissingInputString  = errors.New("missing input string")
	ErrMissingInputObject  = errors.New("missing input object")
	ErrMissingInputFormat  = errors.New("missing input format")
	ErrMissingOutputFormat = errors.New("missing output format")
)

var (
	DeserializeDefaults = map[string]interface{}{
		"limit":             NoLimit,
		"keyValueSeparator": "=",
		"lineSeparator":     "\n",
		"escapePrefix":      "\\",
		"expandHeader":      true,
	}
	SerializeDefaults = map[string]interface{}{
		"limit":             NoLimit,
		"keyValueSeparator": "=",
		"lineSeparator":     "\n",
		"escapePrefix":      "\\",
	}
)
