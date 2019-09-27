// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gob

import (
	"io"
	"reflect"

	"github.com/pkg/errors"
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
)

// WriteInput provides the input for the Write function.
type WriteInput struct {
	Writer io.Writer   // the underlying writer
	Object interface{} // the object to write
	Limit  int
	Fit    bool
}

// Write writes the given object(s) as gob-encoded items.
// As gob is a stream encoding, if provded an object of array or slice, then each contained element is written as its own item.  If not provided an array or slice, then the provided object is encoded as is as a gob item.
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
	w := NewWriter(input.Writer, input.Fit)
	errorRun := p.Output(w).Run()
	if errorRun != nil {
		return errors.Wrap(errorRun, "error serializing gob")
	}
	return nil
}
