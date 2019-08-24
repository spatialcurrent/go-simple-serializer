// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package bson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalEmpty(t *testing.T) {
	obj, err := Unmarshal([]byte{})
	assert.Equal(t, err, ErrEmptyInput)
	assert.Equal(t, obj, nil)
}

func TestUnmarshalMap(t *testing.T) {
	in := []byte{0x26, 0x0, 0x0, 0x0, 0x1, 0x61, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf0, 0x3f, 0x1, 0x62, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x1, 0x63, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x40, 0x0}
	expected := map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}
	obj, err := Unmarshal(in)
	assert.NoError(t, err)
	assert.Equal(t, expected, obj)
}
