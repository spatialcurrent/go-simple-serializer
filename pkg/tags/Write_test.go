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
	keys := make([]interface{}, 0)

	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:            buf,
		Keys:              keys,
		KeyValueSeparator: "=",
		LineSeparator:     "\n",
		Object:            in,
		KeySerializer:     stringify.NewStringer("", false, false, false),
		ValueSerializer:   stringify.NewStringer("", false, false, false),
		Sorted:            true,
		Limit:             -1,
	})
	assert.NoError(t, err)
	out := buf.String()
	assert.Equal(t, "a=1 b=2 c=3\na=4 b=5 c=6\n", out)
}

func TestWriteKeys(t *testing.T) {
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
	keys := []interface{}{"a", "c"}

	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:            buf,
		Keys:              keys,
		KeyValueSeparator: "=",
		LineSeparator:     "\n",
		Object:            in,
		KeySerializer:     stringify.NewStringer("", false, false, false),
		ValueSerializer:   stringify.NewStringer("", false, false, false),
		Sorted:            true,
		Limit:             -1,
	})
	assert.NoError(t, err)
	out := buf.String()
	assert.Equal(t, "a=1 c=3\na=4 c=6\n", out)
}
