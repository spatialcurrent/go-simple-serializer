// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package fmt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteObject(t *testing.T) {
	object := map[string]interface{}{
		"a": 1,
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	w := NewWriter(buf, "%#v", "\n")
	assert.NotNil(t, w)

	err := w.WriteObject(object)
	assert.NoError(t, err)

	err = w.Flush()
	assert.NoError(t, err)

	text := buf.String()
	assert.NotNil(t, text)
	assert.Equal(t, "map[string]interface {}{\"a\":1}\n", text)
}

func TestWriteStruct(t *testing.T) {
	object := struct {
		A string
	}{
		A: "1",
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	w := NewWriter(buf, "%#v", "\n")
	assert.NotNil(t, w)

	err := w.WriteObject(object)
	assert.NoError(t, err)

	err = w.Flush()
	assert.NoError(t, err)

	text := buf.String()
	assert.NotNil(t, text)
	assert.Equal(t, "struct { A string }{A:\"1\"}\n", text)
}

func TestWriterObjects(t *testing.T) {
	objects := []interface{}{
		map[string]interface{}{
			"a": 1,
		},
		map[string]interface{}{
			"b": 5,
		},
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	w := NewWriter(buf, "%#v", "\n")
	assert.NotNil(t, w)

	err := w.WriteObjects(objects)
	assert.NoError(t, err)

	err = w.Flush()
	assert.NoError(t, err)

	text := buf.String()
	assert.NotNil(t, text)
	assert.Equal(t, "map[string]interface {}{\"a\":1}\nmap[string]interface {}{\"b\":5}\n", text)
}
