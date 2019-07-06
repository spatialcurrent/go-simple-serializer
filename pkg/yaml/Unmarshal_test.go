// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package yaml

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalEmpty(t *testing.T) {
	obj, err := Unmarshal([]byte{})
	assert.Equal(t, err, ErrEmptyInput)
	assert.Equal(t, obj, nil)
}

func TestUnmarshalTrue(t *testing.T) {
	obj, err := Unmarshal([]byte("true"))
	assert.NoError(t, err)
	assert.Equal(t, obj, true)
}

func TestUnmarshalFalse(t *testing.T) {
	obj, err := Unmarshal([]byte("false"))
	assert.NoError(t, err)
	assert.Equal(t, obj, false)
}

func TestUnmarshalNull(t *testing.T) {
	obj, err := Unmarshal([]byte("null"))
	assert.NoError(t, err)
	assert.Equal(t, obj, nil)
}

func TestUnmarshalArray(t *testing.T) {
	obj, err := Unmarshal([]byte("[1,2,3]"))
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{1, 2, 3}, obj)
}

func TestUnmarshalMap(t *testing.T) {
	obj, err := Unmarshal([]byte("{\"a\":1,\"b\":2,\"c\":3}"))
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"a": 1, "b": 2, "c": 3}, obj)
}

func TestUnmarshalString(t *testing.T) {
	obj, err := Unmarshal([]byte("hello world"))
	assert.NoError(t, err)
	assert.Equal(t, obj, "hello world")
}

func TestUnmarshalStringQuoted(t *testing.T) {
	obj, err := Unmarshal([]byte("\"hello world\""))
	assert.NoError(t, err)
	assert.Equal(t, obj, "hello world")
}

func TestUnmarshalFloat(t *testing.T) {
	obj, err := Unmarshal([]byte("123.456"))
	assert.NoError(t, err)
	assert.Equal(t, obj, 123.456)
}
