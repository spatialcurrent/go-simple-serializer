// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package json

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

func TestUnmarshalTypeTrue(t *testing.T) {
	obj, err := UnmarshalType([]byte("true"), reflect.TypeOf(true))
	assert.NoError(t, err)
	assert.Equal(t, obj, true)
}

func TestUnmarshalTypeFalse(t *testing.T) {
	obj, err := UnmarshalType([]byte("false"), reflect.TypeOf(true))
	assert.NoError(t, err)
	assert.Equal(t, obj, false)
}

func TestUnmarshalTypeNull(t *testing.T) {
	obj, err := UnmarshalType([]byte("null"), reflect.TypeOf(true))
	assert.NoError(t, err)
	assert.Equal(t, obj, nil)
}

func TestUnmarshalTypeArrayInterface(t *testing.T) {
	obj, err := UnmarshalType([]byte("[1,2,3]"), reflect.TypeOf([]interface{}{}))
	assert.NoError(t, err)
	assert.Equal(t, obj, []interface{}{1.0, 2.0, 3.0})
}

func TestUnmarshalTypeArrayFloat64(t *testing.T) {
	obj, err := UnmarshalType([]byte("[1,2,3]"), reflect.TypeOf([]float64{}))
	assert.NoError(t, err)
	assert.Equal(t, obj, []float64{1.0, 2.0, 3.0})
}

func TestUnmarshalTypeMapStringInterface(t *testing.T) {
	obj, err := UnmarshalType([]byte("{\"a\":1,\"b\":2,\"c\":3}"), reflect.TypeOf(map[string]interface{}{}))
	assert.NoError(t, err)
	assert.Equal(t, obj, map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0})
}

func TestUnmarshalTypeMapStringFloat64(t *testing.T) {
	obj, err := UnmarshalType([]byte("{\"a\":1,\"b\":2,\"c\":3}"), reflect.TypeOf(map[string]float64{}))
	assert.NoError(t, err)
	assert.Equal(t, obj, map[string]float64{"a": 1.0, "b": 2.0, "c": 3.0})
}

func TestUnmarshalTypeString(t *testing.T) {
	obj, err := UnmarshalType([]byte("\"hello world\""), reflect.TypeOf(""))
	assert.NoError(t, err)
	assert.Equal(t, obj, "hello world")
}

func TestUnmarshalTypeFloat64(t *testing.T) {
	obj, err := UnmarshalType([]byte("123.456"), reflect.TypeOf(123.456))
	assert.NoError(t, err)
	assert.Equal(t, obj, 123.456)
}
