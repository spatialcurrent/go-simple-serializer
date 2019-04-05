// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
)

import (
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// DeserializeJSONL deserializes the input JSON lines bytes into a Go object.
//  - https://golang.org/pkg/encoding/json/
func DeserializeJSONL(input io.Reader, inputComment string, inputSkipLines int, inputLimit int, outputType reflect.Type, async bool) (interface{}, error) {

	output := reflect.MakeSlice(outputType, 0, 0)
	if inputLimit == 0 {
		return output.Interface(), nil
	}

	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)
	for i := 0; i < inputSkipLines; i++ {
		if !scanner.Scan() {
			break
		}
	}

	if async {

		objects := make(chan reflect.Value, 10000)
		waitGroupObjects := &sync.WaitGroup{}
		waitGroupObjects.Add(1)
		go func() {
			for object := range objects {
				output = reflect.Append(output, object)
			}
			waitGroupObjects.Done()
		}()

		inputCount := 0
		errGroupLines, _ := errgroup.WithContext(context.Background())
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(inputComment) == 0 || !strings.HasPrefix(line, inputComment) {
				errGroupLines.Go(func() error {
					lineType := GetTypeJSON(line)
					var ptr reflect.Value
					if lineType.Kind() == reflect.Array || lineType.Kind() == reflect.Slice {
						ptr = reflect.New(lineType)
						ptr.Elem().Set(reflect.MakeSlice(lineType, 0, 0))
					} else if lineType.Kind() == reflect.Map {
						ptr = reflect.New(lineType)
						ptr.Elem().Set(reflect.MakeMap(lineType))
					} else {
						return errors.New("error creating object for JSON line " + line)
					}
					err := json.Unmarshal([]byte(line), ptr.Interface())
					if err != nil {
						return errors.Wrap(err, "Error reading object from JSON line")
					}
					objects <- ptr.Elem()
					return nil
				})
				inputCount += 1
				if inputLimit > 0 && inputCount >= inputLimit {
					break
				}
			}
		}
		err := errGroupLines.Wait()
		close(objects)
		waitGroupObjects.Wait()
		return output.Interface(), err
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(inputComment) == 0 || !strings.HasPrefix(line, inputComment) {
			lineType := GetTypeJSON(line)
			var ptr reflect.Value
			if lineType.Kind() == reflect.Array || lineType.Kind() == reflect.Slice {
				ptr = reflect.New(lineType)
				ptr.Elem().Set(reflect.MakeSlice(lineType, 0, 0))
			} else if lineType.Kind() == reflect.Map {
				ptr = reflect.New(lineType)
				ptr.Elem().Set(reflect.MakeMap(lineType))
			} else {
				return nil, errors.New("error creating object for JSON line " + line)
			}
			err := json.Unmarshal([]byte(line), ptr.Interface())
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("Error reading object from JSON line with content %q", line))
			}
			output = reflect.Append(output, ptr.Elem())
			if inputLimit > 0 && output.Len() >= inputLimit {
				break
			}
		}
	}

	return output.Interface(), nil
}
