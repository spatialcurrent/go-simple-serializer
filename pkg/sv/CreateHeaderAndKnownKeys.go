// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"reflect"
)

func CreateHeaderAndKnownKeys(object interface{}, sorted bool) ([]interface{}, map[interface{}]struct{}) {
	return CreateHeaderAndKnownKeysFromValue(reflect.ValueOf(object), sorted)
}
