// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

type ErrUnknownFormat struct {
	Name string // the name of the unknown format
}

// Error returns the error as a string.
func (e ErrUnknownFormat) Error() string {
	return "unknown format " + e.Name
}
