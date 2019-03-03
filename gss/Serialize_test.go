// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"testing"
)

func TestSerialize(t *testing.T) {

	for _, testCase := range serializeTestCases {

		got, err := SerializeString(testCase.Object, testCase.Format, testCase.Header, testCase.Limit)
		if err != nil {
			t.Errorf(errors.Wrap(err, "error running test").Error())
		} else if !reflect.DeepEqual(strings.TrimSpace(got), strings.TrimSpace(testCase.String)) {
			t.Errorf("SerializeString(%v) == %v (%v), want %v (%s)", testCase, got, reflect.TypeOf(got), testCase.String, reflect.TypeOf(testCase.String))
		}
	}

}
