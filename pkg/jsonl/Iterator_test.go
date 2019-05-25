// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"io"
	"strings"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	text := `
  {"a": "b"}
  {"c": "d"}
  {"e": "f"}
  false
  true
  "foo"
  "bar"
  `

	it := NewIterator(&NewIteratorInput{
		Reader:        strings.NewReader(text),
		SkipLines:     0,
		Comment:       "",
		Trim:          true,
		SkipBlanks:    false,
		SkipComments:  false,
		LineSeparator: []byte("\n")[0],
		DropCR:        true,
	})

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
	assert.Equal(t, map[string]interface{}{"c": "d"}, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]interface{}{"e": "f"}, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, false, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, true, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, "foo", obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, "bar", obj)

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
  {"a": "b"}
  #{"c": "d"}
  {"e": "f"}
  `

	it := NewIterator(&NewIteratorInput{
		Reader:        strings.NewReader(text),
		SkipLines:     0,
		Comment:       "#",
		Trim:          true,
		SkipBlanks:    false,
		LineSeparator: []byte("\n")[0],
		DropCR:        true,
	})

	// Empty Line
	obj, err := it.Next()
	assert.NoError(t, err)
	assert.Nil(t, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]interface{}{"a": "b"}, obj)

	// Commented line returns nil object
	obj, err = it.Next()
	assert.NoError(t, err)
	assert.Nil(t, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]interface{}{"e": "f"}, obj)

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
	it := NewIterator(&NewIteratorInput{
		Reader:        strings.NewReader(""),
		SkipLines:     0,
		Comment:       "#",
		Trim:          true,
		SkipBlanks:    false,
		LineSeparator: []byte("\n")[0],
		DropCR:        true,
	})

	// Should return io.EOF to indicate the reader is finished
	obj, err := it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)
}

func TestIteratorBlanks(t *testing.T) {
	it := NewIterator(&NewIteratorInput{
		Reader:        strings.NewReader(strings.Repeat("\n", 5)),
		SkipLines:     0,
		Comment:       "#",
		Trim:          true,
		SkipBlanks:    true,
		LineSeparator: []byte("\n")[0],
		DropCR:        true,
	})

	// Should return io.EOF to indicate the reader is finished
	obj, err := it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)
}
