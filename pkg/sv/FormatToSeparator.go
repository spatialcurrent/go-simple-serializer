// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

// FormatToSeparator converts a format into its corresponding separator.
func FormatToSeparator(format string) (rune, error) {
	if format == "tsv" {
		return '\t', nil
	} else if format == "csv" {
		return ',', nil
	}
	return ',', &ErrInvalidFormat{Format: format}
}
