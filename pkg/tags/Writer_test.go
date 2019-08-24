// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spatialcurrent/go-stringify/pkg/stringify"
	"github.com/stretchr/testify/assert"
)

func TestWriteObject(t *testing.T) {
	object := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	keys := make([]interface{}, 0)
	keySerializer := stringify.NewStringer("", false, false, false)
	valueSerializer := stringify.NewStringer("", false, false, false)

	w := NewWriter(buf, keys, true, "=", "\n", keySerializer, valueSerializer, true, false)
	assert.NotNil(t, w)

	err := w.WriteObject(object)
	assert.NoError(t, err)

	err = w.Flush()
	assert.NoError(t, err)

	text := buf.String()
	assert.NotNil(t, text)
	assert.Equal(t, "a=1 b=2 c=3\n", text)
}

func TestWriteStruct(t *testing.T) {
	object := struct {
		A string
		B string
		C string
	}{
		A: "1",
		B: "2",
		C: "3",
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	keys := make([]interface{}, 0)
	keySerializer := func(object interface{}) (string, error) {
		str, err := stringify.NewStringer("", false, false, false)(object)
		if err != nil {
			return str, err
		}
		return strings.ToLower(str), nil
	}
	valueSerializer := stringify.NewStringer("", false, false, false)

	w := NewWriter(buf, keys, true, "=", "\n", keySerializer, valueSerializer, true, false)
	assert.NotNil(t, w)

	err := w.WriteObject(object)
	assert.NoError(t, err)

	err = w.Flush()
	assert.NoError(t, err)

	text := buf.String()
	assert.NotNil(t, text)
	assert.Equal(t, "a=1 b=2 c=3\n", text)
}

func TestWriterObjects(t *testing.T) {
	objects := []interface{}{
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

	buf := bytes.NewBuffer(make([]byte, 0))
	keys := make([]interface{}, 0)
	keySerializer := stringify.NewStringer("", false, false, false)
	valueSerializer := stringify.NewStringer("", false, false, false)

	w := NewWriter(buf, keys, true, "=", "\n", keySerializer, valueSerializer, true, false)
	assert.NotNil(t, w)

	err := w.WriteObjects(objects)
	assert.NoError(t, err)

	err = w.Flush()
	assert.NoError(t, err)

	text := buf.String()
	assert.NotNil(t, text)
	assert.Equal(t, "a=1 b=2 c=3\na=4 b=5 c=6\n", text)
}
