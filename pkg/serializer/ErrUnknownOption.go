// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package serializer

// ErrUnknownOption is used when an unknown option is provided.
type ErrUnknownOption struct {
	Name string // the name of the unknown option
}

// Error returns the error formatted as a string.
func (e ErrUnknownOption) Error() string {
	return "unknown option " + e.Name
}
