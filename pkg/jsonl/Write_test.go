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

	"github.com/stretchr/testify/assert"

	"github.com/spatialcurrent/go-stringify/pkg/stringify"
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
		LineSeparator: "\n",
		KeySerializer: stringify.NewStringer("", false, false, false),
		Object:        in,
		Pretty:        false,
		Limit:         -1,
	})
	assert.NoError(t, err)
	out := buf.String()
	assert.Equal(t, "{\"a\":\"1\",\"b\":\"2\",\"c\":\"3\"}\n{\"a\":\"4\",\"b\":\"5\",\"c\":\"6\"}\n", out)
}

func TestWriteUpper(t *testing.T) {
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
		LineSeparator: "\n",
		KeySerializer: stringify.NewStringer("", false, false, true),
		Object:        in,
		Pretty:        false,
		Limit:         -1,
	})
	assert.NoError(t, err)
	out := buf.String()
	assert.Equal(t, "{\"A\":\"1\",\"B\":\"2\",\"C\":\"3\"}\n{\"A\":\"4\",\"B\":\"5\",\"C\":\"6\"}\n", out)
}

func TestWritePretty(t *testing.T) {
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
		LineSeparator: "\n",
		KeySerializer: stringify.NewStringer("", false, false, false),
		Object:        in,
		Pretty:        true,
		Limit:         -1,
	})
	assert.NoError(t, err)
	out := buf.String()
	assert.Equal(t, "{\n  \"a\": \"1\",\n  \"b\": \"2\",\n  \"c\": \"3\"\n}\n{\n  \"a\": \"4\",\n  \"b\": \"5\",\n  \"c\": \"6\"\n}\n", out)
}
