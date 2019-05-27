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
	valueSerializer := stringify.DefaultValueStringer("")

	b, err := Marshal(in, valueSerializer, true)
	assert.NoError(t, err)
	assert.Equal(t, "a=1 b=2 c=3", string(b))
}

func TestMarshalStruct(t *testing.T) {
	in := struct {
		A string
		B string
		C string
	}{A: "1", B: "2", C: "3"}
	valueSerializer := stringify.DefaultValueStringer("")
	b, err := Marshal(in, valueSerializer, true)
	assert.NoError(t, err)
	assert.Equal(t, "A=1 B=2 C=3", string(b))
}
