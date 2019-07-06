// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"io"
	"reflect"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// WriteInput provides the input for the Write function.
type WriteInput struct {
	Writer        io.Writer // the underlying writer
	LineSeparator string    // the line separator
	KeySerializer stringify.Stringer
	Object        interface{} // the object to write
	Pretty        bool        // pretty output
	Limit         int
}

// Write writes the given object(s) as JSON Lines (aka jsonl).
// If the type of the input object is of kind Array or Slice, then writes each object on its own line.
// Otherwise, the object is simply written as JSON.
func Write(input *WriteInput) error {
	inputObject := input.Object
	inputObjectValue := reflect.ValueOf(inputObject)
	for reflect.TypeOf(inputObjectValue.Interface()).Kind() == reflect.Ptr {
		inputObjectValue = inputObjectValue.Elem()
	}
	inputObjectValue = reflect.ValueOf(inputObjectValue.Interface()) // sets value to concerete type
	inputObjectKind := inputObjectValue.Type().Kind()

	p := pipe.NewBuilder().OutputLimit(input.Limit)
	if inputObjectKind == reflect.Array || inputObjectKind == reflect.Slice {
		if len(input.LineSeparator) == 0 {
			return ErrMissingLineSeparator
		}
		it, errorIterator := pipe.NewSliceIterator(inputObject)
		if errorIterator != nil {
			return errors.Wrap(errorIterator, "error creating slice iterator")
		}
		p = p.Input(it)
	} else {
		it, errorIterator := pipe.NewSliceIterator([]interface{}{inputObject})
		if errorIterator != nil {
			return errors.Wrap(errorIterator, "error creating slice iterator")
		}
		p = p.Input(it)
	}
	w := NewWriter(input.Writer, input.LineSeparator, input.KeySerializer, input.Pretty)
	errorRun := p.Output(w).Run()
	if errorRun != nil {
		return errors.Wrap(errorRun, "error serializing jsonl")
	}
	return nil
}
