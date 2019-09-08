// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package iterator provides an easy API to create an iterator to read objects from a file.
// Depends on the following packages in go-simple-serializer.
//	- github.com/spatialcurrent/go-simple-serializer/pkg/jsonl
//	- github.com/spatialcurrent/go-simple-serializer/pkg/sv
//	- github.com/spatialcurrent/go-simple-serializer/pkg/tags
package writer

import (
	"fmt"
	"io"

	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-simple-serializer/pkg/gob"
	"github.com/spatialcurrent/go-simple-serializer/pkg/jsonl"
	"github.com/spatialcurrent/go-simple-serializer/pkg/sv"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// Parameters for NewWriter function.
type NewWriterInput struct {
	Writer          io.Writer
	Format          string
	Header          []interface{}
	KeySerializer   stringify.Stringer
	ValueSerializer stringify.Stringer
	LineSeparator   string
	Fit             bool
	Pretty          bool
	Sorted          bool
	Reversed        bool
}

// NewWriter returns a new pipe.Writer for writing formatted objects to an underlying writer.
func NewWriter(input *NewWriterInput) (pipe.Writer, error) {

	switch input.Format {
	case "csv", "tsv":
		separator, err := sv.FormatToSeparator(input.Format)
		if err != nil {
			return nil, err
		}
		w := sv.NewWriter(
			input.Writer,
			separator,
			input.Header,
			input.KeySerializer,
			input.ValueSerializer,
			input.Sorted,
			input.Reversed,
		)
		return w, nil
	case "jsonl":
		return jsonl.NewWriter(input.Writer, input.LineSeparator, input.KeySerializer, input.Pretty), nil
	case "go":
		w := pipe.NewFunctionWriter(func(object interface{}) error {
			_, err := fmt.Fprintf(input.Writer, "%#v\n", object)
			return err
		})
		return w, nil
	case "gob":
		return gob.NewWriter(input.Writer, input.Fit), nil
	}

	return nil, &ErrInvalidFormat{Format: input.Format}
}
