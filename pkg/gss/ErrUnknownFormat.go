// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

// ErrUnknownFormat is used when an unknown format is provided.
type ErrUnknownFormat struct {
	Name string // the name of the unknown format
}

// Error returns the error as a string.
func (e ErrUnknownFormat) Error() string {
	return "unknown format " + e.Name
}
