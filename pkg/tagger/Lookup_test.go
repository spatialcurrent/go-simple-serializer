// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tagger

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testLookup(t *testing.T, in string, key string, expected *Value) {
	v, err := Lookup(reflect.StructTag(in), key)
	assert.NoError(t, err)
	assert.Equal(t, expected, v)
	out, err := v.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, reflect.StructTag(in).Get(key), string(out))
}

func TestLookup(t *testing.T) {
	in := `name:"a" json:"b"`
	expected := &Value{
		Ignore:    false,
		Name:      "a",
		OmitEmpty: false,
	}
	testLookup(t, in, "name", expected)
}

func TestLookupIgnore(t *testing.T) {
	in := `name:"-" json:"b"`
	expected := &Value{
		Ignore:    true,
		Name:      "",
		OmitEmpty: false,
	}
	testLookup(t, in, "name", expected)
}

func TestLookupOmitEmpty(t *testing.T) {
	in := `name:"a,omitempty" json:"b"`
	expected := &Value{
		Ignore:    false,
		Name:      "a",
		OmitEmpty: true,
	}
	testLookup(t, in, "name", expected)
}
