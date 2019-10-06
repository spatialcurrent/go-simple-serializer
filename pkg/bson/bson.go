// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package bson provides an API for BSON serialization.  This package wraps the mgo bson package.
//	- https://godoc.org/gopkg.in/mgo.v2/bson
package bson

import (
	"reflect"

	"github.com/pkg/errors"
)

var (
	DefaultType = reflect.TypeOf(map[string]interface{}{})
)

var (
	ErrEmptyInput  = errors.New("empty input")
	ErrInvalidRune = errors.New("invalid rune")
)
