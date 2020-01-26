// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package scanner

type multiScanner struct {
	scanners []Scanner
	eof      bool
}

func (ms *multiScanner) Buffer(buf []byte, max int) {
	for _, s := range ms.scanners {
		s.Buffer(buf, max)
	}
}

func (ms *multiScanner) Bytes() []byte {
	if ms.eof {
		return make([]byte, 0)
	}
	if err := ms.scanners[0].Err(); err != nil {
		return make([]byte, 0)
	}
	if len(ms.scanners) == 0 {
		return make([]byte, 0)
	}
	return ms.scanners[0].Bytes()
}

func (ms *multiScanner) Err() error {
	return ms.scanners[0].Err()
}

func (ms *multiScanner) Scan() bool {
	if ms.eof {
		return false
	}
	if err := ms.scanners[0].Err(); err != nil {
		return false
	}
	if len(ms.scanners) == 0 {
		return false
	}
	valid := ms.scanners[0].Scan()
	if err := ms.scanners[0].Err(); err != nil {
		return valid
	}
	if len(ms.scanners) == 1 {
		if !valid {
			ms.eof = true
			return false
		}
		return true
	}
	if !valid {
		ms.scanners = ms.scanners[1:]
		return ms.Scan()
	}
	return true
}

func (ms *multiScanner) Text() string {
	if ms.eof {
		return ""
	}
	if err := ms.scanners[0].Err(); err != nil {
		return ""
	}
	return ms.scanners[0].Text()
}

// MultiScanner returns a new Scanner that reads from a series of scanners sequentially.
func MultiScanner(scanners ...Scanner) Scanner {
	s := make([]Scanner, len(scanners))
	copy(s, scanners)
	return &multiScanner{
		scanners: s,
		eof:      false,
	}
}
