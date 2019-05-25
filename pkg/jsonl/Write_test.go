// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"bytes"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	in := []map[string]interface{}{
		map[string]interface{}{
			"a": "1",
			"b": "2",
			"c": "3",
		},
		map[string]interface{}{
			"a": "4",
			"b": "5",
			"c": "6",
		},
	}

	buf := new(bytes.Buffer)
	err := Write(&WriteInput{
		Writer:        buf,
		LineSeparator: []byte("\n")[0],
		Object:        in,
	})
	assert.Nil(t, err)
	out := buf.String()
	assert.Equal(t, "{\"a\":\"1\",\"b\":\"2\",\"c\":\"3\"}\n{\"a\":\"4\",\"b\":\"5\",\"c\":\"6\"}\n", out)
}
