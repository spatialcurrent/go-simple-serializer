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

	"github.com/pkg/errors"

	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// WriteInput provides the input for the Write function.
type WriteInput struct {
	Writer            io.Writer     // the underlying writer
	Keys              []interface{} // subset of keys to print
	ExpandKeys        bool          // dynamically expand keys
	KeyValueSeparator string        // the key-value separator
	LineSeparator     string        // the line separator
	Object            interface{}   // the object to write
	KeySerializer     stringify.Stringer
	ValueSerializer   stringify.Stringer
	Sorted            bool // sort keys
	Reversed          bool
	Limit             int
}

// Write writes the given object(s) as lines of tags.
// If the type of the input object is of kind Array or Slice, then writes each object on its own line.
// If the type of the input object is of kind Map or Struct, then writes a single line of tags.
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
		w := NewWriter(
			input.Writer,
			input.Keys,
			input.ExpandKeys,
			input.KeyValueSeparator,
			input.LineSeparator,
			input.KeySerializer,
			input.ValueSerializer,
			input.Sorted,
			input.Reversed,
		)
		errorRun := p.Input(it).Output(w).Run()
		if errorRun != nil {
			return errors.Wrap(errorRun, "error serializing arry or slice as tags")
		}
		return nil
	}

	// If not an array of slice, then just marshal.

	b, errMarshal := Marshal(
		inputObject,
		input.Keys,
		input.ExpandKeys,
		input.KeyValueSeparator,
		input.KeySerializer,
		input.ValueSerializer,
		input.Sorted,
		input.Reversed)
	if errMarshal != nil {
		return errors.Wrap(errMarshal, "error serializing to tags")
	}

	_, errWrite := input.Writer.Write(b)
	if errWrite != nil {
		return errors.Wrap(errWrite, "error writing to underlying writer")
	}

	return nil
}
