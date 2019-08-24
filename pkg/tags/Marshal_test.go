// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"testing"

	"github.com/spatialcurrent/go-stringify/pkg/stringify"
	"github.com/stretchr/testify/assert"
)

func TestMarshalMap(t *testing.T) {

	in := map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}
	keys := make([]interface{}, 0)
	keySerializer := stringify.NewStringer("", false, false, false)
	valueSerializer := stringify.NewStringer("", false, false, false)

	b, err := Marshal(in, keys, true, "=", keySerializer, valueSerializer, true, false)
	assert.NoError(t, err)
	assert.Equal(t, "a=1 b=2 c=3", string(b))
}

func TestMarshalMapKeys(t *testing.T) {

	in := map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}
	keys := []interface{}{"a", "c"}
	keySerializer := stringify.NewStringer("", false, false, false)
	valueSerializer := stringify.NewStringer("", false, false, false)

	b, err := Marshal(in, keys, false, "=", keySerializer, valueSerializer, true, false)
	assert.NoError(t, err)
	assert.Equal(t, "a=1 c=3", string(b))
}

func TestMarshalMapKeysExpand(t *testing.T) {

	in := map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}
	keys := []interface{}{"a", "c"}
	keySerializer := stringify.NewStringer("", false, false, false)
	valueSerializer := stringify.NewStringer("", false, false, false)

	b, err := Marshal(in, keys, true, "=", keySerializer, valueSerializer, true, false)
	assert.NoError(t, err)
	assert.Equal(t, "a=1 c=3 b=2", string(b))
}

func TestMarshalMapUpper(t *testing.T) {

	in := map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}
	keys := make([]interface{}, 0)
	keySerializer := stringify.NewStringer("", false, false, true)
	valueSerializer := stringify.NewStringer("", false, false, false)

	b, err := Marshal(in, keys, true, "=", keySerializer, valueSerializer, true, false)
	assert.NoError(t, err)
	assert.Equal(t, "A=1 B=2 C=3", string(b))
}

func TestMarshalStruct(t *testing.T) {
	in := struct {
		A string
		B string
		C string
	}{A: "1", B: "2", C: "3"}
	keys := make([]interface{}, 0)
	keySerializer := stringify.NewStringer("", false, false, false)
	valueSerializer := stringify.NewStringer("", false, false, false)
	b, err := Marshal(in, keys, true, "=", keySerializer, valueSerializer, true, false)
	assert.NoError(t, err)
	assert.Equal(t, "A=1 B=2 C=3", string(b))
}

func TestMarshalStructKeys(t *testing.T) {
	in := struct {
		A string
		B string
		C string
	}{A: "1", B: "2", C: "3"}
	keys := []interface{}{"A", "C"}
	keySerializer := stringify.NewStringer("", false, false, false)
	valueSerializer := stringify.NewStringer("", false, false, false)
	b, err := Marshal(in, keys, false, "=", keySerializer, valueSerializer, true, false)
	assert.NoError(t, err)
	assert.Equal(t, "A=1 C=3", string(b))
}

func TestMarshalStructLower(t *testing.T) {
	in := struct {
		A string
		B string
		C string
	}{A: "1", B: "2", C: "3"}
	keys := make([]interface{}, 0)
	keySerializer := stringify.NewStringer("", false, true, false)
	valueSerializer := stringify.NewStringer("", false, false, false)
	b, err := Marshal(in, keys, true, "=", keySerializer, valueSerializer, true, false)
	assert.NoError(t, err)
	assert.Equal(t, "a=1 b=2 c=3", string(b))
}
