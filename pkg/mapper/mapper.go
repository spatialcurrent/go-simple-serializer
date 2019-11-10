// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package mapper provides a simple api for serializing structs to maps and pointers to their concerete types.
// This package is useful for providing a single execution path for serializing to file formats.
// For example, using map struct tags, you can provide a single representation to JSON, YAML, and TOML serializers.
package mapper

type Marshaler interface {
	MarshalMap() (interface{}, error)
}

type Unmarshaler interface {
	UnmarshalMap(data interface{}) error
}
