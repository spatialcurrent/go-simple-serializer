// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package tagger includes a struct for marshaling and unmarshaling struct tags.
package tagger

import (
	"strings"
)

type Value struct {
	Ignore    bool
	Name      string
	OmitEmpty bool
}

// UnmarshalText unmarshals the given struct tag value into a Value object.
func (t *Value) UnmarshalText(text []byte) error {
	str := string(text)

	if len(str) == 0 {
		t.Ignore = false
		t.Name = ""
		t.OmitEmpty = false
		return nil
	}

	if str == "-" {
		t.Ignore = true
		t.Name = ""
		t.OmitEmpty = false
		return nil
	}

	if str == ",omitempty" {
		t.Ignore = false
		t.Name = ""
		t.OmitEmpty = true
		return nil
	}

	if strings.Contains(str, ",") {
		t.Ignore = false
		parts := strings.Split(str, ",")
		t.Name = parts[0]
		attrs := map[string]struct{}{}
		for _, p := range parts[1:] {
			attrs[p] = struct{}{}
		}
		_, omitEmpty := attrs["omitempty"]
		t.OmitEmpty = omitEmpty
		return nil
	}

	t.Ignore = false
	t.Name = str
	t.OmitEmpty = false
	return nil
}

// MarshalText returns the Value formatted as a struct tag value.
// Always returns a nil error.
func (t Value) MarshalText() ([]byte, error) {
	if t.Ignore {
		return []byte("-"), nil
	}
	str := t.Name
	if t.OmitEmpty {
		str += ",omitempty"
	}
	return []byte(str), nil
}
