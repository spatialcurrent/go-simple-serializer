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

type testMarshaler string

func (t testMarshaler) MarshalMap() (interface{}, error) {
	return map[string]interface{}{"value": string(t)}, nil
}

func TestMarshalNil(t *testing.T) {
	out, err := Marshal(nil)
	assert.NoError(t, err)
	assert.Nil(t, out)
}

func TestMarshalString(t *testing.T) {
	out, err := Marshal("hello world")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", out)
}

func TestMarshalStringPointer(t *testing.T) {
	foo := "foo"
	out, err := Marshal(&foo)
	assert.NoError(t, err)
	assert.Equal(t, "foo", out)
}

func TestMarshalSliceString(t *testing.T) {
	out, err := Marshal([]string{"foo", "bar"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, out)
}

func TestMarshalSliceStringPointer(t *testing.T) {
	foo := "foo"
	out, err := Marshal([]*string{&foo})
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo"}, out)
}

func TestMarshalMapStringInterace(t *testing.T) {
	m, err := Marshal(map[string]interface{}{"foo": "bar"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, m)
}

func TestMarshalMapInterfaceFloat(t *testing.T) {
	m, err := Marshal(map[interface{}]float64{"yo": 1.2})
	assert.NoError(t, err)
	assert.Equal(t, map[interface{}]float64{"yo": 1.2}, m)
}

func TestMarshalMapStringString(t *testing.T) {
	m, err := Marshal(map[string]string{"foo": "bar"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"foo": "bar"}, m)
}

func TestMarshalMapStruct(t *testing.T) {
	m, err := Marshal(struct{ Foo string }{Foo: "bar"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"Foo": "bar"}, m)
}

func TestMarshalMapStructTagged(t *testing.T) {
	m, err := Marshal(struct {
		Foo string `map:"foo"`
	}{Foo: "bar"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, m)
}

func TestMarshalMapStructMarshaler(t *testing.T) {
	m, err := Marshal(testMarshaler("bar"))
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"value": "bar"}, m)
}

func TestMarshalMapStructMarshalerInside(t *testing.T) {
	m, err := Marshal(map[string]interface{}{
		"foo": testMarshaler("bar"),
	})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"foo": map[string]interface{}{"value": "bar"}}, m)
}
