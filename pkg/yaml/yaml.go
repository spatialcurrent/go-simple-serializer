// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package yaml provides an API for YAML serialization that can automatically infers types.
//
//	- https://yaml.org/
//	- https://en.wikipedia.org/wiki/YAML
package yaml

import (
	"github.com/pkg/errors"
)

var (
	True           = []byte("true")
	False          = []byte("false")
	Null           = []byte("null")
	BoundaryMarker = []byte("---")
	Y              = []byte("y")
)

var (
	ErrEmptyInput  = errors.New("empty input")
	ErrInvalidRune = errors.New("invalid rune")
)
