// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package tags provides a simple API for reading and writing to lines of tags.
// tags also supports iterators for efficiently reading through a stream.
// See the examples below for usage.
//
//  - https://wiki.openstreetmap.org/wiki/Tags
package tags

import (
	"reflect"

	"github.com/pkg/errors"
)

const (
	quote rune = '"'
	space rune = ' '
)

var (
	DefaultType = reflect.TypeOf(map[string]string{})
)

var (
	ErrEmptyInput               = errors.New("empty input")
	ErrInvalidRune              = errors.New("invalid rune")
	ErrMissingKeyValueSeparator = errors.New("missing key-value separator")
	ErrMissingLineSeparator     = errors.New("missing line separator")
	ErrMissingKeySerializer     = errors.New("missing key serializer")
	ErrMissingValueSerializer   = errors.New("missing value serializer")
	ErrMissingType              = errors.New("missing type")
	ErrInvalidUTF8              = errors.New("invalid utf-8")
)
