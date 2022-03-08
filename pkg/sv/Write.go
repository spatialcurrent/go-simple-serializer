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
	"sort"

	"github.com/spatialcurrent/go-pipe/pkg/pipe"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// WriteInput provides the input for the Write function.
type WriteInput struct {
	Writer          io.Writer // the underlying writer
	Separator       rune      // the values separator
	Header          []interface{}
	KeySerializer   stringify.Stringer
	ValueSerializer stringify.Stringer
	Object          interface{} // the object to write
	ExpandHeader    bool        // dynamically expand header, requires caching output in memory
	Sorted          bool        // sort columns
	Reversed        bool        // if sorted, sort in reverse alphabetical order.
	Limit           int
}

// Write writes the given object(s) as separated values, e.g., CSV or T
// If the type of the input object is of kind Array or Slice, then writes each object on its own line.
// Otherwise, just writes a CSV with a header and one row for the object.
func Write(input *WriteInput) error {

	// If the input is a slice of strings, simply write it as a row
	if slc, ok := input.Object.([]string); ok {
		errWriteTable := WriteTable(&WriteTableInput{
			Writer:    input.Writer,
			Separator: input.Separator,
			Header:    make([]string, 0),
			Rows:      [][]string{slc},
			Sorted:    input.Sorted, // if sorted and no specific wilcard position
			Reversed:  input.Reversed,
		})
		if errWriteTable != nil {
			return fmt.Errorf("error writing table to underlying writer: %w", errWriteTable)
		}
	}

	inputObject := input.Object
	inputObjectValue := reflect.ValueOf(inputObject)
	for reflect.TypeOf(inputObjectValue.Interface()).Kind() == reflect.Ptr {
		inputObjectValue = inputObjectValue.Elem()
	}
	inputObjectValue = reflect.ValueOf(inputObjectValue.Interface()) // sets value to concerete type
	inputObjectKind := inputObjectValue.Type().Kind()

	if input.ExpandHeader {

		// set the key serializer
		keySerializer := input.KeySerializer
		if keySerializer == nil {
			keySerializer = stringify.NewStringer("", false, false, false)
		}

		// set the value serializer
		valueSerializer := input.ValueSerializer
		if valueSerializer == nil {
			valueSerializer = stringify.NewStringer("", false, false, false)
		}

		// initialize header, wildcard, and known keys
		header := input.Header
		wildcard := false
		knownKeys := map[interface{}]struct{}{}
		if len(header) > 0 {
			for _, k := range header {
				if k == "*" {
					wildcard = true
				} else {
					knownKeys[k] = struct{}{}
				}
			}
			if input.Sorted {
				sort.Slice(header, func(i, j int) bool {
					return (fmt.Sprint(header[i]) < fmt.Sprint(header[j])) && (!input.Reversed)
				})
			}
		}

		rows := make([][]string, 0)

		switch inputObjectKind {
		case reflect.Map, reflect.Struct:
			if len(header) == 0 {
				header, _ = CreateHeaderAndKnownKeysFromValue(inputObjectValue, input.Sorted, input.Reversed)
			} else if wildcard {
				header, _ = ExpandHeaderWithWildcard(header, knownKeys, inputObjectValue, input.Sorted, input.Reversed)
			} else {
				header, _ = ExpandHeader(header, knownKeys, inputObjectValue, input.Sorted, input.Reversed)
			}
			row, err := ToRowFromValue(inputObjectValue, header, valueSerializer)
			if err != nil {
				return fmt.Errorf("error serializing object to row: %w", err)
			}
			rows = append(rows, row)
		case reflect.Array, reflect.Slice:
			if inputObjectValue.Len() == 0 {
				// If there are no records then just return an empty string
				return nil
			}
			if len(header) == 0 {
				header, knownKeys = CreateHeaderAndKnownKeysFromValue(inputObjectValue.Index(0), input.Sorted, input.Reversed)
				for i := 0; i < inputObjectValue.Len() && (input.Limit < 0 || i <= input.Limit); i++ {
					header, knownKeys = ExpandHeader(header, knownKeys, concerete(inputObjectValue.Index(i)), input.Sorted, input.Reversed)
					row, err := ToRowFromValue(concerete(inputObjectValue.Index(i)), header, valueSerializer)
					if err != nil {
						return fmt.Errorf("error serializing object to row: %w", err)
					}
					rows = append(rows, row)
				}
			} else if wildcard {
				// With a wildcard, must expand header before creating rows
				for i := 0; i < inputObjectValue.Len() && (input.Limit < 0 || i <= input.Limit); i++ {
					header, knownKeys = ExpandHeaderWithWildcard(header, knownKeys, concerete(inputObjectValue.Index(i)), input.Sorted, input.Reversed)
				}
				header = RemoveWildcard(header)
				for i := 0; i < inputObjectValue.Len() && (input.Limit < 0 || i <= input.Limit); i++ {
					row, err := ToRowFromValue(concerete(inputObjectValue.Index(i)), header, valueSerializer)
					if err != nil {
						return fmt.Errorf("error serializing object to row: %w", err)
					}
					rows = append(rows, row)
				}
			} else {
				// If the length of the initial header is greater than zero and does not include a wildcard.
				// Can expand header inline
				for i := 0; i < inputObjectValue.Len() && (input.Limit < 0 || i <= input.Limit); i++ {
					header, knownKeys = ExpandHeader(header, knownKeys, concerete(inputObjectValue.Index(i)), input.Sorted, input.Reversed)
					row, err := ToRowFromValue(concerete(inputObjectValue.Index(i)), header, valueSerializer)
					if err != nil {
						return fmt.Errorf("error serializing object to row: %w", err)
					}
					rows = append(rows, row)
				}
			}
		}
		outputHeader, errStringifyHeader := stringify.StringifySlice(header, keySerializer)
		if errStringifyHeader != nil {
			return fmt.Errorf("error stringifying header %q: %w", header, errStringifyHeader)
		}
		errWriteTable := WriteTable(&WriteTableInput{
			Writer:    input.Writer,
			Separator: input.Separator,
			Header:    outputHeader,
			Rows:      rows,
			Sorted:    input.Sorted && !wildcard, // if sorted and no specific wilcard position
			Reversed:  input.Reversed,
		})
		if errWriteTable != nil {
			return fmt.Errorf("error writing table to underlying writer: %w", errWriteTable)
		}
	}

	// if streaming and not expanding header.
	p := pipe.NewBuilder().OutputLimit(input.Limit)
	if inputObjectKind == reflect.Array || inputObjectKind == reflect.Slice {
		fmt.Println("Input Object:", input.Header, inputObject)
		it, errorIterator := pipe.NewSliceIterator(inputObject)
		if errorIterator != nil {
			return fmt.Errorf("error creating slice iterator: %w", errorIterator)
		}
		p = p.Input(it)
		p = p.Output(NewWriter(
			input.Writer,
			input.Separator,
			input.Header,
			input.KeySerializer,
			input.ValueSerializer,
			input.Sorted,
			input.Reversed,
		))
		errorRun := p.Run()
		if errorRun != nil {
			return fmt.Errorf("error serializing separated values: %w", errorRun)
		}
		return nil
	}

	w := NewWriter(
		input.Writer,
		input.Separator,
		input.Header,
		input.KeySerializer,
		input.ValueSerializer,
		input.Sorted,
		input.Reversed,
	)

	errorWrite := w.WriteObject(inputObject)
	if errorWrite != nil {
		return fmt.Errorf("error serializing separated values: %w", errorWrite)
	}

	errorFlush := w.Flush()
	if errorFlush != nil {
		return fmt.Errorf("error serializing separated values: %w", errorFlush)
	}

	return nil
}
