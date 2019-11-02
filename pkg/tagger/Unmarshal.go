// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package tagger includes a struct for marshaling and unmarshaling struct tags.
package tagger

// Unmarshal unmarshal the struct tag value into the Value object.
func Unmarshal(text []byte, value *Value) error {
	return value.UnmarshalText(text)
}
