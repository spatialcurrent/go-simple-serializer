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

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
)

func CreateHeaderAndKnownKeysFromValue(objectValue reflect.Value, sorted bool) ([]interface{}, map[interface{}]struct{}) {
	for objectValue.Type().Kind() == reflect.Ptr {
		objectValue = objectValue.Elem()
	}
	objectKind := objectValue.Type().Kind()
	if objectKind == reflect.Map {
		keys := inspector.GetKeysFromValue(objectValue, sorted)
		knownKeys := map[interface{}]struct{}{}
		for _, key := range keys {
			knownKeys[key] = struct{}{}
		}
		return keys, knownKeys
	} else if objectKind == reflect.Struct {
		fieldNames := inspector.GetFieldNamesFromValue(objectValue, sorted)
		header := make([]interface{}, 0, len(fieldNames))
		knownKeys := map[interface{}]struct{}{}
		for _, fieldName := range fieldNames {
			header = append(header, fieldName)
			knownKeys[fieldName] = struct{}{}
		}
		return header, knownKeys
	}
	return make([]interface{}, 0), map[interface{}]struct{}{}
}
