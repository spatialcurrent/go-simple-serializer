// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package splitter

// DropCarriageReturn drops the final byte of the slice if it is a carriage return.
//
// Examples:
//	 - DropCarriageReturn([]byte("abc\r")) == []byte("abc")
//	 - DropCarriageReturn([]byte("abc")) == []byte("abc")
func DropCarriageReturn(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
