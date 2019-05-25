// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"reflect"
)

import (
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
)

// Read reads the json lines from the input reader of the type given.
func Read(input *ReadInput) (interface{}, error) {
	it := NewIterator(&NewIteratorInput{
		Reader:        input.Reader,
		SkipLines:     input.SkipLines,
		SkipBlanks:    input.SkipBlanks,
		SkipComments:  input.SkipComments,
		Comment:       input.Comment,
		Trim:          input.Trim,
		Limit:         input.Limit,
		LineSeparator: input.LineSeparator,
		DropCR:        input.DropCR,
	})
	output := reflect.MakeSlice(input.Type, 0, 0).Interface()
	w := pipe.NewSliceWriterWithValues(output)
	err := pipe.NewBuilder().Input(it).Output(w).Run()
	return w.Values(), err
}
