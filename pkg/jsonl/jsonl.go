// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"github.com/pkg/errors"
)

var (
	ErrMissingType          = errors.New("missing type")
	ErrMissingLineSeparator = errors.New("missing line separator")
)
