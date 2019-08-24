// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gssjs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerializeJsonNil(t *testing.T) {
	str, err := Serialize(nil, "json", map[string]interface{}{})
	require.NoError(t, err)
	require.Equal(t, "null", str)
}

func TestSerializeJsonTrue(t *testing.T) {
	str, err := Serialize(true, "json", map[string]interface{}{})
	require.NoError(t, err)
	require.Equal(t, "true", str)
}

func TestSerializeJsonFalse(t *testing.T) {
	str, err := Serialize(false, "json", map[string]interface{}{})
	require.NoError(t, err)
	require.Equal(t, "false", str)
}

func TestSerializeJsonSlice(t *testing.T) {
	str, err := Serialize([]interface{}{"a", "b", "c"}, "json", map[string]interface{}{})
	require.NoError(t, err)
	require.Equal(t, "[\"a\",\"b\",\"c\"]", str)
}

func TestSerializeJsonSlicePretty(t *testing.T) {
	str, err := Serialize([]interface{}{"a", "b", "c"}, "json", map[string]interface{}{
		"pretty": true,
	})
	require.NoError(t, err)
	require.Equal(t, "[\n  \"a\",\n  \"b\",\n  \"c\"\n]", str)
}

func TestSerializeJsonMap(t *testing.T) {
	str, err := Serialize(map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}, "json", map[string]interface{}{
		"pretty": false,
	})
	require.NoError(t, err)
	require.Equal(t, "{\"a\":1,\"b\":2,\"c\":3}", str)
}

func TestSerializeJsonMapPretty(t *testing.T) {
	str, err := Serialize(map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}, "json", map[string]interface{}{
		"pretty": true,
	})
	require.NoError(t, err)
	require.Equal(t, "{\n  \"a\": 1,\n  \"b\": 2,\n  \"c\": 3\n}", str)
}

func TestSerializeJsonString(t *testing.T) {
	str, err := Serialize("hello world", "json", map[string]interface{}{})
	require.NoError(t, err)
	require.Equal(t, "\"hello world\"", str)
}

func TestSerializeJsonStringPtr(t *testing.T) {
	in := "hello world"
	str, err := Serialize(&in, "json", map[string]interface{}{})
	require.NoError(t, err)
	require.Equal(t, "\"hello world\"", str)
}

func TestSerializeJsonFloat(t *testing.T) {
	str, err := Serialize(123.456, "json", map[string]interface{}{})
	require.NoError(t, err)
	require.Equal(t, "123.456", str)
}

func TestSerializeJsonStruct(t *testing.T) {
	in := struct {
		A string
		B string
		C string
	}{
		A: "1",
		B: "2",
		C: "3",
	}
	str, err := Serialize(in, "json", map[string]interface{}{})
	require.NoError(t, err)
	require.Equal(t, "{\"A\":\"1\",\"B\":\"2\",\"C\":\"3\"}", str)
}

func TestSerializeBytesTags(t *testing.T) {
	in := map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}
	b, err := Serialize(in, "tags", map[string]interface{}{
		"sorted": true,
	})
	require.NoError(t, err)
	require.Equal(t, "a=1 b=2 c=3", string(b))
}

func TestSerializeBytesTagsMultiple(t *testing.T) {
	in := []interface{}{
		map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0},
		map[string]interface{}{"a": "x", "b": "y", "c": "z"},
	}
	b, err := Serialize(in, "tags", map[string]interface{}{
		"sorted": true,
	})
	require.NoError(t, err)
	require.Equal(t, "a=1 b=2 c=3\na=x b=y c=z\n", string(b))
}

/*

func TestSerializeBytesUnknown(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: []interface{}{"a", "b", "c"},
		Format: "foo",
	})
	assert.IsType(t, &ErrUnknownFormat{}, errors.Cause(err))
	assert.Equal(t, "could not serialize object: unknown format foo", err.Error())
	assert.Equal(t, []byte{}, b)
}

func TestSerializeBytesCsvSlice(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: []interface{}{map[string]string{"a": "x", "b": "y", "c": "z"}},
		Format: "csv",
		Limit:  NoLimit,
		Sorted: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "a,b,c\nx,y,z\n", string(b))
}

func TestSerializeBytesTsvSlice(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: []interface{}{map[string]string{"a": "x", "b": "y", "c": "z"}},
		Format: "tsv",
		Limit:  NoLimit,
		Sorted: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "a\tb\tc\nx\ty\tz\n", string(b))
}



func TestSerializeBytesTagsMultiple(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: ,
		Format:        "tags",
		LineSeparator: "\n",
		Sorted:        true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "a=1 b=2 c=3\na=x b=y c=z\n", string(b))
}

*/
