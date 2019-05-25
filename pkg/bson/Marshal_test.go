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
)

import (
	"github.com/stretchr/testify/assert"
)

func TestMarshalMap(t *testing.T) {
	in := map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0}
	b, err := Marshal(in)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	returned, err := Unmarshal(b)
	assert.NoError(t, err)
	assert.Equal(t, in, returned) // check roundtrip
}

func TestMarshalStruct(t *testing.T) {
	in := struct {
		A string
		B string
		C string
	}{A: "1", B: "2", C: "3"}
	b, err := Marshal(in)
	assert.NoError(t, err)
	returned, err := UnmarshalType(b, reflect.TypeOf(in))
	assert.NoError(t, err)
	assert.Equal(t, in, returned) // check roundtrip
}
