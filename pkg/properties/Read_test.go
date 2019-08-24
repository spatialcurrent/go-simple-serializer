// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package properties

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	in := `
  a=1
  b:2
  c true
  d=nil
  e=
  `

	out, err := Read(&ReadInput{
		Type:            reflect.TypeOf(map[string]string{}),
		Reader:          strings.NewReader(in),
		LineSeparator:   []byte("\n")[0],
		Comment:         "",
		Trim:            true,
		UnescapeSpace:   false,
		UnescapeEqual:   false,
		UnescapeColon:   false,
		UnescapeNewLine: false,
	})
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"a": "1", "b": "2", "c": "true", "d": "nil", "e": ""}, out)
}
