// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package rapid

import (
	"io"
	"reflect"

	"github.com/spatialcurrent/go-pipe/pkg/pipe"
)

// ReadInput provides the input for the Read function.
type ReadInput struct {
	Type   reflect.Type // the output type
	Reader io.Reader    // the underlying reader
	DropCR bool         // drop carriage return
	Limit  int
}

// Read reads the json lines from the input reader of the type given.
func Read(input *ReadInput) (interface{}, error) {

	var inputType reflect.Type
	if input.Type != nil {
		inputType = input.Type.Elem()
	}

	outputType := reflect.TypeOf([]interface{}{})
	if input.Type != nil {
		outputType = input.Type
	}

	it := NewIterator(&NewIteratorInput{
		Type:   inputType,
		Reader: input.Reader,
		Limit:  input.Limit,
		DropCR: input.DropCR,
	})

	output := reflect.MakeSlice(outputType, 0, 0).Interface()

	w := pipe.NewSliceWriterWithValues(output)

	err := pipe.NewBuilder().Input(it).Output(w).Run()
	return w.Values(), err
}
