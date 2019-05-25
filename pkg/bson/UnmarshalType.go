// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package bson

import (
	"fmt"
	// import the mgo bson library
	mgobson "gopkg.in/mgo.v2/bson"
	"reflect"
)

import (
	"github.com/pkg/errors"
)

// UnmarshalType parses a slice of bytes into an object of a given type.
// If no input is given, then returns ErrEmptyInput.
func UnmarshalType(b []byte, outputType reflect.Type) (interface{}, error) {

	if len(b) == 0 {
		return nil, ErrEmptyInput
	}

	k := outputType.Kind()

	if k == reflect.Map {
		if outputType.Key().Kind() != reflect.String {
			return nil, &ErrInvalidKeys{Value: outputType.Elem()}
		}
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeMap(outputType))
		err := mgobson.Unmarshal(b, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling BSON %q", string(b)))
		}
		return ptr.Elem().Interface(), nil
	}

	if k == reflect.Struct {
		ptr := reflect.New(outputType)
		err := mgobson.Unmarshal(b, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling BSON %q", string(b)))
		}
		return ptr.Elem().Interface(), nil
	}

	if outputType == reflect.TypeOf(mgobson.D{}) {
		ptr := reflect.New(outputType)
		ptr.Elem().Set(reflect.MakeMap(outputType))
		err := mgobson.Unmarshal(b, ptr.Interface())
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error unmarshaling BSON %q", string(b)))
		}
		return ptr.Elem().Interface(), nil
	}

	return nil, &ErrInvalidKind{Value: outputType, Expected: []reflect.Kind{reflect.Map, reflect.Struct}}
}
