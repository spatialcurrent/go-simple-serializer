// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package escaper

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestEscape(t *testing.T) {
	in := "Hello Beautiful\\World Again"
	e := New().Prefix("\\")
	out := e.Escape(in)
	assert.Equal(t, "Hello Beautiful\\\\World Again", out)
}

func TestEscapePrefix(t *testing.T) {
	in := "Hello Beautiful\\World Again"
	e := New().Prefix("*")
	out := e.Escape(in)
	assert.Equal(t, "Hello Beautiful\\World Again", out)
}

func TestEscapeEqual(t *testing.T) {
	in := "Hello=World"
	e := New().Prefix("\\").Sub("=")
	out := e.Escape(in)
	assert.Equal(t, "Hello\\=World", out)
}

func TestEscapeEqualColonSpace(t *testing.T) {
	in := "Hello=World\nCiao Ciao"
	e := New().Prefix("\\").Sub("=", ":", " ")
	out := e.Escape(in)
	assert.Equal(t, "Hello\\=World\nCiao\\ Ciao", out)
}

func TestUnescape(t *testing.T) {
	in := "Hello Beautiful\\\\World Again"
	e := New().Prefix("\\")
	out := e.Unescape(in)
	assert.Equal(t, "Hello Beautiful\\World Again", out)
}

func TestUnescapePrefix(t *testing.T) {
	in := "Hello Beautiful\\World Again"
	e := New().Prefix("*")
	out := e.Unescape(in)
	assert.Equal(t, "Hello Beautiful\\World Again", out)
}

func TestUnescapeEqual(t *testing.T) {
	in := "Hello\\=World"
	e := New().Prefix("\\").Sub("=")
	out := e.Unescape(in)
	assert.Equal(t, "Hello=World", out)
}

func TestUnescapeEqualColonSpace(t *testing.T) {
	in := "Hello\\=World\nCiao\\ Ciao"
	e := New().Prefix("\\").Sub("=", ":", " ")
	out := e.Unescape(in)
	assert.Equal(t, "Hello=World\nCiao Ciao", out)
}
