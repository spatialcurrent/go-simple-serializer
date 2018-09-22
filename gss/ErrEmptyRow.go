// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gss

type ErrEmptyRow struct{}

func (e ErrEmptyRow) Error() string {
	return "empty row"
}
