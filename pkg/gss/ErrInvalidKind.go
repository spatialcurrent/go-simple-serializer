// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"fmt"
	"reflect"
)

type ErrInvalidKind struct {
	Value reflect.Kind
	Valid []reflect.Kind
}

func (e ErrInvalidKind) Error() string {
	return fmt.Sprintf("invalid reflect.Kind %q, expecting one of %q", e.Value, e.Valid)
}
