// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package toml provides an API for TOML serialization.  This package wraps the BurntSushi toml package.
//	- https://github.com/BurntSushi/toml
package toml

import (
	"errors"
	"reflect"
)

var (
	DefaultType = reflect.TypeOf(map[string]interface{}{})
)

var (
	ErrEmptyInput = errors.New("empty input")
	// TOML cannot marshal nil values, because of a design decision.
	// https://github.com/toml-lang/toml/issues/30
	ErrNilObject = errors.New("nil object, toml cannot marshal nil values")
)
