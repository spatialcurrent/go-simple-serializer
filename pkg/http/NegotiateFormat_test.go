// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package http

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

func TestNegotiateFormatJSON(t *testing.T) {
	reg := NewDefaultRegistry()
	r := httptest.NewRequest("GET", "https://example.com/foo/bar", nil)
	r.Header.Set("Accept", "application/json")
	c, f, err := NegotiateFormat(r, reg)
	assert.NoError(t, err)
	assert.Equal(t, "application/json", c)
	assert.Equal(t, serializer.FormatJSON, f)
}

func TestNegotiateFormatBSON(t *testing.T) {
	reg := NewDefaultRegistry()
	r := httptest.NewRequest("GET", "https://example.com/foo/bar", nil)
	r.Header.Set("Accept", "application/ubjson, application/json")
	c, f, err := NegotiateFormat(r, reg)
	assert.NoError(t, err)
	assert.Equal(t, "application/ubjson", c)
	assert.Equal(t, serializer.FormatBSON, f)
}

func TestNegotiateFormatWeight(t *testing.T) {
	reg := NewDefaultRegistry()
	r := httptest.NewRequest("GET", "https://example.com/foo/bar", nil)
	r.Header.Set("Accept", "text/csv;q=0.8, application/json;q=0.9")
	c, f, err := NegotiateFormat(r, reg)
	assert.NoError(t, err)
	assert.Equal(t, "application/json", c)
	assert.Equal(t, serializer.FormatJSON, f)
}
