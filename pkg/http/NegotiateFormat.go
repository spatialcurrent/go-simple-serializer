// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package http

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/spatialcurrent/go-simple-serializer/pkg/registry"
)

var (
	ErrMissingAcceptHeader = errors.New("missing accept header")
	ErrMissingRegistry     = errors.New("missing file type registry")
)

// NegotiateFormat negotitates the format for the response based on the incoming request and the given file type registry.
// Returns the matching content type, followed by the format known to GSS, and then an error if any.
func NegotiateFormat(r *http.Request, reg *registry.Registry) (string, string, error) {

	accept := strings.TrimSpace(r.Header.Get(HeaderAccept))

	if len(accept) == 0 {
		return "", "", ErrMissingAcceptHeader
	}

	if reg == nil {
		return "", "", ErrMissingRegistry
	}

	// Parse accept header into map of weights to accepted values
	values := map[float64][]string{}
	for _, str := range strings.SplitN(accept, ",", -1) {
		v := strings.TrimSpace(str)
		if strings.Contains(v, ";q=") {
			parts := strings.SplitN(v, ";q=", 2)
			w, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return "", "", errors.Wrapf(err, "could not parse quality value for value %q", v)
			}
			if _, ok := values[w]; !ok {
				values[w] = make([]string, 0)
			}
			values[w] = append(values[w], strings.TrimSpace(parts[0]))

		} else {
			if _, ok := values[1.0]; !ok {
				values[1.0] = make([]string, 0)
			}
			values[1.0] = append(values[1.0], v)
		}
	}

	// Create list of weights
	weights := make([]float64, 0, len(values))
	for w := range values {
		weights = append(weights, w)
	}

	// Sort by weigt in descending order
	sort.SliceStable(weights, func(i, j int) bool {
		return weights[i] > weights[j]
	})

	// Iterate through accepted values in order of highest weight first
	for _, w := range weights {
		for _, contentType := range values[w] {
			if item, ok := reg.LookupContentType(contentType); ok {
				return contentType, item.Format, nil
			}
		}
	}
	return "", "", &ErrNotNegotiable{Value: accept}
}
