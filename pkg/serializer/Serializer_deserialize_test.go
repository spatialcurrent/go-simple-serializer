// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package serializer

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

/*
func TestSerializerCSV(t *testing.T) {
	in := map[interface{}]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	s := New(FormatCSV).ValueSerializer(stringify.NewStringer("", false, false, false))
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "a,b,c\bx,y,z", string(out))
}

func TestSerializerTSV(t *testing.T) {
	in := map[interface{}]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	s := New(FormatTSV).ValueSerializer(stringify.NewStringer("", false, false, false))
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "a\tb\tc\bx\ty\tz", string(out))
}

*/

func TestSerializerDeserializeJSON(t *testing.T) {
	in := "{\"foo\":\"bar\"}"
	s := New(FormatJSON)
	out, err := s.Deserialize([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, out)
}

func TestSerializerDeserializeJSONL(t *testing.T) {
	in := "{\"a\":\"1\",\"b\":\"2\",\"c\":\"3\"}\n{\"a\":\"4\",\"b\":\"5\",\"c\":\"6\"}\n"
	expected := []interface{}{
		map[string]interface{}{
			"a": "1",
			"b": "2",
			"c": "3",
		},
		map[string]interface{}{
			"a": "4",
			"b": "5",
			"c": "6",
		},
	}
	s := New(FormatJSONL).LineSeparator("\n")
	out, err := s.Deserialize([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestSerializerDeserializeTags(t *testing.T) {
	in := "hello=\"beautiful world\""
	s := New(FormatTags).LineSeparator("\n")
	out, err := s.Deserialize([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, []map[string]string{map[string]string{"hello": "beautiful world"}}, out)
}

func TestSerializerDeserializeToml(t *testing.T) {
	in := "a = 1.0\nb = 2.0\nc = 3.0\n"
	s := New(FormatTOML)
	out, err := s.Deserialize([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}, out)
}

func TestSerializerDeserializeYaml(t *testing.T) {
	in := "a: 1\nb: 2\nc: 3\n"
	s := New(FormatYAML)
	out, err := s.Deserialize([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"a": 1, "b": 2, "c": 3}, out)
}
