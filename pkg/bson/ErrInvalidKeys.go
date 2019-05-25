// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package bson

import (
	"fmt"
	"reflect"
)

type ErrInvalidKeys struct {
	Value reflect.Type
}

func (e ErrInvalidKeys) Error() string {
	return fmt.Sprintf("keys are of type %q, expecting string", e.Value)
}
