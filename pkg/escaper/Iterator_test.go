// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package escaper

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	in := `properties:addr\:city:Maur:addr\:country`
	it := NewIterator(
		strings.NewReader(in),
		[]byte("\\"),
		[][]byte{
			[]byte("\\"),
			[]byte(":"),
		},
		[]byte(":"),
	)
	out, err := it.Next()
	assert.NoError(t, err)
	assert.Equal(t, "properties", string(out))
	out, err = it.Next()
	assert.NoError(t, err)
	assert.Equal(t, "addr:city", string(out))
	out, err = it.Next()
	assert.NoError(t, err)
	assert.Equal(t, "Maur", string(out))
	out, err = it.Next()
	assert.NoError(t, err)
	assert.Equal(t, "addr:country", string(out))
	out, err = it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, "", string(out))
	out, err = it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, "", string(out))
}
