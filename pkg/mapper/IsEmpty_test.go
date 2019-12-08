// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmptyString(t *testing.T) {
	assert.True(t, IsEmpty(""))
	assert.False(t, IsEmpty("hello world"))
}

func TestIsEmptySlice(t *testing.T) {
	assert.True(t, IsEmpty([]string{}))
	assert.False(t, IsEmpty([]string{"hello world"}))
}

func TestIsEmptyMap(t *testing.T) {
	assert.True(t, IsEmpty(map[string]string{}))
	assert.False(t, IsEmpty(map[string]string{"hello": "world"}))
}

func TestIsEmptyInt(t *testing.T) {
	assert.True(t, IsEmpty(0))
	assert.False(t, IsEmpty(1))
}
