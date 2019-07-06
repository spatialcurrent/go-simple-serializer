// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package toml

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalEmpty(t *testing.T) {
	obj, err := Unmarshal([]byte{})
	assert.Equal(t, err, ErrEmptyInput)
	assert.Equal(t, obj, nil)
}

func TestUnmarshalMapStringInterface(t *testing.T) {
	obj, err := Unmarshal([]byte("a = 1.0\nb = 2.0\nc = 3.0\n"))
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}, obj)
}
