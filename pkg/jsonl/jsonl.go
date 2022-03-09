// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package jsonl provides a simple API for reading and writing to JSON Lines (aka jsonl).
// jsonl also supports iterators for efficiently reading through a stream.
// jsonl uses the github.com/spatialcurrent/go-simple-serializer/pkg/json for marshaling/unmarshaling JSON.
// See the examples below for usage.
//
//  - https://godoc.org/pkg/github.com/spatialcurrent/go-simple-serializer/pkg/json
package jsonl

import (
	"errors"
)

var (
	ErrMissingType          = errors.New("missing type")
	ErrMissingLineSeparator = errors.New("missing line separator")
)
