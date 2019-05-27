// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

// SerializeInput provides the input for the SerializeString and SerializeBytes function.
type SerializeInput struct {
	Object            interface{}
	Format            string
	Header            []interface{}
	Limit             int
	Pretty            bool
	Sorted            bool
	LineSeparator     string
	KeyValueSeparator string
	ValueSerializer   func(object interface{}) (string, error)
	EscapePrefix      string
	EscapeSpace       bool
	EscapeNewLine     bool
	EscapeEqual       bool
	EscapeColon       bool
}
