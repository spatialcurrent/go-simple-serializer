// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package rapid

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/spatialcurrent/go-simple-serializer/pkg/splitter"
)

type Decoder struct {
	dictionary     *Dictionary
	scanner        *bufio.Scanner
	headerKey      byte
	headerPath     byte
	headerComment  byte
	typeBool       byte
	typeString     byte
	typeInt        byte
	typeFloat      byte
	typeArray      byte
	separator      byte
	boundaryMarker []byte
	valueTrue      byte
	valueFalse     byte
	eof            bool
}

func NewDecoder(r io.Reader, lineSeparator byte, dropCR bool) *Decoder {
	scanner := bufio.NewScanner(r)
	scanner.Split(splitter.ScanLines(lineSeparator, dropCR))
	return &Decoder{
		dictionary:     NewDictionary(),
		scanner:        scanner,
		headerKey:      'K',
		headerPath:     'P',
		headerComment:  '#',
		typeBool:       'B',
		typeString:     'S',
		typeInt:        'I',
		typeFloat:      'F',
		typeArray:      'A',
		separator:      ':',
		boundaryMarker: []byte("---"),
		valueTrue:      '1',
		valueFalse:     '0',
		eof:            false,
	}
}

func (d *Decoder) add(v interface{}, chain []interface{}) error {
	if len(chain) == 0 {
		return nil
	}

	if m, ok := v.(map[string]interface{}); ok {
		if len(chain) < 2 {
			return fmt.Errorf("invalid chain for map[string]interface{}: %#v", chain)
		}
		if len(chain) == 2 {
			m[fmt.Sprint(chain[0])] = chain[1]
			return nil
		}
		if len(chain) > 2 {
			k := fmt.Sprint(chain[0])
			n, ok := m[k]
			if ok {
				if slc, ok := chain[1].([]interface{}); ok {
					d.add(slc, chain[2:])
				} else {
					d.add(n, chain[1:])
				}
			} else {
				if slc, ok := chain[1].([]interface{}); ok {
					d.add(slc, chain[2:])
					m[k] = slc
				} else {
					n := map[string]interface{}{}
					d.add(n, chain[1:])
					m[k] = n
				}
			}
		}
	}

	if slc, ok := v.([]interface{}); ok {
		if len(chain) < 2 {
			return fmt.Errorf("invalid chain for []interface{}: %#v", chain)
		}
		if len(chain) == 2 {
			i, ok := chain[0].(int)
			if !ok {
				return fmt.Errorf("error decoding array index from %q", fmt.Sprint(chain[0]))
			}
			slc[i] = chain[1]
			return nil
		}
		/*if len(chain) > 2 {
			i, err := strconv.Atoi(fmt.Sprint(chain[0]))
			if err != nil {
				return fmt.Errorf("error decoding array index from %q: %w", fmt.Sprint(chain[0]), err)
			}
			x := slc[i]
			if x != nil {
				d.add(x, chain[1:])
			} else {
				n := map[string]interface{}{}
				d.add(n, chain[1:])
				slc[i] = n
			}
		}*/
	}
	return nil
}

func (d *Decoder) Decode(v interface{}) error {

	if d.eof {
		return io.EOF
	}

	if m, ok := v.(*map[string]interface{}); ok {

		values := map[string]interface{}{}

		dict := NewDictionary()

		eof := true
		for d.scanner.Scan() {
			b := d.scanner.Bytes()
			if bytes.Equal(b, d.boundaryMarker) {
				eof = false
				break
			}
			if len(b) == 0 {
				continue
			}

			switch b[0] {
			case d.headerComment:
				continue
			case d.headerKey:
				i := bytes.Index(b[1:], []byte{d.separator})
				if i != 0 {
					return fmt.Errorf("invalid byte, expecting separator after key header: %q", string(b))
				}
				j := bytes.Index(b[2:], []byte{d.separator})
				if j != 1 {
					return fmt.Errorf("invalid index %d, expecting separator after key type: %q", j, string(b))
				}
				switch b[2] {
				case d.typeBool:
					if b[4] == d.valueFalse {
						dict.AddKey(false)
					} else {
						dict.AddKey(true)
					}
				case d.typeString:
					dict.AddKey(string(b[4:]))
				case d.typeInt:
					i, err := strconv.Atoi(string(b[4:]))
					if err != nil {
						return fmt.Errorf("invalid integer %q: %w", string(b[4:]), err)
					}
					dict.AddKey(i)
				case d.typeFloat:
					f, err := strconv.ParseFloat(string(b[4:]), 64)
					if err != nil {
						return fmt.Errorf("invalid integer %q: %w", string(b[4:]), err)
					}
					dict.AddKey(f)
				case d.typeArray:
					length, err := strconv.Atoi(string(b[4:]))
					if err != nil {
						return fmt.Errorf("invalid array size %q: %w", string(b[4:]), err)
					}
					dict.AddKey(make([]interface{}, length)) // add a
				}
			case d.headerPath:
				i := bytes.Index(b[1:], []byte{d.separator})
				if i != 0 {
					return fmt.Errorf("invalid byte, expecting separator after path header: %q", string(b))
				}
				chain := make([]interface{}, 0)
				for _, str := range bytes.Split(b[2:], []byte{d.separator}) {
					i, err := strconv.Atoi(string(str))
					if err != nil {
						return fmt.Errorf("invalid index: %w", err)
					}
					v, ok := dict.GetKey(i)
					if !ok {
						return fmt.Errorf("invalid index %d", i)
					}
					chain = append(chain, v)
				}
				err := d.add(values, chain)
				if err != nil {
					return fmt.Errorf("error adding chain: %q", string(b))
				}
			}
		}

		// if scanner returns an error then return it
		if err := d.scanner.Err(); err != nil {
			d.eof = true
			return err
		}

		d.eof = eof

		// set output
		*m = values

		return nil
	}
	return fmt.Errorf("error decoding %#v, unsupported type %T", v, v)
}
