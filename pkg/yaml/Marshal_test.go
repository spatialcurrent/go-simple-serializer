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

func TestMarshalNil(t *testing.T) {
	b, err := Marshal(nil)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "null\n")
}

func TestMarshalTrue(t *testing.T) {
	b, err := Marshal(true)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "true\n")
}

func TestMarshalFalse(t *testing.T) {
	b, err := Marshal(false)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "false\n")
}

func TestMarshalArray(t *testing.T) {
	b, err := Marshal([]interface{}{"a", "b", "c"})
	assert.NoError(t, err)
	assert.Equal(t, "- a\n- b\n- c\n", string(b))
}

func TestMarshalMap(t *testing.T) {
	b, err := Marshal(map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0})
	assert.NoError(t, err)
	assert.Equal(t, "a: 1\nb: 2\nc: 3\n", string(b))
}

func TestMarshalString(t *testing.T) {
	b, err := Marshal("hello world")
	assert.NoError(t, err)
	assert.Equal(t, "hello world\n", string(b))
}

func TestMarshalStringPtr(t *testing.T) {
	str := "hello world"
	b, err := Marshal(&str)
	assert.NoError(t, err)
	assert.Equal(t, "hello world\n", string(b))
}

func TestMarshalFloat(t *testing.T) {
	b, err := Marshal(123.456)
	assert.NoError(t, err)
	assert.Equal(t, "123.456\n", string(b))
}

func TestMarshalStruct(t *testing.T) {
	b, err := Marshal(struct {
		A string
		B string
		C string
	}{A: "1", B: "2", C: "3"})
	assert.NoError(t, err)
	assert.Equal(t, "a: \"1\"\nb: \"2\"\nc: \"3\"\n", string(b))
}
