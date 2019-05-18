// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"encoding/csv"
)

// WriteSV writes the given rows as separated values.
func Write(input *WriteInput) error {

	// Create a new CSV writer.
	csvWriter := csv.NewWriter(input.Writer)

	// set the values separator
	csvWriter.Comma = input.Separator

	// Write the header to the underlying writer.
	err := csvWriter.Write(input.Header)
	if err != nil {
		return err
	}

	// WriteAll will write all the rows to the underlying io.Writer and then flushes.
	err = csvWriter.WriteAll(input.Rows)
	if err != nil {
		return err
	}

	return nil
}
