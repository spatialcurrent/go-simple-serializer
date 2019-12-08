// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package tagger includes a struct for marshaling and unmarshaling struct tags.
package tagger

import (
	"fmt"

	"github.com/pkg/errors"
)

// KeyValue represents a struct tag key:"value" pair.
type KeyValue struct {
	Key   string
	Value Value
}

// MarshalText returns the KeyValue formatted as a struct tag key:"value" pairs.
// Always returns a nil error.
func (t KeyValue) MarshalText() ([]byte, error) {
	v, err := t.Value.MarshalText()
	if err != nil {
		return make([]byte, 0), errors.Wrapf(err, "error marshaling %#v", t)
	}
	return []byte(fmt.Sprintf("%s:%q", t.Key, string(v))), nil
}
