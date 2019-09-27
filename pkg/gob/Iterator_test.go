// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gob

import (
	"bytes"
	"encoding/gob"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIterator(t *testing.T) {
	objects := []interface{}{
		map[string]string{"a": "x"},
		map[string]string{"b": "y"},
		map[string]string{"c": "z"},
	}

	buf := new(bytes.Buffer)

	w := gob.NewEncoder(buf)

	err := w.Encode(objects[0])
	require.NoError(t, err)
	err = w.Encode(objects[1])
	require.NoError(t, err)
	err = w.Encode(objects[2])
	require.NoError(t, err)

	it := NewIterator(&NewIteratorInput{
		Reader: buf,
		Type:   reflect.TypeOf(map[string]string{}),
		Limit:  -1,
	})

	obj, err := it.Next()
	require.NoError(t, err)
	require.NotNil(t, obj)
	require.Equal(t, objects[0], obj)

	obj, err = it.Next()
	require.NoError(t, err)
	require.NotNil(t, obj)
	require.Equal(t, objects[1], obj)

	obj, err = it.Next()
	require.NoError(t, err)
	require.NotNil(t, obj)
	require.Equal(t, objects[2], obj)

	// Should return io.EOF to indicate the reader is finished
	obj, err = it.Next()
	require.Equal(t, io.EOF, err)
	require.Nil(t, obj)

	// Should still return io.EOF to indicate the reader is finished
	obj, err = it.Next()
	require.Equal(t, io.EOF, err)
	require.Nil(t, obj)
}

func TestIteratorEmpty(t *testing.T) {
	it := NewIterator(&NewIteratorInput{
		Reader: new(bytes.Buffer),
		Type:   reflect.TypeOf(map[string]string{}),
		Limit:  -1,
	})

	// Should return io.EOF to indicate the reader is finished
	obj, err := it.Next()
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, obj)
}
