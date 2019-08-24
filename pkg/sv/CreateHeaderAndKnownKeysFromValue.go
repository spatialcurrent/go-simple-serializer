// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"reflect"

	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
)

// CreateHeaderAndKnownKeysFromValue returns an object's keys or field names as a slice and set.
func CreateHeaderAndKnownKeysFromValue(objectValue reflect.Value, sorted bool, reversed bool) ([]interface{}, map[interface{}]struct{}) {
	for reflect.TypeOf(objectValue.Interface()).Kind() == reflect.Ptr {
		objectValue = objectValue.Elem()
	}
	objectValue = reflect.ValueOf(objectValue.Interface()) // sets value to concerete type
	switch objectValue.Type().Kind() {
	case reflect.Map:
		keys := inspector.GetKeys(objectValue.Interface(), sorted, reversed)
		knownKeys := map[interface{}]struct{}{}
		for _, key := range keys {
			knownKeys[key] = struct{}{}
		}
		return keys, knownKeys
	case reflect.Struct:
		fieldNames := inspector.GetFieldNames(objectValue.Interface(), sorted, reversed)
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
