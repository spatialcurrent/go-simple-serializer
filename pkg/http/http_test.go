// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package http

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/registry"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

func NewDefaultRegistry() *registry.Registry {
	r := registry.New()
	r.Add(registry.Item{
		Format:       serializer.FormatBSON,
		ContentTypes: []string{"application/ubjson"},
		Extensions:   []string{},
	})
	r.Add(registry.Item{
		Format:       serializer.FormatJSON,
		ContentTypes: []string{"application/json", "text/json"},
		Extensions:   []string{"json"},
	})
	r.Add(registry.Item{
		Format:       serializer.FormatCSV,
		ContentTypes: []string{"text/csv"},
		Extensions:   []string{"csv"},
	})
	r.Add(registry.Item{
		Format:       serializer.FormatYAML,
		ContentTypes: []string{"text/yaml"},
		Extensions:   []string{"yaml", "yml"},
	})
	return r
}
