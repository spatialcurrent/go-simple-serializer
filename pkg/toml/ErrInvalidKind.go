// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package toml

import (
	"fmt"
	"reflect"
)

type ErrInvalidKind struct {
	Value    reflect.Type
	Expected []reflect.Kind
}

func (e ErrInvalidKind) Error() string {
	return fmt.Sprintf("type %q is of invalid kind, expecting one of %q", e.Value, e.Expected)
}
