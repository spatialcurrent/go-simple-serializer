// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"bytes"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

func TestWrite(t *testing.T) {
	in := []interface{}{
		map[string]interface{}{
			"a": 1,
			"b": 2,
			"c": 3,
		},
		map[string]interface{}{
			"a": 4,
			"b": 5,
			"c": 6,
		},
	}

	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:          buf,
		LineSeparator:   "\n",
		Object:          in,
		ValueSerializer: stringify.DefaultValueStringer(""),
		Sorted:          true,
	})
	assert.NoError(t, err)
	out := buf.String()
	assert.Equal(t, "a=1 b=2 c=3\na=4 b=5 c=6\n", out)
}
