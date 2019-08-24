// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package iterator

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIteratorJsonl(t *testing.T) {
	text := `
  {"a": "b"}
  {"c": "d"}
  {"e": "f"}
  false
  true
  "foo"
  "bar"
  `

	it, err := NewIterator(&NewIteratorInput{
		Reader:        strings.NewReader(text),
		Format:        "jsonl",
		SkipLines:     0,
		Comment:       "",
		Trim:          true,
		SkipBlanks:    false,
		SkipComments:  false,
		LineSeparator: []byte("\n")[0],
		DropCR:        true,
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

func TestIteratorTags(t *testing.T) {
	text := `
  a=b x=y
  c=d y=z
  e=f h=i
  `

	it, err := NewIterator(&NewIteratorInput{
		Reader:            strings.NewReader(text),
		Type:              reflect.TypeOf([]map[string]interface{}{}),
		Format:            "tags",
		SkipLines:         0,
		Comment:           "",
		Trim:              true,
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
	assert.Equal(t, map[string]interface{}{"a": "b", "x": "y"}, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]interface{}{"c": "d", "y": "z"}, obj)

	obj, err = it.Next()
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, map[string]interface{}{"e": "f", "h": "i"}, obj)

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
