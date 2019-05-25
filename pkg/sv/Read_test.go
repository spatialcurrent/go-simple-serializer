// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"reflect"
	"strings"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	text := `
a,b,c
mary,RST,46
joe,XYZ,79
`

	expected := []interface{}{
		map[string]string{
			"a": "mary",
			"b": "RST",
			"c": "46",
		},
		map[string]string{
			"a": "joe",
			"b": "XYZ",
			"c": "79",
		},
	}

	out, err := Read(&ReadInput{
		Type:       reflect.TypeOf([]interface{}{}),
		Reader:     strings.NewReader(text),
		Separator:  ',',
		Comment:    "",
		SkipLines:  0,
		LazyQuotes: false,
	})
	assert.Nil(t, err)
	assert.Equal(t, expected, out)
}

func TestReadTyped(t *testing.T) {
	text := `
a,b,c
mary,RST,46
joe,XYZ,79
`

	expected := []map[string]interface{}{
		map[string]interface{}{
			"a": "mary",
			"b": "RST",
			"c": "46",
		},
		map[string]interface{}{
			"a": "joe",
			"b": "XYZ",
			"c": "79",
		},
	}

	out, err := Read(&ReadInput{
		Type:       reflect.TypeOf([]map[string]interface{}{}),
		Reader:     strings.NewReader(text),
		Separator:  ',',
		Comment:    "",
		SkipLines:  0,
		LazyQuotes: false,
	})
	assert.Nil(t, err)
	assert.Equal(t, expected, out)
}
