// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"reflect"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-pipe/pkg/pipe"
)

// Read reads the json lines from the input reader into a slice.
func Read(input *ReadInput) (interface{}, error) {
	iteratorType := input.Type.Elem()
	if iteratorType.Kind() == reflect.Interface {
		iteratorType = reflect.TypeOf(map[string]string{})
	}
	it, errorIterator := NewIterator(&NewIteratorInput{
		Reader:     input.Reader,
		Type:       iteratorType,
		Separator:  input.Separator,
		Comment:    input.Comment,
		SkipLines:  input.SkipLines,
		LazyQuotes: input.LazyQuotes,
		Limit:      input.Limit,
	})
	if errorIterator != nil {
		return nil, errors.Wrap(errorIterator, "error creating iterator")
	}
	output := reflect.MakeSlice(input.Type, 0, 0).Interface()
	w := pipe.NewSliceWriterWithValues(output)
	errorRun := pipe.NewBuilder().Input(it).Output(w).Run()
	return w.Values(), errorRun
}
