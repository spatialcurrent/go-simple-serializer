// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package registry provides a library for managing a list of known file types.
package registry

type Item struct {
	Format       string
	ContentTypes []string
	Extensions   []string
}

type Registry struct {
	items        []Item
	formats      map[string]Item
	contentTypes map[string]Item
	extensions   map[string]Item
}

func (r Registry) LookupFormat(format string) (Item, bool) {
	item, ok := r.formats[format]
	return item, ok
}

func (r Registry) LookupContentType(contentType string) (Item, bool) {
	item, ok := r.contentTypes[contentType]
	return item, ok
}

func (r Registry) LookupExtension(extension string) (Item, bool) {
	item, ok := r.extensions[extension]
	return item, ok
}

func (r *Registry) Add(item Item) {
	r.items = append(r.items, item)
	if len(item.Format) > 0 {
		r.formats[item.Format] = item
	}
	if len(item.ContentTypes) > 0 {
		for _, key := range item.ContentTypes {
			r.contentTypes[key] = item
		}
	}
	if len(item.Extensions) > 0 {
		for _, key := range item.Extensions {
			r.extensions[key] = item
		}
	}
}

func New() *Registry {
	return &Registry{
		items:        make([]Item, 0),
		formats:      map[string]Item{},
		contentTypes: map[string]Item{},
		extensions:   map[string]Item{},
	}
}
