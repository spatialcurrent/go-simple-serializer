// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

// SerializeInput provides the input for the SerializeString and SerializeBytes function.
type SerializeInput struct {
	Object interface{}
	Format string
	Header []string
	Limit  int
	Pretty bool
}
