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

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

func TestWriteHeader(t *testing.T) {

	buf := bytes.NewBuffer(make([]byte, 0))

	w := NewWriter(buf, ',', []interface{}{"a", "b", "d"}, stringify.DecimalValueStringer(""))
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

	w := NewWriter(buf, ',', []interface{}{"a", "b", "d"}, stringify.DecimalValueStringer(""))
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
	objects := []interface{}{
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

	w := NewWriter(buf, ',', []interface{}{"a", "b", "d"}, stringify.DecimalValueStringer(""))
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

func TestWriteStructs(t *testing.T) {
	objects := []interface{}{
		struct {
			A string
			B string
			C string
		}{
			A: "1",
			B: "2",
			C: "3",
		},
		struct {
			A string
			B string
			C string
		}{
			A: "4",
			B: "5",
			C: "6",
		},
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	w := NewWriter(buf, ',', []interface{}{"a", "b", "d"}, stringify.DecimalValueStringer(""))
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
