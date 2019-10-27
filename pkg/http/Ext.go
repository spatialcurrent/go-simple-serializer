// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package http

import (
	"net/http"
	"path/filepath"

	"github.com/pkg/errors"
)

var (
	ErrMissingURL = errors.New("missing URL")
)

// Ext returns the file name extension in the URL path.
// The extension begins after the last period in the file element of the path.
// If no period is in the last element or a period is the last character, then returns a blank string.
func Ext(r *http.Request) (string, error) {
	if r.URL == nil {
		return "", ErrMissingURL
	}
	ext := filepath.Ext(r.URL.Path)
	if len(ext) == 0 {
		return "", nil
	}
	if ext == "." {
		return "", nil
	}
	return ext[1:], nil
}
