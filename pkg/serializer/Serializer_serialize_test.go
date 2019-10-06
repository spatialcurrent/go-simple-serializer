// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package serializer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

func TestSerializerSerializeBSON(t *testing.T) {
	in := map[string]interface{}{
		"foo": "bar",
	}
	s := New(FormatBSON)
	b, err := s.Serialize(in)
	assert.NoError(t, err)
	out, err := s.Deserialize(b)
	assert.NoError(t, err)
	assert.Equal(t, in, out)
}

func TestSerializerCSVMap(t *testing.T) {
	in := map[interface{}]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	s := New(FormatCSV).Limit(NoLimit).Sorted(true)
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "a,b,c\nx,y,z\n", string(out))
}

func TestSerializerCSVSliceString(t *testing.T) {
	in := []string{"a", "b", "c"}
	s := New(FormatCSV).Limit(NoLimit)
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "a,b,c\n", string(out))
}

func TestSerializerCSVSlice(t *testing.T) {
	in := []map[string]interface{}{
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
	s := New(FormatCSV).Limit(NoLimit).Sorted(true)
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "a,b,c\n1,2,3\n4,5,6\n", string(out))
}

func TestSerializerCSVSliceExpandHeader(t *testing.T) {
	in := []map[string]interface{}{
		map[string]interface{}{
			"a": "1",
			"c": "3",
		},
		map[string]interface{}{
			"b": "5",
			"c": "6",
		},
	}
	s := New(FormatCSV).Limit(NoLimit).Sorted(true).ExpandHeader(true)
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "a,b,c\n1,,3\n,5,6\n", string(out))
}

func TestSerializerCSVSliceExpandHeaderWithWildcard(t *testing.T) {
	in := []map[string]interface{}{
		map[string]interface{}{
			"a": "1",
			"c": "3",
		},
		map[string]interface{}{
			"b": "5",
			"c": "6",
		},
	}
	s := New(FormatCSV).Limit(NoLimit).ExpandHeader(true).Header([]interface{}{"*", "c"})
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "a,b,c\n1,,3\n,5,6\n", string(out))
}

/*
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

func TestSerializerSerializeGoMap(t *testing.T) {
	in := map[interface{}]interface{}{
		"foo": "bar",
	}
	s := New(FormatGo)
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "map[interface {}]interface {}{\"foo\":\"bar\"}", string(out))
}

func TestSerializerGoStruct(t *testing.T) {
	in := struct {
		A string
		B string
		C string
	}{
		A: "1",
		B: "2",
		C: "3",
	}
	s := New(FormatGo)
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "struct { A string; B string; C string }{A:\"1\", B:\"2\", C:\"3\"}", string(out))
}

func TestSerializerSerializeJSON(t *testing.T) {
	in := map[interface{}]interface{}{
		"foo": "bar",
	}
	s := New(FormatJSON).Sorted(true)
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\"}", string(out))
}

func TestSerializerSerializeJSONL(t *testing.T) {
	in := []map[string]interface{}{
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
	s := New(FormatJSONL).Limit(NoLimit).LineSeparator("\n")
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "{\"a\":\"1\",\"b\":\"2\",\"c\":\"3\"}\n{\"a\":\"4\",\"b\":\"5\",\"c\":\"6\"}\n", string(out))
}

func TestSerializerSerializeTags(t *testing.T) {
	in := map[interface{}]interface{}{
		"hello": "beautiful world",
	}
	s := New(FormatTags).
		KeyValueSeparator("=").
		LineSeparator("\n").
		ValueSerializer(stringify.NewStringer("", false, false, false)).
		Limit(NoLimit).
		Sorted(true)
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "hello=\"beautiful world\"", string(out))
}

func TestSerializerSerializeToml(t *testing.T) {
	in := map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}
	s := New(FormatTOML)
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "a = 1.0\nb = 2.0\nc = 3.0\n", string(out))
}

func TestSerializerSerializeYaml(t *testing.T) {
	in := map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}
	s := New(FormatYAML)
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "a: 1\nb: 2\nc: 3\n", string(out))
}
