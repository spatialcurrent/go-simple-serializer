// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

// SortTable sorts the table columns alphabetically.
// If reversed is true, then sorts in reverse alphabetical order.
func SortTable(header []string, rows [][]string, reversed bool) ([]string, [][]string) {
	newHeader := Row(append(make([]string, 0), header...))
	newHeader.Sort(reversed)
	transform := map[int]int{}
	for i, str := range header {
		transform[i] = newHeader.Index(str)
	}
	newRows := make([][]string, 0, len(rows))
	for _, row := range rows {
		newRow := Row(row).FillRight(len(newHeader)).Transform(transform)
		newRows = append(newRows, newRow)
	}
	return newHeader, newRows
}
