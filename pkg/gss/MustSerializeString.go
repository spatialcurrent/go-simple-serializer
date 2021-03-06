// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

// MustSerializeString serializes an object to its representation given by format and panics if there is any error.
func MustSerializeString(input *SerializeBytesInput) string {
	b, err := SerializeBytes(input)
	if err != nil {
		panic(err)
	}
	return string(b)
}
