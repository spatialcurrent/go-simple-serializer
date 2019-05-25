// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

// FillRight fills a slice on the right with blanks strings until it has length n.
func FillRight(row []string, n int) []string {
	newRow := make([]string, 0, n)
	newRow = append(newRow, row...)
	if len(row) < n {
		newRow = append(newRow, make([]string, n-len(row))...)
	}
	return newRow
}
