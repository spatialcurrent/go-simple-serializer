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

	"github.com/spatialcurrent/go-pipe/pkg/pipe"
)

// ReadInput provides the input for the Read function.
type ReadInput struct {
	Type   reflect.Type // the output type
	Reader io.Reader    // the underlying reader
	Limit  int
}

// Read reads the json lines from the input reader of the type given.
func Read(input *ReadInput) (interface{}, error) {

	if input.Type == nil {
		return nil, ErrMissingType
	}

	it := NewIterator(&NewIteratorInput{
		Type:   input.Type.Elem(),
		Reader: input.Reader,
		Limit:  input.Limit,
	})

	output := reflect.MakeSlice(input.Type, 0, 0).Interface()

	w := pipe.NewSliceWriterWithValues(output)

	err := pipe.NewBuilder().Input(it).Output(w).Run()

	return w.Values(), err
}
