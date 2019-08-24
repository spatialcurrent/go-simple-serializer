// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package toml

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalTypeEmpty(t *testing.T) {
	obj, err := UnmarshalType([]byte{}, reflect.TypeOf(map[string]string{}))
	assert.Equal(t, err, ErrEmptyInput)
	assert.Equal(t, obj, nil)
}

func TestUnmarshalTypeMapStringInterface(t *testing.T) {
	obj, err := UnmarshalType([]byte("a = 1.0\nb = 2.0\nc = 3.0\n"), reflect.TypeOf(map[string]interface{}{}))
	assert.NoError(t, err)
	assert.Equal(t, obj, map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0})
}

func TestUnmarshalTypeMapStringFloat64(t *testing.T) {
	obj, err := UnmarshalType([]byte("a = 1.0\nb = 2.0\nc = 3.0\n"), reflect.TypeOf(map[string]float64{}))
	assert.NoError(t, err)
	assert.Equal(t, obj, map[string]float64{"a": 1.0, "b": 2.0, "c": 3.0})
}
