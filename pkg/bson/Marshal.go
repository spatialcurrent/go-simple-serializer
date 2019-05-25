// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package bson

import (
	// import the mgo bson library
	mgobson "gopkg.in/mgo.v2/bson"
)

import (
	"github.com/spatialcurrent/go-stringify/pkg/stringify"
)

const (
	prefix = ""
	indent = "  "
)

// Marshal formats an object into a slice of bytes of BSON.
func Marshal(obj interface{}) ([]byte, error) {
	return mgobson.Marshal(stringify.StringifyMapKeys(obj))
}
