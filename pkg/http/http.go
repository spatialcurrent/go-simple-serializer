// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package http

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

const (
	HeaderAccept = "accept"
)

const (
	NoSkip    = 0  // used as SkipLines parameter to indicate no skipping when reading
	NoLimit   = -1 // used to indicate that there is no limit on reading or writing, depending on context.
	NoComment = "" // used to indicate that there is no comment prefix to consider.
)

var (
	// NoHeader is used to indicate that no defined header is given.
	// Derive the header from the input data.
	NoHeader = []interface{}{}
	// Formats is a list of all the formats supported by GSS
	Formats = serializer.Formats
)

var (
	ContentTypeBSON  = "application/ubjson"
	ContentTypeCSV   = "text/csv"
	ContentTypeJSON  = "application/json"
	ContentTypeTOML  = "application/toml"
	ContentTypeTSV   = "text/tab-separated-values"
	ContentTypeYAML  = "text/yaml"
	ContentTypePlain = "text/plain; charset=utf-8"
)
