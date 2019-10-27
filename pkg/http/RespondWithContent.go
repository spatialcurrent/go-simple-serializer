// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package http

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// RespondWithContent writes the given content to the response writer, and returns an error if any.
// If filename is not empty, then the "Content-Disposition" header is set to "attachment; filename=<FILENAME>".
func RespondWithContent(w http.ResponseWriter, body []byte, contentType string, status int, filename string) error {

	if len(filename) > 0 {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	}

	w.Header().Set("Content-Type", contentType)

	if status != http.StatusOK {
		w.WriteHeader(status)
	}

	_, err := w.Write(body)
	if err != nil {
		return errors.Wrap(err, "error writing response body")
	}

	return nil
}
