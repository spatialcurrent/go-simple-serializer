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

func TestUnmarshalSlice(t *testing.T) {
	in := []interface{}{
		map[string]interface{}{
			"a": "x",
			"b": "y",
			"c": "z",
			"d": []interface{}{"x", "y", "z"},
		},
	}
	expected := []testStruct{
		testStruct{
			A: "x",
			B: "y",
			C: "",
			D: []string{"x", "y", "z"},
		},
	}
	out := []testStruct{}
	err := UnmarshalSlice(in, &out)
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestUnmarshalSlicePointers(t *testing.T) {
	in := []interface{}{
		map[string]interface{}{
			"a": "x",
			"b": "y",
			"c": "z",
			"d": []interface{}{"x", "y", "z"},
		},
	}
	expected := []*testStruct{
		&testStruct{
			A: "x",
			B: "y",
			C: "",
			D: []string{"x", "y", "z"},
		},
	}
	out := []*testStruct{}
	err := UnmarshalSlice(in, &out)
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}
