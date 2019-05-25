// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package json

import (
	stdjson "encoding/json" // import the standard json library as stdjson
)

const (
	prefix = ""
	indent = "  "
)

// Marshal formats an object into a slice of bytes of JSON.
// If the pretty parameter is set, then prints pretty output with an indent.
func Marshal(obj interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return stdjson.MarshalIndent(obj, prefix, indent)
	}
	return stdjson.Marshal(obj)
}
