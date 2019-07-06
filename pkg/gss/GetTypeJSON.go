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

func GetTypeJSON(str string) reflect.Type {
	if len(str) > 0 && str[0] == '[' {
		return interfaceSliceType
	}
	return mapStringInterfaceType
}
