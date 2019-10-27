// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package jsonl

import (
	"io"
)

// WriteFlusher is a simple interface that wraps io.Writer and Flusher.
type WriteFlusher interface {
	io.Writer
	Flusher
}
