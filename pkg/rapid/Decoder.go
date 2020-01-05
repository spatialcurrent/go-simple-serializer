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
	//"strings"

	"github.com/spatialcurrent/go-simple-serializer/pkg/escaper"
	"github.com/spatialcurrent/go-simple-serializer/pkg/splitter"
)

type Decoder struct {
	dictionary     *Dictionary
	scanner        *bufio.Scanner
	iterator       *escaper.Iterator
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
	resetMarker    []byte
	valueTrue      byte
	valueFalse     byte
	charset        map[rune]int
	eof            bool
}

func NewDecoder(r io.Reader, lineSeparator byte, dropCR bool) *Decoder {
	charset := map[rune]int{}
	for i := 0; i < 256; i++ {
		if r := rune(i); r != ':' && r != '\'' && r != '\n' {
			charset[r] = len(charset)
		}
	}
	scanner := bufio.NewScanner(r)
	scanner.Split(splitter.ScanLines(lineSeparator, dropCR))
	return &Decoder{
		dictionary: NewDictionary(),
		scanner:    scanner,
		iterator: escaper.NewIterator(
			nil,
			[]byte("\\"),
			[][]byte{
				[]byte("\\"),
				[]byte(":"),
			},
			[]byte(":"),
		),
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
		resetMarker: []byte("==="),
		valueTrue:      '1',
		valueFalse:     '0',
		charset: charset,
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
			if str, ok := chain[0].(string); ok {
				m[str] = chain[1]
			}
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
	}

	return nil
}

func (d *Decoder) Decode(v interface{}) error {

	if d.eof {
		return io.EOF
	}

	if m, ok := v.(*map[string]interface{}); ok {

		values := *m

		//d.dictionary.Reset()

		eof := true
		for d.scanner.Scan() {

			b := d.scanner.Bytes()

			if len(b) == 0 {
				continue
			}

			if bytes.Equal(b, d.boundaryMarker) {
				eof = false
				break
			}

			if bytes.Equal(b, d.resetMarker) {
				d.dictionary.Reset()
				eof = false
				break
			}

			switch b[0] {
			case d.headerComment:
				continue
			case d.headerKey:
				if i := bytes.Index(b, []byte{d.separator}); i != 2 {
					return fmt.Errorf("invalid index %d, expecting separator after key type: %q", i, string(b))
				}
				switch b[1] {
				case d.typeBool:
					if b[3] == d.valueFalse {
						d.dictionary.AddKey(false)
					} else {
						d.dictionary.AddKey(true)
					}
				case d.typeString:
					d.iterator.Reset(bytes.NewReader(b[3:]))
					for {
						key, err := d.iterator.Next()
						if err != nil {
							if err == io.EOF {
								break
							}
							return fmt.Errorf("error iterating through strings %q: %w", string(b[3:]), err)
						}
						d.dictionary.AddKey(string(bytes.ReplaceAll(key, []byte("\\n"), []byte("\n"))))
						//d.dictionary.AddKey(strings.ReplaceAll(string(key), "\\n", "\n"))
					}
				case d.typeInt:
					for _, key := range bytes.Split(b[3:], []byte{d.separator}) {
						i, err := strconv.Atoi(string(key))
						if err != nil {
							return fmt.Errorf("invalid integer %q: %w", string(b[3:]), err)
						}
						d.dictionary.AddKey(i)
					}
				case d.typeFloat:
					for _, key := range bytes.Split(b[3:], []byte{d.separator}) {
						f, err := strconv.ParseFloat(string(key), 64)
						if err != nil {
							return fmt.Errorf("invalid integer %q: %w", string(b[3:]), err)
						}
						d.dictionary.AddKey(f)
					}
				case d.typeArray:
					length, err := strconv.Atoi(string(b[3:]))
					if err != nil {
						return fmt.Errorf("invalid array size %q: %w", string(b[4:]), err)
					}
					d.dictionary.AddKey(make([]interface{}, length)) // add a
				}
			case d.headerPath:
				i := bytes.Index(b[1:], []byte{d.separator})
				if i != 0 {
					return fmt.Errorf("invalid byte, expecting separator after path header: %q", string(b))
				}
				//
				parts := bytes.Split(b[2:], []byte{d.separator})
				//
				chain := make([]interface{}, 0, len(parts))
				//
				for _, str := range parts {
					/*
					i, err := strconv.Atoi([]byte(string(str)))
					if err != nil {
						return fmt.Errorf("invalid path index: %w", err)
					}
					*/
					i := d.charset[rune(str[0])]
					v, ok := d.dictionary.GetKey(i)
					if !ok {
						return fmt.Errorf("invalid path index %d: key not present in dictionary with keys %#v", i, d.dictionary.Keys())
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

		return nil
	}

	return fmt.Errorf("error decoding %#v, unsupported type %T", v, v)
}
