// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"reflect"
)

// DeserializeInput provides the input for the DeserializeBytes function.
type DeserializeInput struct {
	Bytes      []byte
	Format     string
	Header     []string
	Comment    string
	LazyQuotes bool
	SkipLines  int
	Limit      int
	Type       reflect.Type
	Async      bool
	Verbose    bool
}
