// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIterator(t *testing.T) {
	text := `
	a=b
	hello="beautiful world"
	hello="beautiful \"wide\" world"
  `

	it, err := NewIterator(&NewIteratorInput{
		Reader:            strings.NewReader(text),
		SkipLines:         0,
		Comment:           "",
		SkipBlanks:        false,
		SkipComments:      false,
		KeyValueSeparator: "=",
		LineSeparator:     []byte("\n")[0],
		DropCR:            true,
	})
	require.NoError(t, err)
	require.NotNil(t, it)

	// Empty Line
	obj, err := it.Next()
	assert.NoError(t, err)
	assert.Nil(t, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]string{"a": "b"}, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]string{"hello": "beautiful world"}, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]string{"hello": "beautiful \"wide\" world"}, obj)

	// Empty line returns nil object
	obj, err = it.Next()
	assert.NoError(t, err)
	assert.Nil(t, obj)

	// Should return io.EOF to indicate the reader is finished
	obj, err = it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)

	// Should still return io.EOF to indicate the reader is finished
	obj, err = it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)
}

func TestIteratorType(t *testing.T) {
	text := `
	a=b
	hello="beautiful world"
	hello="beautiful \"wide\" world"
  `

	it, err := NewIterator(&NewIteratorInput{
		Reader:            strings.NewReader(text),
		Type:              reflect.TypeOf(map[string]interface{}{}),
		SkipLines:         0,
		Comment:           "",
		SkipBlanks:        false,
		SkipComments:      false,
		KeyValueSeparator: "=",
		LineSeparator:     []byte("\n")[0],
		DropCR:            true,
	})
	require.NoError(t, err)
	require.NotNil(t, it)

	// Empty Line
	obj, err := it.Next()
	assert.NoError(t, err)
	assert.Nil(t, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]interface{}{"a": "b"}, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]interface{}{"hello": "beautiful world"}, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]interface{}{"hello": "beautiful \"wide\" world"}, obj)

	// Empty line returns nil object
	obj, err = it.Next()
	assert.NoError(t, err)
	assert.Nil(t, obj)

	// Should return io.EOF to indicate the reader is finished
	obj, err = it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)

	// Should still return io.EOF to indicate the reader is finished
	obj, err = it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)

}

func TestIteratorComment(t *testing.T) {
	text := `
  a=b
  #c=d
  e=f
  `

	it, err := NewIterator(&NewIteratorInput{
		Reader:            strings.NewReader(text),
		SkipLines:         0,
		Comment:           "#",
		SkipBlanks:        false,
		KeyValueSeparator: "=",
		LineSeparator:     []byte("\n")[0],
		DropCR:            true,
	})
	require.NoError(t, err)
	require.NotNil(t, it)

	// Empty Line
	obj, err := it.Next()
	assert.NoError(t, err)
	assert.Nil(t, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]string{"a": "b"}, obj)

	// Commented line returns nil object
	obj, err = it.Next()
	assert.NoError(t, err)
	assert.Nil(t, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]string{"e": "f"}, obj)

	// Empty Line
	obj, err = it.Next()
	assert.NoError(t, err)
	assert.Nil(t, obj)

	// Should return io.EOF to indicate the reader is finished
	obj, err = it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)
}

func TestIteratorEmpty(t *testing.T) {
	it, err := NewIterator(&NewIteratorInput{
		Reader:            strings.NewReader(""),
		SkipLines:         0,
		Comment:           "#",
		SkipBlanks:        false,
		KeyValueSeparator: "=",
		LineSeparator:     []byte("\n")[0],
		DropCR:            true,
	})
	require.NoError(t, err)
	require.NotNil(t, it)

	// Should return io.EOF to indicate the reader is finished
	obj, err := it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)
}

func TestIteratorBlanks(t *testing.T) {
	it, err := NewIterator(&NewIteratorInput{
		Reader:            strings.NewReader(strings.Repeat("\n", 5)),
		SkipLines:         0,
		Comment:           "#",
		SkipBlanks:        true,
		KeyValueSeparator: "=",
		LineSeparator:     []byte("\n")[0],
		DropCR:            true,
	})

	require.NoError(t, err)
	require.NotNil(t, it)

	// Should return io.EOF to indicate the reader is finished
	obj, err := it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)
}
