// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package gob provides support for reading and writing a stream of gob-encoded obejcts.
package gob

import (
	"errors"
)

var (
	ErrMissingType = errors.New("missing type")
	ErrInvalidType = errors.New("invalid type, expecting map or struct")
)
