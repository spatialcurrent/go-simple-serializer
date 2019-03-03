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
	"testing"
)

func TestDeserialize(t *testing.T) {

	for _, testCase := range deserializeTestCases {

		got, err := DeserializeString(testCase.String, testCase.Format, testCase.Header, testCase.Comment, testCase.LazyQuotes, testCase.SkipLines, testCase.Limit, testCase.Type, false, false)
		if err != nil {
			t.Errorf(errors.Wrap(err, "error running test").Error())
		} else if !reflect.DeepEqual(got, testCase.Object) {
			t.Errorf("Deserialize(%v) == \n%v\n (%v), want \n%v\n (%s)", testCase, got, reflect.TypeOf(got), testCase.Object, reflect.TypeOf(testCase.Object))
		}
	}

}
