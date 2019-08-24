// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteTable(t *testing.T) {
	header := []string{"a", "b", "c"}
	rows := [][]string{
		[]string{"x", "y", "z"},
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	err := WriteTable(&WriteTableInput{
		Writer:    buf,
		Separator: ',',
		Header:    header,
		Rows:      rows,
		Sorted:    false,
		Reversed:  false,
	})
	assert.NoError(t, err)
	text := buf.String()
	assert.Equal(t, "a,b,c\nx,y,z\n", text)
}

func TestWriteTableSorted(t *testing.T) {
	header := []string{"a", "b", "c"}
	rows := [][]string{
		[]string{"x", "y", "z"},
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	err := WriteTable(&WriteTableInput{
		Writer:    buf,
		Separator: ',',
		Header:    header,
		Rows:      rows,
		Sorted:    true,
		Reversed:  false,
	})
	assert.NoError(t, err)
	text := buf.String()
	assert.Equal(t, "a,b,c\nx,y,z\n", text)
}

func TestWriteTableSortedReversed(t *testing.T) {
	header := []string{"a", "b", "c"}
	rows := [][]string{
		[]string{"x", "y", "z"},
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	err := WriteTable(&WriteTableInput{
		Writer:    buf,
		Separator: ',',
		Header:    header,
		Rows:      rows,
		Sorted:    true,
		Reversed:  true,
	})
	assert.NoError(t, err)
	text := buf.String()
	assert.Equal(t, "c,b,a\nz,y,x\n", text)
}
