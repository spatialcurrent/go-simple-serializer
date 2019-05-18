// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

// SerializeJSON serializes an object to its JSON representation.
func SerializeJSON(obj interface{}, pretty bool) (string, error) {
	return SerializeString(&SerializeInput{
		Object: obj,
		Format: "json",
		Header: []string{},
		Limit:  NoLimit,
		Pretty: pretty,
	})
}
