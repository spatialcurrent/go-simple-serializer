// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortTable(t *testing.T) {
	header := []string{"a", "c", "b"}
	rows := [][]string{
		[]string{"1", "3", "2"},
		[]string{"4", "6", "5"},
	}
	newHeader, newRows := SortTable(header, rows, false)
	assert.Equal(t, []string{"a", "b", "c"}, newHeader)
	assert.Equal(t, [][]string{[]string{"1", "2", "3"}, []string{"4", "5", "6"}}, newRows)
}

func TestSortTableReverse(t *testing.T) {
	header := []string{"a", "c", "b"}
	rows := [][]string{
		[]string{"1", "3", "2"},
		[]string{"4", "6", "5"},
	}
	newHeader, newRows := SortTable(header, rows, true)
	assert.Equal(t, []string{"c", "b", "a"}, newHeader)
	assert.Equal(t, [][]string{[]string{"3", "2", "1"}, []string{"6", "5", "4"}}, newRows)
}
