// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testValue(t *testing.T, in string, expected *Value) {
	v := &Value{}
	err := v.UnmarshalText([]byte(in))
	assert.NoError(t, err)
	assert.Equal(t, expected, v)
	out, err := v.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, in, string(out))
}

func TestEmptyString(t *testing.T) {
	testValue(t, "", &Value{
		Ignore:    false,
		Name:      "",
		OmitEmpty: false,
	})
}

func TestOmitEmpty(t *testing.T) {
	testValue(t, ",omitempty", &Value{
		Ignore:    false,
		Name:      "",
		OmitEmpty: true,
	})
}

func TestIgnore(t *testing.T) {
	testValue(t, "-", &Value{
		Ignore:    true,
		Name:      "",
		OmitEmpty: false,
	})
}

func TestName(t *testing.T) {
	testValue(t, "foo", &Value{
		Ignore:    false,
		Name:      "foo",
		OmitEmpty: false,
	})
}

func TestNameOmitEmpty(t *testing.T) {
	testValue(t, "foo,omitempty", &Value{
		Ignore:    false,
		Name:      "foo",
		OmitEmpty: true,
	})
}
