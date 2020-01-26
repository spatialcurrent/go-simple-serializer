// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package slicer provides an easy API to slice a string using a custom slice of indicies.
// See the examples below for usage.
package slicer

import (
	"fmt"
)

// Slicer provides an easy API to slice a series of strings using a custom x.
type Slicer struct {
	indicies []int
}

// New returns a new escaper with the default prefix "\" and no additional substrings.
// This default escaper will only escape the prefix itself, so "\" becomes "\\".
func New(indicies ...int) (*Slicer, error) {
	if len(indicies) > 2 {
		return nil, fmt.Errorf("invalid number of indicies %d", len(indicies))
	}
	return &Slicer{indicies: indicies}, nil
}

// Slice returns the sliced version of the given string.
func (s *Slicer) Slice(in string) string {
	if len(s.indicies) == 2 {
		if s.indicies[1] < 0 {
			return in[s.indicies[0] : len(in)+s.indicies[1]]
		}
		return in[s.indicies[0]:s.indicies[1]]
	}
	if len(s.indicies) == 1 {
		return in[s.indicies[0]:]
	}
	return in[:]
}
