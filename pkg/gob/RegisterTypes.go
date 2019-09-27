// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gob

import (
	stdgob "encoding/gob"
)

// RegisterTypes registers some default types that are not registered by default by the standard library gob package.
func RegisterTypes() {
	stdgob.Register(map[string]interface{}(nil))
	stdgob.Register([]interface{}(nil))
}
