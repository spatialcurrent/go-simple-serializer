// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSerializeBytesJsonNil(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: nil,
		Format: "json",
	})
	assert.NoError(t, err)
	assert.Equal(t, "null", string(b))
}

func TestSerializeBytesJsonTrue(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: true,
		Format: "json",
	})
	assert.NoError(t, err)
	assert.Equal(t, "true", string(b))
}

func TestSerializeBytesJsonFalse(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: false,
		Format: "json",
	})
	assert.NoError(t, err)
	assert.Equal(t, string(b), "false")
}

func TestSerializeBytesJsonSlice(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: []interface{}{"a", "b", "c"},
		Format: "json",
	})
	assert.NoError(t, err)
	assert.Equal(t, "[\"a\",\"b\",\"c\"]", string(b))
}

func TestSerializeBytesJsonSlicePretty(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: []interface{}{"a", "b", "c"},
		Format: "json",
		Pretty: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "[\n  \"a\",\n  \"b\",\n  \"c\"\n]", string(b))
}

func TestSerializeBytesJsonMap(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0},
		Format: "json",
		Pretty: false,
	})
	assert.NoError(t, err)
	assert.Equal(t, "{\"a\":1,\"b\":2,\"c\":3}", string(b))
}

func TestSerializeBytesJsonMapPretty(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0},
		Format: "json",
		Pretty: true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "{\n  \"a\": 1,\n  \"b\": 2,\n  \"c\": 3\n}", string(b))
}

func TestSerializeBytesJsonString(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: "hello world",
		Format: "json",
	})
	assert.NoError(t, err)
	assert.Equal(t, "\"hello world\"", string(b))
}

func TestSerializeBytesJsonStringPtr(t *testing.T) {
	str := "hello world"
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: &str,
		Format: "json",
	})
	assert.NoError(t, err)
	assert.Equal(t, "\"hello world\"", string(b))
}

func TestSerializeBytesJsonFloat(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: 123.456,
		Format: "json",
	})
	assert.NoError(t, err)
	assert.Equal(t, "123.456", string(b))
}

func TestSerializeBytesJsonStruct(t *testing.T) {
	in := struct {
		A string
		B string
		C string
	}{
		A: "1",
		B: "2",
		C: "3",
	}
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: in,
		Format: "json",
	})
	assert.NoError(t, err)
	assert.Equal(t, "{\"A\":\"1\",\"B\":\"2\",\"C\":\"3\"}", string(b))
}

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

func TestSerializeBytesTags(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object:            map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0},
		KeyValueSeparator: "=",
		Format:            "tags",
		Sorted:            true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "a=1 b=2 c=3", string(b))
}

func TestSerializeBytesTagsMultiple(t *testing.T) {
	b, err := SerializeBytes(&SerializeBytesInput{
		Object: []interface{}{
			map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0},
			map[string]interface{}{"a": "x", "b": "y", "c": "z"},
		},
		KeyValueSeparator: "=",
		Format:            "tags",
		LineSeparator:     "\n",
		Sorted:            true,
		Limit:             -1,
	})
	assert.NoError(t, err)
	assert.Equal(t, "a=1 b=2 c=3\na=x b=y c=z\n", string(b))
}
