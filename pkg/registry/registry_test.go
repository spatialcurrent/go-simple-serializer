// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistryLookupFormatNil(t *testing.T) {
	r := New()
	item, ok := r.LookupFormat("json")
	assert.False(t, ok)
	assert.Equal(t, Item{}, item)
}

func TestRegistryLookupFormatJSON(t *testing.T) {
	r := New()
	expected := Item{
		Format:       "json",
		ContentTypes: []string{"application/json", "text/json"},
		Extensions:   []string{"json"},
	}
	r.Add(expected)
	item, ok := r.LookupFormat("json")
	assert.True(t, ok)
	assert.Equal(t, expected, item)
}

func TestRegistryLookupFormatCSV(t *testing.T) {
	r := New()
	r.Add(Item{
		Format:       "json",
		ContentTypes: []string{"application/json", "text/json"},
		Extensions:   []string{"json"},
	})
	expected := Item{
		Format:       "csv",
		ContentTypes: []string{"text/csv"},
		Extensions:   []string{"csv"},
	}
	r.Add(expected)
	item, ok := r.LookupFormat("csv")
	assert.True(t, ok)
	assert.Equal(t, expected, item)
}

func TestRegistryLookupContentTypeCSV(t *testing.T) {
	r := New()
	r.Add(Item{
		Format:       "json",
		ContentTypes: []string{"application/json", "text/json"},
		Extensions:   []string{"json"},
	})
	expected := Item{
		Format:       "csv",
		ContentTypes: []string{"text/csv"},
		Extensions:   []string{"csv"},
	}
	r.Add(expected)
	item, ok := r.LookupContentType("text/csv")
	assert.True(t, ok)
	assert.Equal(t, expected, item)
}
