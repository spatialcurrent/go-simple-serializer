// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"encoding/csv"
	"io"
)

// WriteTableInput provides the input for the WriteTable function.
type WriteTableInput struct {
	Writer    io.Writer  // the underlying writer
	Separator rune       // the values separator
	Header    []string   // the table header
	Rows      [][]string // the row of values to write to the underlying writer
	Sorted    bool       // sort the columns
	Reversed  bool       // if sorted, sort in reverse order
}

// WriteTable writes the given rows as separated values.
func WriteTable(input *WriteTableInput) error {

	header := input.Header
	rows := input.Rows

	if input.Sorted {
		header, rows = SortTable(header, rows, input.Reversed) // Also requires filling in all the rows with missing values
	}

	// Create a new CSV writer.
	csvWriter := csv.NewWriter(input.Writer)

	// set the values separator
	csvWriter.Comma = input.Separator

	if len(header) > 0 {
		// Write the header to the underlying writer.
		errHeader := csvWriter.Write(header)
		if errHeader != nil {
			return errHeader
		}
	}

	// WriteAll will write all the rows to the underlying io.Writer and then flushes.
	errRows := csvWriter.WriteAll(rows)
	if errRows != nil {
		return errRows
	}

	return nil
}
