// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package sv

import (
	"sort"
)

type Row []string

func (r Row) Index(str string) int {
	for i, element := range r {
		if element == str {
			return i
		}
	}
	return -1
}

func (r Row) Sort(reversed bool) {
	sort.Slice(r, func(i, j int) bool {
		if reversed {
			return r[i] > r[j]
		}
		return r[i] < r[j]
	})
}

func (r Row) Transform(m map[int]int) Row {
	newRow := make([]string, len(r))
	for i, str := range r {
		newRow[m[i]] = str
	}
	return newRow
}

func (r Row) FillRight(n int) Row {
	return FillRight(append(make([]string, 0), r...), n)
}
