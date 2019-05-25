// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	text := `
a,b,c
mary,RST,46
joe,XYZ,79
`

	it, err := NewIterator(&NewIteratorInput{
		Reader:     strings.NewReader(text),
		Type:       reflect.TypeOf(map[string]string{}),
		Separator:  ',',
		Comment:    "",
		SkipLines:  0,
		LazyQuotes: false,
	})
	assert.NoError(t, err)

	// First Object
	obj, err := it.Next()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"a": "mary", "b": "RST", "c": "46"}, obj)

	// Second Object
	obj, err = it.Next()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"a": "joe", "b": "XYZ", "c": "79"}, obj)

	// Should return io.EOF to indicate the reader is finished
	obj, err = it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)

	// Should still return io.EOF to indicate the reader is finished
	obj, err = it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)
}

func TestIteratorHeaderSkip(t *testing.T) {
	text := `
a,b,c
mary,RST,46
joe,XYZ,79
`

	it, err := NewIterator(&NewIteratorInput{
		Reader:     strings.NewReader(text),
		Type:       reflect.TypeOf(map[string]string{}),
		Separator:  ',',
		Comment:    "",
		SkipLines:  1, // skip the first line: a,b,c
		LazyQuotes: false,
		Header:     []string{"name", "type", "age"},
	})
	assert.NoError(t, err)

	// First Object
	obj, err := it.Next()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"name": "mary", "type": "RST", "age": "46"}, obj)

	// Second Object
	obj, err = it.Next()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"name": "joe", "type": "XYZ", "age": "79"}, obj)

	// Should return io.EOF to indicate the reader is finished
	obj, err = it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)

	// Should still return io.EOF to indicate the reader is finished
	obj, err = it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)
}
