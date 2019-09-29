// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package properties

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/spatialcurrent/go-simple-serializer/pkg/escaper"
	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

// WriteInput provides the input for the Write function.
type WriteInput struct {
	Writer            io.Writer          // the underlying writer
	LineSeparator     string             // the newline byte
	KeyValueSeparator string             // the separator for key-value pairs
	Object            interface{}        // the object to write
	KeySerializer     stringify.Stringer // serializer for object properties
	ValueSerializer   stringify.Stringer // serializer for object properties
	Sorted            bool               // sort output
	Reversed          bool               // if sorted, sort in reverse alphabetical order
	EscapePrefix      string             // escape prefix, if empty then doesn't escape
	EscapeSpace       bool               // escape spaces
	EscapeEqual       bool               // escape =
	EscapeColon       bool               // escape :
	EscapeNewLine     bool               // escape \n
}

// Write writes the given rows as separated values.
func Write(input *WriteInput) error {

	if len(input.LineSeparator) == 0 {
		return ErrMissingLineSeparator
	}

	if len(input.KeyValueSeparator) == 0 {
		return ErrMissingKeyValueSeparator
	}

	inputObject := input.Object
	inputObjectValue := reflect.ValueOf(inputObject)
	for reflect.TypeOf(inputObjectValue.Interface()).Kind() == reflect.Ptr {
		inputObjectValue = inputObjectValue.Elem()
	}
	inputObjectValue = reflect.ValueOf(inputObjectValue.Interface()) // sets value to concerete type
	inputObjectKind := inputObjectValue.Type().Kind()

	// Initialize Escaper
	e := escaper.New()
	if len(input.EscapePrefix) > 0 {
		e = e.Prefix(input.EscapePrefix)
		if input.EscapeSpace {
			e = e.Sub(" ")
		}
		if input.EscapeEqual {
			e = e.Sub("=")
		}
		if input.EscapeColon {
			e = e.Sub(":")
		}
		if input.EscapeNewLine {
			e = e.Sub("\n")
		}
	}

	keySerializer := input.KeySerializer
	if keySerializer == nil {
		keySerializer = stringify.NewStringer("", false, false, false)
	}

	valueSerializer := input.ValueSerializer
	if valueSerializer == nil {
		valueSerializer = stringify.NewStringer("", false, false, false)
	}

	outputWriter := input.Writer
	kvSeparator := input.KeyValueSeparator
	lineSeparator := input.LineSeparator

	if inputObjectKind == reflect.Map {
		m := inputObjectValue
		keys := inspector.GetKeysFromValue(inputObjectValue, input.Sorted, input.Reversed)
		for i, key := range keys {
			keyString, errorKey := keySerializer(key)
			if errorKey != nil {
				return errors.Wrap(errorKey, "error serializing property key")
			}
			if strings.Contains(keyString, kvSeparator) {
				switch kvSeparator {
				case " ":
					if !input.EscapeSpace {
						return fmt.Errorf("if using key-value separator \" \" and a key contains \" \", then you must escape \" \"")
					}
				case "=":
					if !input.EscapeEqual {
						return fmt.Errorf("if using key-value separator \"=\" and a key contains \"=\", then you must escape \"=\"")
					}
				case ":":
					if !input.EscapeColon {
						return fmt.Errorf("if using key-value separator \":\" and a key contains \":\", then you must escape \":\"")
					}
				}
			}
			valueString, errorValue := valueSerializer(m.MapIndex(reflect.ValueOf(key)).Interface())
			if errorValue != nil {
				return errors.Wrap(errorValue, "error serializing property value")
			}
			line := e.Escape(keyString) + kvSeparator + e.Escape(valueString)
			if i < m.Len()-1 {
				line = line + lineSeparator
			}
			_, errorWrite := outputWriter.Write([]byte(line))
			if errorWrite != nil {
				return errors.Wrap(errorWrite, "error writing property to output writer")
			}
		}
		return nil
	}

	if inputObjectKind == reflect.Struct {
		s := inputObjectValue
		fieldNames := inspector.GetFieldNamesFromValue(inputObjectValue, input.Sorted, input.Reversed)
		for i, fieldName := range fieldNames {
			keyString, errorKey := keySerializer(fieldName)
			if errorKey != nil {
				return errors.Wrap(errorKey, "error serializing property key")
			}
			fieldValue, errorValue := valueSerializer(s.FieldByName(fieldName).Interface())
			if errorValue != nil {
				return errors.Wrap(errorValue, "error serializing property value")
			}
			// don't need to escape field name since go field names are already valid
			line := keyString + kvSeparator + e.Escape(fieldValue)
			if i < len(fieldNames)-1 {
				line = line + lineSeparator
			}
			_, errorWrite := outputWriter.Write([]byte(line))
			if errorWrite != nil {
				return errors.Wrap(errorWrite, "error writing property to output writer")
			}
		}
		return nil
	}

	return &ErrInvalidKind{Value: inputObjectValue.Type(), Expected: []reflect.Kind{reflect.Map, reflect.Struct}}
}
