// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gob

import (
	"github.com/pkg/errors"
)

var (
	ErrMissingType = errors.New("missing type")
	ErrInvalidType = errors.New("invalid type, expecting map or struct")
)
