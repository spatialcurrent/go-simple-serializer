// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package json provides an API for JSON serialization that automatically infers types.
// See the examples below for usage.
package json

import (
	"github.com/pkg/errors"
)

var (
	ErrEmptyInput  = errors.New("empty input")
	ErrInvalidRune = errors.New("invalid rune")

	BytesTrue  = []byte("true")
	BytesFalse = []byte("false")
	BytesNull  = []byte("null")
)
