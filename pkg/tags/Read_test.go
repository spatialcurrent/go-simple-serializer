// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	in := `
	a=b
	hello="beautiful world"
	hello="beautiful \"wide\" world"
  `

	expected := []interface{}{
		map[string]string{"a": "b"},
		map[string]string{"hello": "beautiful world"},
		map[string]string{"hello": "beautiful \"wide\" world"},
	}

	out, err := Read(&ReadInput{
		Type:              reflect.TypeOf([]interface{}{}),
		Reader:            strings.NewReader(in),
		SkipLines:         0,
		SkipBlanks:        true,
		SkipComments:      false,
		KeyValueSeparator: "=",
		LineSeparator:     []byte("\n")[0],
		DropCR:            true,
		Comment:           "",
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestReadBlanks(t *testing.T) {
	in := `
	a=b
	hello="beautiful world"
	hello="beautiful \"wide\" world"
  `

	expected := []interface{}{
		nil,
		map[string]string{"a": "b"},
		map[string]string{"hello": "beautiful world"},
		map[string]string{"hello": "beautiful \"wide\" world"},
		nil,
	}

	out, err := Read(&ReadInput{
		Type:              reflect.TypeOf([]interface{}{}),
		Reader:            strings.NewReader(in),
		SkipLines:         0,
		SkipBlanks:        false,
		SkipComments:      false,
		KeyValueSeparator: "=",
		LineSeparator:     []byte("\n")[0],
		DropCR:            true,
		Comment:           "",
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, out)
}
