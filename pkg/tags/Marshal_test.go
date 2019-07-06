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

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

func TestMarshalMap(t *testing.T) {

	in := map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}
	keySerializer := stringify.NewStringer("", false, false, false)
	valueSerializer := stringify.NewStringer("", false, false, false)

	b, err := Marshal(in, keySerializer, valueSerializer, true, false)
	assert.NoError(t, err)
	assert.Equal(t, "a=1 b=2 c=3", string(b))
}

func TestMarshalMapUpper(t *testing.T) {

	in := map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}
	keySerializer := stringify.NewStringer("", false, false, true)
	valueSerializer := stringify.NewStringer("", false, false, false)

	b, err := Marshal(in, keySerializer, valueSerializer, true, false)
	assert.NoError(t, err)
	assert.Equal(t, "A=1 B=2 C=3", string(b))
}

func TestMarshalStruct(t *testing.T) {
	in := struct {
		A string
		B string
		C string
	}{A: "1", B: "2", C: "3"}
	keySerializer := stringify.NewStringer("", false, false, false)
	valueSerializer := stringify.NewStringer("", false, false, false)
	b, err := Marshal(in, keySerializer, valueSerializer, true, false)
	assert.NoError(t, err)
	assert.Equal(t, "A=1 B=2 C=3", string(b))
}

func TestMarshalStructLower(t *testing.T) {
	in := struct {
		A string
		B string
		C string
	}{A: "1", B: "2", C: "3"}
	keySerializer := stringify.NewStringer("", false, true, false)
	valueSerializer := stringify.NewStringer("", false, false, false)
	b, err := Marshal(in, keySerializer, valueSerializer, true, false)
	assert.NoError(t, err)
	assert.Equal(t, "a=1 b=2 c=3", string(b))
}
