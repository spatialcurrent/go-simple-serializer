// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

// ConvertInput provides the input for the Convert function.
type ConvertInput struct {
	InputBytes      []byte
	InputFormat     string
	InputHeader     []string
	InputComment    string
	InputLazyQuotes bool
	InputSkipLines  int
	InputLimit      int
	OutputFormat    string
	OutputHeader    []string
	OutputLimit     int
	Async           bool
	Verbose         bool
}
