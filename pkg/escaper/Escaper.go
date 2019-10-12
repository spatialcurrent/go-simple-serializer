// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package escaper provides an easy API to escapae and unescape a string using a custom escape prefix.
// See the examples below for usage.
package escaper

import (
	"strings"
)

// Escaper provides an easy API to escape and unescape a string using a custom escape prefix.
type Escaper struct {
	prefix     string   // the escape prefix string.
	substrings []string // slice of strings that will be escaped/unescaped.
}

// New returns a new escaper with the default prefix "\" and no additional substrings.
// This default escaper will only escape the prefix itself, so "\" becomes "\\".
func New() *Escaper {
	return &Escaper{
		prefix:     "",
		substrings: make([]string, 0), // as default, don't escape any strings other than the prefix itself
	}
}

// Clone returns a new copy of this escaper.
func (e *Escaper) Clone() *Escaper {
	substrings := make([]string, 0, len(e.substrings))
	copy(substrings, e.substrings)
	return &Escaper{
		prefix:     e.prefix,
		substrings: substrings,
	}
}

// Prefix sets the escape prefix to use for escaping substrings.
// The escaper will also escape the prefix itself.
// Any string can be used.
func (e *Escaper) Prefix(prefix string) *Escaper {
	e.prefix = prefix
	return e
}

// Sub adds the given substrings to the Escaper to be escaped and unescaped.
func (e *Escaper) Sub(substrings ...string) *Escaper {
	e.substrings = append(e.substrings, substrings...)
	return e
}

// Escape returns the escaped version of the given string.
func (e *Escaper) Escape(in string) string {
	out := in
	if len(e.prefix) > 0 {
		out = strings.Replace(out, e.prefix, e.prefix+e.prefix, -1) // escape the prefix itself
		for _, str := range e.substrings {
			out = strings.Replace(out, str, e.prefix+str, -1)
		}
	}
	return out
}

// Unescape returns the unescaped version of the given string.
func (e *Escaper) Unescape(in string) string {
	out := in
	if len(e.prefix) > 0 {
		for _, str := range e.substrings {
			out = strings.Replace(out, e.prefix+str, str, -1)
		}
		out = strings.Replace(out, e.prefix+e.prefix, e.prefix, -1) // unescape the prefix itself
	}
	return out
}
