// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"io"
	"reflect"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
)

// WriteInput provides the input for the Write function.
type WriteInput struct {
	Writer          io.Writer   // the underlying writer
	LineSeparator   string      // the line separator
	Object          interface{} // the object to write
	ValueSerializer func(object interface{}) (string, error)
	Sorted          bool // sort keys
}

// Write writes the given object(s) as JSON Lines (aka jsonl).
// If the type of the input object is of kind Array or Slice, then writes each object on its own line.
// Otherwise, the object is simply written as JSON.
func Write(input *WriteInput) error {
	inputObject := input.Object
	inputObjectValue := reflect.ValueOf(inputObject)
	k := inputObjectValue.Type().Kind()
	if k == reflect.Ptr {
		inputObjectValue = inputObjectValue.Elem()
		k = inputObjectValue.Type().Kind()
	}
	p := pipe.NewBuilder()
	if k == reflect.Array || k == reflect.Slice {
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
	errorRun := p.Output(NewWriter(input.Writer, input.LineSeparator, input.ValueSerializer, input.Sorted)).Run()
	if errorRun != nil {
		return errors.Wrap(errorRun, "error serializing jsonl")
	}
	return nil
}
