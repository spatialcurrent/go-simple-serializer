// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/spatialcurrent/go-simple-serializer/pkg/registry"
	"github.com/spatialcurrent/go-simple-serializer/pkg/serializer"
)

// Respond writes the given data to the respond writer, and returns an error if any.
// If filename is not empty, then the "Content-Disposition" header is set to "attachment; filename=<FILENAME>".
func Respond(w http.ResponseWriter, r *http.Request, reg *registry.Registry, data interface{}, status int, filename string) error {

	contentType, format, err := NegotiateFormat(r, reg)
	if err != nil {
		ext, err := Ext(r)
		if err != nil || len(ext) == 0 {
			return errors.Errorf("could not negotiate format or parse file extension from %#v", r)
		}
		if item, ok := reg.LookupExtension(ext); ok {
			contentType = item.ContentTypes[0]
			format = item.Format
		} else {
			return errors.Errorf("could not negotiate format or parse file extension from %#v", r)
		}
	}

	s := serializer.New(format)

	body, err := s.Serialize(data)
	if err != nil {
		return errors.Wrap(err, "error serializing response body")
	}

	return RespondWithContent(w, body, contentType, status, filename)
}
