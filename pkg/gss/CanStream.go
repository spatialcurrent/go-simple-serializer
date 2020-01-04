// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

// CanStream returns true if you can process the data as a stream from the given input format to the output format.
// There are a few logical rules for deciding if streaming is possible.
// If you are sorting the output, then you cannot stream the data.
// If the output format has a header, then the input format must also have a header to stream.
func CanStream(inputFormat string, outputFormat string, outputSorted bool) bool {

	if outputSorted {
		return false
	}

	switch inputFormat {
	case serializer.FormatCSV, serializer.FormatTSV:
		switch outputFormat {
		case serializer.FormatCSV, serializer.FormatJSONL, serializer.FormatFmt, serializer.FormatGo, serializer.FormatGob, serializer.FormatRapid, serializer.FormatTags, serializer.FormatTSV:
			return true
		}
	case serializer.FormatJSONL, serializer.FormatGob, serializer.FormatRapid, serializer.FormatTags:
		switch outputFormat {
		case serializer.FormatJSONL, serializer.FormatFmt, serializer.FormatGo, serializer.FormatGob, serializer.FormatRapid, serializer.FormatTags:
			return true
		}
	}

	return false

}
