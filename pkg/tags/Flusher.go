// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package tags

// Flusher interfaces is a simple interface that wraps the Flush() function.
type Flusher interface {
	Flush() error
}
