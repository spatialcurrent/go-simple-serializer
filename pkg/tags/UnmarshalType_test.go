// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

import (
	"reflect"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalTypeMapInterfaceInterface(t *testing.T) {
	obj, err := UnmarshalType([]byte("a=1 b=2 c=3"), '=', reflect.TypeOf(map[interface{}]interface{}{}))
	assert.NoError(t, err)
	assert.Equal(t, obj, map[interface{}]interface{}{"a": "1", "b": "2", "c": "3"})
}

func TestUnmarshalTypeMapStringInterface(t *testing.T) {
	obj, err := UnmarshalType([]byte("a=1 b=2 c=3"), '=', reflect.TypeOf(map[string]interface{}{}))
	assert.NoError(t, err)
	assert.Equal(t, obj, map[string]interface{}{"a": "1", "b": "2", "c": "3"})
}

func TestUnmarshalTypeMapStringString(t *testing.T) {
	obj, err := UnmarshalType([]byte("a=1 b=2 c=3"), '=', reflect.TypeOf(map[string]string{}))
	assert.NoError(t, err)
	assert.Equal(t, obj, map[string]string{"a": "1", "b": "2", "c": "3"})
}
