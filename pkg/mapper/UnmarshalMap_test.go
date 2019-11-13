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

func TestUmarshalMap(t *testing.T) {
	in := map[string]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	expected := map[string]string{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	out := map[string]string{}
	err := UnmarshalMap(in, &out)
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestUmarshalMapPointers(t *testing.T) {
	in := map[string]interface{}{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	expected := map[string]string{
		"a": "x",
		"b": "y",
		"c": "z",
	}
	out := map[string]string{}
	err := UnmarshalMap(&in, &out)
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}
