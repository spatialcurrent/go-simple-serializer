// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tagger

import (
	"reflect"

	"github.com/pkg/errors"
)

// Lookup lookups and parses a struct tag key-value pair if found.
// If the lookup key is not found, then returns (nil, nil).
// Returns an error if there was an error unmarshaling the value of the struct tag key-value pair.
func Lookup(structTag reflect.StructTag, key string) (*Value, error) {
	if str, ok := structTag.Lookup(key); ok {
		v := &Value{}
		err := Unmarshal([]byte(str), v)
		if err != nil {
			return nil, errors.Wrapf(err, "error unmarshaling struct tag value %q", str)
		}
		return v, nil
	}
	return nil, nil
}
