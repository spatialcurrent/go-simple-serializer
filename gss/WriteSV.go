// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"encoding/csv"
	"io"
)

func WriteSV(w io.Writer, format string, header []string, rows [][]string) error {
	csvWriter := csv.NewWriter(w)
	if format == "tsv" {
		csvWriter.Comma = '\t'
	}
	newHeader := make([]string, 0)
	for _, k := range header {
		if k != "*" {
			newHeader = append(newHeader, k)
		}
	}
	err := csvWriter.Write(newHeader)
	if err != nil {
		return err
	}
	err = csvWriter.WriteAll(rows)
	if err != nil {
		return err
	}
	return nil
}
