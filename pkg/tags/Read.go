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
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
)

// ReadInput provides the input for the Read function.
type ReadInput struct {
	Type          reflect.Type // the output type
	Reader        io.Reader    // the underlying reader
	SkipLines     int
	SkipBlanks    bool
	SkipComments  bool
	Comment       string // the comment prefix
	LineSeparator byte   // the newline byte
	DropCR        bool   // drop carriage return
	Limit         int
}

// Read reads the lines of tags from the input Reader into the given type.
// If no type is given, returns a slice of type []map[string]string.
func Read(input *ReadInput) (interface{}, error) {
	inputType := reflect.TypeOf([]map[string]string{})
	if input.Type != nil {
		inputType = input.Type
	}
	it := NewIterator(&NewIteratorInput{
		Reader:        input.Reader,
		SkipLines:     input.SkipLines,
		SkipBlanks:    input.SkipBlanks,
		SkipComments:  input.SkipComments,
		Comment:       input.Comment,
		Limit:         input.Limit,
		LineSeparator: input.LineSeparator,
		DropCR:        input.DropCR,
	})
	output := reflect.MakeSlice(inputType, 0, 0).Interface()
	w := pipe.NewSliceWriterWithValues(output)
	err := pipe.NewBuilder().Input(it).Output(w).Run()
	return w.Values(), err
}
