// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package sv includes iterators for reading from separated-values sources and writing using separated-values formats, such as CSV and TSV.
package sv

import (
	"reflect"
)

const (
	Wildcard = "*"
)

func concerete(object reflect.Value) reflect.Value {
	return reflect.ValueOf(object.Interface())
}
