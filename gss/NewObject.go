// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"strings"
	"unicode"
)

func NewObject(content string, format string) (interface{}, error) {

	if format == "json" {
		content = strings.TrimLeftFunc(content, unicode.IsSpace)
		if len(content) > 0 && content[0] == '[' {
			return []map[string]interface{}{}, nil
		}
		return map[string]interface{}{}, nil
	} else if format == "yaml" {
		content = strings.TrimLeftFunc(content, unicode.IsSpace)
		if len(content) > 0 && content[0] == '-' {
			return []map[string]interface{}{}, nil
		}
		return map[string]interface{}{}, nil
	} else if format == "bson" || format == "hcl" || format == "hcl2" || format == "properties" || format == "toml" {
		return map[string]interface{}{}, nil
	} else if format == "jsonl" || format == "csv" || format == "tsv" {
		return []map[string]interface{}{}, nil
	}

	return nil, nil
}
