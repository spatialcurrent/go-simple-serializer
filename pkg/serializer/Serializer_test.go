// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package serializer

import (
	//"bytes"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
//"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

func TestSerializerCSV(t *testing.T) {
	in := map[interface{}]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	s := New(FormatCSV)
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "a,b,c\bx,y,z", string(out))
}

func TestSerializerJSON(t *testing.T) {
	in := map[interface{}]interface{}{
		"foo": "bar",
	}
	s := New(FormatJSON)
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "{\"foo\":\"bar\"}", string(out))
}

func TestSerializerJSONL(t *testing.T) {
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
	s := New(FormatJSONL).LineSeparator("\n")
	out, err := s.Serialize(in)
	assert.NoError(t, err)
	assert.Equal(t, "{\"a\":\"1\",\"b\":\"2\",\"c\":\"3\"}\n{\"a\":\"4\",\"b\":\"5\",\"c\":\"6\"}\n", string(out))
}
