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

func TestWriteObject(t *testing.T) {
	object := map[string]interface{}{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	keySerializer := stringify.NewStringer("", false, false, false)
	w := NewWriter(buf, "\n", keySerializer, false)
	assert.NotNil(t, w)

	err := w.WriteObject(object)
	assert.NoError(t, err)

	err = w.Flush()
	assert.NoError(t, err)

	text := buf.String()
	assert.NotNil(t, text)
	assert.Equal(t, "{\"a\":\"1\",\"b\":\"2\",\"c\":\"3\"}\n", text)
}

func TestWriterObjects(t *testing.T) {
	objects := []map[string]interface{}{
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

	buf := bytes.NewBuffer(make([]byte, 0))

	keySerializer := stringify.NewStringer("", false, false, false)
	w := NewWriter(buf, "\n", keySerializer, false)
	assert.NotNil(t, w)

	err := w.WriteObjects(objects)
	assert.NoError(t, err)

	err = w.Flush()
	assert.NoError(t, err)

	text := buf.String()
	assert.NotNil(t, text)
	assert.Equal(t, "{\"a\":\"1\",\"b\":\"2\",\"c\":\"3\"}\n{\"a\":\"4\",\"b\":\"5\",\"c\":\"6\"}\n", text)
}
