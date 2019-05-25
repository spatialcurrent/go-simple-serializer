// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"reflect"
	"strings"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	in := `
	{"a": "b"}
  {"c": "d"}
  {"e": "f"}
  false
  true
  "foo"
  "bar"
  `

	expected := []interface{}{
		map[string]interface{}{"a": "b"},
		map[string]interface{}{"c": "d"},
		map[string]interface{}{"e": "f"},
		false,
		true,
		"foo",
		"bar",
	}

	out, err := Read(&ReadInput{
		Type:          reflect.TypeOf([]interface{}{}),
		Reader:        strings.NewReader(in),
		SkipLines:     0,
		SkipBlanks:    true,
		SkipComments:  false,
		LineSeparator: []byte("\n")[0],
		DropCR:        true,
		Comment:       "",
		Trim:          true,
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestReadBlanks(t *testing.T) {
	in := `
	{"a": "b"}
  {"c": "d"}
  {"e": "f"}
  false
  true
  "foo"
  "bar"
  `

	expected := []interface{}{
		nil,
		map[string]interface{}{"a": "b"},
		map[string]interface{}{"c": "d"},
		map[string]interface{}{"e": "f"},
		false,
		true,
		"foo",
		"bar",
		nil,
	}

	out, err := Read(&ReadInput{
		Type:          reflect.TypeOf([]interface{}{}),
		Reader:        strings.NewReader(in),
		SkipLines:     0,
		SkipBlanks:    false,
		SkipComments:  false,
		LineSeparator: []byte("\n")[0],
		DropCR:        true,
		Comment:       "",
		Trim:          true,
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}
