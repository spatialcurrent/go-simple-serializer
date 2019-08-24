// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package yaml

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
	assert.Equal(t, []interface{}{1, 2, 3}, obj)
}

func TestUnmarshalTypeArrayFloat64(t *testing.T) {
	obj, err := UnmarshalType([]byte("[1,2,3]"), reflect.TypeOf([]float64{}))
	assert.NoError(t, err)
	assert.Equal(t, []float64{1.0, 2.0, 3.0}, obj)
}

func TestUnmarshalTypeMapStringInterface(t *testing.T) {
	obj, err := UnmarshalType([]byte("a: 1\nb: 2\nc: 3\n"), reflect.TypeOf(map[string]interface{}{}))
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"a": 1, "b": 2, "c": 3}, obj)
}

func TestUnmarshalTypeMapStringInt(t *testing.T) {
	obj, err := UnmarshalType([]byte("a: 1\nb: 2\nc: 3\n"), reflect.TypeOf(map[string]int{}))
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, obj)
}

func TestUnmarshalTypeMapStringFloat64(t *testing.T) {
	obj, err := UnmarshalType([]byte("a: 1\nb: 2\nc: 3\n"), reflect.TypeOf(map[string]float64{}))
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{"a": 1, "b": 2, "c": 3}, obj)
}

func TestUnmarshalTypeMapStringInterfaceInline(t *testing.T) {
	obj, err := UnmarshalType([]byte("{\"a\":1,\"b\":2,\"c\":3}"), reflect.TypeOf(map[string]interface{}{}))
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"a": 1, "b": 2, "c": 3}, obj)
}

func TestUnmarshalTypeMapStringIntInline(t *testing.T) {
	obj, err := UnmarshalType([]byte("{\"a\":1,\"b\":2,\"c\":3}"), reflect.TypeOf(map[string]int{}))
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, obj)
}

func TestUnmarshalTypeMapStringFloat64Inline(t *testing.T) {
	obj, err := UnmarshalType([]byte("{\"a\":1,\"b\":2,\"c\":3}"), reflect.TypeOf(map[string]float64{}))
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{"a": 1, "b": 2, "c": 3}, obj)
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
