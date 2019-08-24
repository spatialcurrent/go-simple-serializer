// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalNil(t *testing.T) {
	b, err := Marshal(nil, false)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "null")
}

func TestMarshalTrue(t *testing.T) {
	b, err := Marshal(true, false)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "true")
}

func TestMarshalFalse(t *testing.T) {
	b, err := Marshal(false, false)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "false")
}

func TestMarshalArray(t *testing.T) {
	b, err := Marshal([]interface{}{"a", "b", "c"}, false)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "[\"a\",\"b\",\"c\"]")
}

func TestMarshalArrayPretty(t *testing.T) {
	b, err := Marshal([]interface{}{"a", "b", "c"}, true)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "[\n  \"a\",\n  \"b\",\n  \"c\"\n]")
}

func TestMarshalMap(t *testing.T) {
	b, err := Marshal(map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}, false)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "{\"a\":1,\"b\":2,\"c\":3}")
}

func TestMarshalString(t *testing.T) {
	b, err := Marshal("hello world", false)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "\"hello world\"")
}

func TestMarshalStringPtr(t *testing.T) {
	str := "hello world"
	b, err := Marshal(&str, false)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "\"hello world\"")
}

func TestMarshalFloat(t *testing.T) {
	b, err := Marshal(123.456, false)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "123.456")
}

func TestMarshalStruct(t *testing.T) {
	b, err := Marshal(struct {
		A string
		B string
		C string
	}{A: "1", B: "2", C: "3"}, false)
	assert.NoError(t, err)
	assert.Equal(t, string(b), "{\"A\":\"1\",\"B\":\"2\",\"C\":\"3\"}")
}
