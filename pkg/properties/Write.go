// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package properties

import (
	"fmt"
	"reflect"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-simple-serializer/pkg/escaper"
	"github.com/spatialcurrent/go-simple-serializer/pkg/inspector"
)

// Write writes the given rows as separated values.
func Write(input *WriteInput) error {

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

	inputObject := input.Object
	valueSerializer := input.ValueSerializer
	inputObjectValue := reflect.ValueOf(inputObject)
	k := inputObjectValue.Type().Kind()
	if k == reflect.Ptr {
		inputObjectValue = inputObjectValue.Elem()
		k = inputObjectValue.Type().Kind()
	}
	outputWriter := input.Writer
	kvSeparator := input.KeyValueSeparator
	lineSeparator := input.LineSeparator

	if k == reflect.Map {
		m := inputObjectValue
		keys := inspector.GetKeysFromValue(inputObjectValue, input.Sorted)
		for i, key := range keys {
			keyString, errorKey := valueSerializer(key)
			if errorKey != nil {
				return errors.Wrap(errorKey, "error serializing property key")
			}
			if strings.Contains(keyString, kvSeparator) {
				switch kvSeparator {
				case " ":
					if !input.EscapeSpace {
						return fmt.Errorf("if using key-value separator \" \" and a key contains \" \", then you must escape \" \".")
					}
				case "=":
					if !input.EscapeEqual {
						return fmt.Errorf("if using key-value separator \"=\" and a key contains \"=\", then you must escape \"=\".")
					}
				case ":":
					if !input.EscapeColon {
						return fmt.Errorf("if using key-value separator \":\" and a key contains \":\", then you must escape \":\".")
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
	} else if k == reflect.Struct {
		s := inputObjectValue
		fieldNames := inspector.GetFieldNamesFromValue(inputObjectValue, input.Sorted)
		for i, fieldName := range fieldNames {
			fieldValue, errorValue := valueSerializer(s.FieldByName(fieldName).Interface())
			if errorValue != nil {
				return errors.Wrap(errorValue, "error serializing property value")
			}
			// don't need to escape field name since go field names are already valid
			line := fieldName + kvSeparator + e.Escape(fieldValue)
			if i < len(fieldNames)-1 {
				line = line + lineSeparator
			}
			_, errorWrite := outputWriter.Write([]byte(line))
			if errorWrite != nil {
				return errors.Wrap(errorWrite, "error writing property to output writer")
			}
		}
	} else {
		switch inputObject := inputObject.(type) {
		case string:
			_, err := outputWriter.Write([]byte(inputObject))
			if err != nil {
				return errors.Wrap(err, "error writing property to output writer")
			}
		default:
			inputObjectString, errorValue := valueSerializer(inputObject)
			if errorValue != nil {
				return errors.Wrap(errorValue, "error serializing property key")
			}
			_, errorWrite := outputWriter.Write([]byte(inputObjectString))
			if errorWrite != nil {
				return errors.Wrap(errorWrite, "error writing property to output writer")
			}
		}
	}
	return nil
}
