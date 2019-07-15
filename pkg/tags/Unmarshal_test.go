// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalEmpty(t *testing.T) {
	obj, err := Unmarshal([]byte{}, '=')
	assert.Equal(t, err, ErrEmptyInput)
	assert.Equal(t, obj, nil)
}

func TestUnmarshalSingle(t *testing.T) {
	text := "a=b"
	expected := map[string]string{"a": "b"}
	obj, err := Unmarshal([]byte(text), '=')
	assert.NoError(t, err)
	assert.Equal(t, expected, obj)
}

func TestUnmarshalMap(t *testing.T) {
	text := "a=b c=\"beautiful world\" d=\"beautiful \\\"wide\\\" world\""
	expected := map[string]string{"a": "b", "c": "beautiful world", "d": "beautiful \"wide\" world"}
	obj, err := Unmarshal([]byte(text), '=')
	assert.NoError(t, err)
	assert.Equal(t, expected, obj)
}
