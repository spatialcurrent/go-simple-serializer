// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

// SerializeString serializes an object to its representation given by format.
func SerializeString(input interface{}, format string, header []string, limit int) (string, error) {
	b, err := SerializeBytes(input, format, header, limit)
	return string(b), err
}
