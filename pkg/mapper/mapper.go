// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package mapper provides a simple api for serializing objects to maps.
package mapper

type Marshaler interface {
	MarshalMap() (interface{}, error)
}

type Unmarshaler interface {
	UnmarshalMap(data interface{}) error
}
