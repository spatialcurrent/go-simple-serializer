// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package mapper

type testStruct struct {
	A string   `map:"a"`
	B string   `map:"b,omitempty"`
	C string   `map:"-"`
	D []string `map:"d,omitempty"`
}

type testObjects struct {
	Objects []*testStruct `map:"objects"`
}
