// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"fmt"
	"io"
	"reflect"

	"github.com/spatialcurrent/go-object/pkg/object"
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
)

// ReadInput provides the input for the Read function.
type ReadInput struct {
	Type       reflect.Type // the output type
	Reader     io.Reader
	Separator  rune // the values separator
	Header     object.ObjectArray
	SkipLines  int
	Comment    string
	LazyQuotes bool
	Limit      int
}

// Read reads the separated values from the input reader into a slice.
func Read(input *ReadInput) (interface{}, error) {

	// If input.Type is nil, then use []map[string]string{}.
	inputType := reflect.TypeOf([]map[string]string{})
	if input.Type != nil {
		inputType = input.Type
	}

	// The iterator requires the type to return for each element,
	// rather than the type of the array itself.
	iteratorType := inputType.Elem()
	if iteratorType.Kind() == reflect.Interface {
		iteratorType = reflect.TypeOf(map[string]string{})
	}

	it, errorIterator := NewIterator(&NewIteratorInput{
		Reader:     input.Reader,
		Type:       iteratorType,
		Separator:  input.Separator,
		Header:     input.Header,
		Comment:    input.Comment,
		SkipLines:  input.SkipLines,
		LazyQuotes: input.LazyQuotes,
		Limit:      input.Limit,
	})
	if errorIterator != nil {
		return nil, fmt.Errorf("error creating iterator: %w", errorIterator)
	}
	output := reflect.MakeSlice(inputType, 0, 0).Interface()
	w := pipe.NewSliceWriterWithValues(output)
	errorRun := pipe.NewBuilder().Input(it).Output(w).Run()
	return w.Values(), errorRun
}
