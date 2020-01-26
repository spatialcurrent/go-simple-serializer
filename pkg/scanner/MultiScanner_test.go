// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package scanner

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	ms := MultiScanner(
		New(strings.NewReader("hello\nworld"), '\n', true),
		New(strings.NewReader("hello\nbeautiful\nworld"), '\n', true),
	)
	valid := ms.Scan()
	assert.True(t, valid)
	assert.NoError(t, ms.Err())
	assert.Equal(t, "hello", ms.Text())
	valid = ms.Scan()
	assert.True(t, valid)
	assert.NoError(t, ms.Err())
	assert.Equal(t, "world", ms.Text())
	valid = ms.Scan()
	assert.True(t, valid)
	assert.NoError(t, ms.Err())
	assert.Equal(t, "hello", ms.Text())
	valid = ms.Scan()
	assert.True(t, valid)
	assert.NoError(t, ms.Err())
	assert.Equal(t, "beautiful", ms.Text())
	valid = ms.Scan()
	assert.True(t, valid)
	assert.NoError(t, ms.Err())
	assert.Equal(t, "world", ms.Text())
	valid = ms.Scan()
	assert.False(t, valid)
	assert.NoError(t, ms.Err())
	assert.Equal(t, "", ms.Text())
	valid = ms.Scan()
	assert.False(t, valid)
	assert.NoError(t, ms.Err())
	assert.Equal(t, "", ms.Text())
}
