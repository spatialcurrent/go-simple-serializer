// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	A string   `map:"a"`
	B string   `map:"b,omitempty"`
	C string   `map:"-"`
	D []string `map:"d,omitempty"`
}

func TestUnmarshalMap(t *testing.T) {
	in := map[string]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	expected := &testStruct{
		A: "x",
		B: "y",
		C: "",
		D: nil,
	}
	out := &testStruct{}
	err := Unmarshal(in, out)
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestUnmarshalMapSlice(t *testing.T) {
	in := map[string]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
		"d": []string{"x", "y", "z"},
	}
	expected := &testStruct{
		A: "x",
		B: "y",
		C: "",
		D: []string{"x", "y", "z"},
	}
	out := &testStruct{}
	err := Unmarshal(in, out)
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestUnmarshalMapInterfaceSlice(t *testing.T) {
	in := map[string]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
		"d": []interface{}{"x", "y", "z"},
	}
	expected := &testStruct{
		A: "x",
		B: "y",
		C: "",
		D: []string{"x", "y", "z"},
	}
	out := &testStruct{}
	err := Unmarshal(in, out)
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}
