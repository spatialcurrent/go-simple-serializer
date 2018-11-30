// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

// MustSerializeString serializes an object to its representation given by format and panics if there is any error.
func MustSerializeString(input interface{}, format string, header []string, limit int) string {
	b, err := SerializeBytes(input, format, header, limit)
	if err != nil {
		panic(err)
	}
	return string(b)
}
