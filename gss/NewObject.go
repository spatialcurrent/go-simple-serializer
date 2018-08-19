package gss

import (
	"strings"
	"unicode"
)

func NewObject(content string, format string) interface{} {
	if format == "json" {
		content = strings.TrimLeftFunc(content, unicode.IsSpace)
		if len(content) > 0 && content[0] == '[' {
			return []map[string]interface{}{}
		}
		return map[string]interface{}{}
	} else if format == "yaml" {
		content = strings.TrimLeftFunc(content, unicode.IsSpace)
		if len(content) > 0 && content[0] == '-' {
			return []map[string]interface{}{}
		}
		return map[string]interface{}{}
	} else if format == "bson" || format == "hcl" || format == "hcl2" || format == "properties" || format == "toml" {
		return map[string]interface{}{}
	}

	return nil
}
