// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"bytes"
	//"io"
	//"strings"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestWriteHeader(t *testing.T) {

	buf := bytes.NewBuffer(make([]byte, 0))

	w := NewWriter(buf, ',', []interface{}{"a", "b", "d"}, DecimalValueSerializer(""))
	assert.NotNil(t, w)

	err := w.WriteHeader()
	assert.Nil(t, err)

	err = w.Flush()
	assert.Nil(t, err)

	text := buf.String()
	assert.NotNil(t, text)
	assert.Equal(t, "a,b,d\n", text)
}

func TestWriteObject(t *testing.T) {
	object := map[string]interface{}{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	w := NewWriter(buf, ',', []interface{}{"a", "b", "d"}, DecimalValueSerializer(""))
	assert.NotNil(t, w)

	err := w.WriteHeader()
	assert.Nil(t, err)

	err = w.WriteObject(object)
	assert.Nil(t, err)

	err = w.Flush()
	assert.Nil(t, err)

	text := buf.String()
	assert.NotNil(t, text)
	assert.Equal(t, "a,b,d\n1,2,\n", text)
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

	w := NewWriter(buf, ',', []interface{}{"a", "b", "d"}, DecimalValueSerializer(""))
	assert.NotNil(t, w)

	err := w.WriteHeader()
	assert.Nil(t, err)

	err = w.WriteObjects(objects)
	assert.Nil(t, err)

	err = w.Flush()
	assert.Nil(t, err)

	text := buf.String()
	assert.NotNil(t, text)
	assert.Equal(t, "a,b,d\n1,2,\n4,5,\n", text)
}
