// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// Package writer provides an easy API to create a writer to write objects to a file.
// Depends on the following packages in go-simple-serializer.
//	- github.com/spatialcurrent/go-simple-serializer/pkg/jsonl
//	- github.com/spatialcurrent/go-simple-serializer/pkg/sv
//	- github.com/spatialcurrent/go-simple-serializer/pkg/tags
package writer

import (
	"io"

	"github.com/pkg/errors"

	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-simple-serializer/pkg/fmt"
	"github.com/spatialcurrent/go-simple-serializer/pkg/gob"
	"github.com/spatialcurrent/go-simple-serializer/pkg/jsonl"
	"github.com/spatialcurrent/go-simple-serializer/pkg/rapid"
	"github.com/spatialcurrent/go-simple-serializer/pkg/sv"
	"github.com/spatialcurrent/go-simple-serializer/pkg/tags"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

var (
	ErrMissingKeyValueSeparator = errors.New("missing key-value separator")
	ErrMissingLineSeparator     = errors.New("missing line separator")
)

// NewWriterInput includes the parameters for NewWriter function.
type NewWriterInput struct {
	Writer            io.Writer
	Format            string
	FormatSpecifier   string
	Header            []interface{}
	ExpandHeader      bool // in context, only used by tags as ExpandKeys
	KeySerializer     stringify.Stringer
	ValueSerializer   stringify.Stringer
	KeyValueSeparator string
	LineSeparator     string
	Fit               bool
	Pretty            bool
	Sorted            bool
	Reversed          bool
}

// NewWriter returns a new pipe.Writer for writing formatted objects to an underlying writer.
func NewWriter(input *NewWriterInput) (pipe.Writer, error) {

	switch input.Format {
	case "go", "jsonl", "tags":
		if len(input.LineSeparator) == 0 {
			return nil, ErrMissingLineSeparator
		}
	}

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
	case "fmt":
		w := fmt.NewWriter(input.Writer, input.FormatSpecifier, input.LineSeparator)
		return w, nil
	case "go":
		w := fmt.NewWriter(input.Writer, "%#v", input.LineSeparator)
		return w, nil
	case "gob":
		return gob.NewWriter(input.Writer, input.Fit), nil
	case "jsonl":
		w := jsonl.NewWriter(
			input.Writer,
			input.LineSeparator,
			input.KeySerializer,
			input.Pretty,
		)
		return w, nil
	case "rapid":
		return rapid.NewWriter(input.Writer), nil
	case "tags":
		if len(input.KeyValueSeparator) == 0 {
			return nil, ErrMissingKeyValueSeparator
		}
		w := tags.NewWriter(
			input.Writer,
			input.Header,
			input.ExpandHeader,
			input.KeyValueSeparator,
			input.LineSeparator,
			input.KeySerializer,
			input.ValueSerializer,
			input.Sorted,
			input.Reversed,
		)
		return w, nil
	}

	return nil, &ErrInvalidFormat{Format: input.Format}
}
