// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package slicer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlicerZero(t *testing.T) {
	in := "Hello World"
	s, err := New()
	assert.NoError(t, err)
	out := s.Slice(in)
	assert.Equal(t, "Hello World", out)
}

func TestSlicerOne(t *testing.T) {
	in := "Hello World"
	s, err := New(6)
	assert.NoError(t, err)
	out := s.Slice(in)
	assert.Equal(t, "World", out)
}

func TestSlicerTwo(t *testing.T) {
	in := "Hello World"
	s, err := New(0, 5)
	assert.NoError(t, err)
	out := s.Slice(in)
	assert.Equal(t, "Hello", out)
}

func TestSlicerTwoNegative(t *testing.T) {
	in := "Hello World"
	s, err := New(0, -6)
	assert.NoError(t, err)
	out := s.Slice(in)
	assert.Equal(t, "Hello", out)
}
