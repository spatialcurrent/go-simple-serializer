// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

// SerializeString serializes an object to its representation given by format.
func SerializeString(input *SerializeInput) (string, error) {
	b, err := SerializeBytes(input)
	return string(b), err
}
