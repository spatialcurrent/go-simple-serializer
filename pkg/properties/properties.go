// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package properties includes functions for reading and writing from properties files.
// See the examples below for usage.
//
// Reference:
//  - https://en.wikipedia.org/wiki/.properties
//
package properties

import (
	"errors"
)

var (
	ErrMissingLineSeparator     = errors.New("missing line separator")
	ErrMissingKeyValueSeparator = errors.New("missing key-value separator")
)
