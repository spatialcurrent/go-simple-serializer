// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package bson

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
	in := []byte{0x26, 0x0, 0x0, 0x0, 0x1, 0x61, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf0, 0x3f, 0x1, 0x62, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x1, 0x63, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x40, 0x0}
	expected := map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}
	obj, err := UnmarshalType(in, reflect.TypeOf(map[string]interface{}{}))
	assert.NoError(t, err)
	assert.Equal(t, expected, obj)
}

func TestUnmarshalTypeMapInterfaceInterface(t *testing.T) {
	in := []byte{0x26, 0x0, 0x0, 0x0, 0x1, 0x61, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf0, 0x3f, 0x1, 0x62, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40, 0x1, 0x63, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x40, 0x0}
	obj, err := UnmarshalType(in, reflect.TypeOf(map[interface{}]interface{}{}))
	assert.IsType(t, err, &ErrInvalidKeys{})
	assert.Nil(t, obj)
}
