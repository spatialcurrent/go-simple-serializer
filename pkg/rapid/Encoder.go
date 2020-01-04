// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package rapid

import (
	"fmt"
	"io"
)

type Encoder struct {
	writer         io.Writer
	headerPath     string
	headerKey      string
	typeString     string
	typeInt        string
	typeFloat      string
	typeBool       string
	typeArray      string
	separator      string
	trailer        string
	boundaryMarker string
	count          int
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer:         w,
		headerPath:     "P",
		headerKey:      "K",
		typeString:     "S",
		typeInt:        "I",
		typeFloat:      "F",
		typeBool:       "B",
		typeArray:      "A",
		separator:      ":",
		trailer:        "\n",
		boundaryMarker: "---\n",
		count:          0,
	}
}

func (e *Encoder) encodeHeader(header string) (int, error) {
	return fmt.Fprint(e.writer, header)
}

func (e *Encoder) encodeCount(count int) (int, error) {
	return fmt.Fprintf(e.writer, "%d", count)
}

func (e *Encoder) encodeIndex(index int) (int, error) {
	return fmt.Fprintf(e.writer, "%d", index)
}

func (e *Encoder) encodeValue(value interface{}) (int, error) {
	if str, ok := value.(string); ok {
		return fmt.Fprintf(e.writer, "%s%s%s", e.typeString, e.separator, str)
	}
	if b, ok := value.(bool); ok {
		if b {
			return fmt.Fprintf(e.writer, "%s%s1", e.typeBool, e.separator)
		}
		return fmt.Fprintf(e.writer, "%s%s0", e.typeBool, e.separator)
	}
	if i, ok := value.(int); ok {
		return fmt.Fprintf(e.writer, "%s%s%d", e.typeInt, e.separator, i)
	}
	if f, ok := value.(float64); ok {
		return fmt.Fprintf(e.writer, "%s%s%f", e.typeFloat, e.separator, f)
	}
	if f, ok := value.(float32); ok {
		return fmt.Fprintf(e.writer, "%s%s%f", e.typeFloat, e.separator, f)
	}
	if slc, ok := value.([]interface{}); ok {
		return fmt.Fprintf(e.writer, "%s%s%d", e.typeArray, e.separator, len(slc))
	}
	return 0, fmt.Errorf("could not encode value %#v", value)
}

/*
func (e *Encoder) encodeIndicies(indicies []int) error {
	count := len(indicies)
	if _, err := e.encodeCount(count); err != nil {
		return fmt.Errorf("error encoding count %d: %w", count, err)
	}
	if count > 0 {
		if _, err := e.encodeSeparator(); err != nil {
			return fmt.Errorf("error encoding separator : %w", err)
		}
		for i, v := range indicies {
			if _, err := e.encodeIndex(v); err != nil {
				return fmt.Errorf("error encoding index %d: %w", v, err)
			}
			if i < count-1 {
				if _, err := e.encodeSeparator(); err != nil {
					return fmt.Errorf("error encoding separator : %w", err)
				}
			}
		}
	}
	return nil
}*/

func (e *Encoder) encodeSeparator() (int, error) {
	return fmt.Fprint(e.writer, e.separator)
}

func (e *Encoder) encodeTrailer() (int, error) {
	return fmt.Fprint(e.writer, e.trailer)
}

func (e *Encoder) encodeBoundaryMarker() (int, error) {
	return fmt.Fprint(e.writer, e.boundaryMarker)
}

func (e *Encoder) encodeKey(key interface{}) error {
	if _, err := e.encodeHeader(e.headerKey); err != nil {
		return fmt.Errorf("error encoding header %q: %w", e.headerKey, err)
	}
	if _, err := e.encodeSeparator(); err != nil {
		return fmt.Errorf("error encoding separator : %w", err)
	}
	if _, err := e.encodeValue(key); err != nil {
		return fmt.Errorf("error encoding key value %q: %w", key, err)
	}
	if _, err := e.encodeTrailer(); err != nil {
		return fmt.Errorf("error encoding trailer: %w", err)
	}
	return nil
}

func (e *Encoder) encodeIndicies(h string, indicies []int) error {
	if _, err := e.encodeHeader(h); err != nil {
		return fmt.Errorf("error encoding header %q: %w", h, err)
	}
	if len(indicies) > 0 {
		if _, err := e.encodeSeparator(); err != nil {
			return fmt.Errorf("error encoding separator : %w", err)
		}
		for i, v := range indicies {
			if _, err := e.encodeIndex(v); err != nil {
				return fmt.Errorf("error encoding index %d: %w", v, err)
			}
			if i < len(indicies)-1 {
				if _, err := e.encodeSeparator(); err != nil {
					return fmt.Errorf("error encoding separator : %w", err)
				}
			}
		}
	}
	if _, err := e.encodeTrailer(); err != nil {
		return fmt.Errorf("error encoding trailer: %w", err)
	}
	return nil
}

/*
func (e *Encoder) encodeFeature(paths [][]int) error {
	if _, err := e.encodeHeader(e.headerFeature); err != nil {
		return fmt.Errorf("error encoding header %q: %w", e.headerFeature, err)
	}
	if _, err := e.encodeSeparator(); err != nil {
		return fmt.Errorf("error encoding separator : %w", err)
	}
	if _, err := e.encodeCount(len(paths)); err != nil {
		return fmt.Errorf("error encoding count %d: %w", len(paths), err)
	}
	if len(paths) > 0 {
		if _, err := e.encodeSeparator(); err != nil {
			return fmt.Errorf("error encoding separator : %w", err)
		}
		for i, path := range paths {
			if _, err := e.encodeHeader(e.headerPath); err != nil {
				return fmt.Errorf("error encoding header %q: %w", e.headerPath, err)
			}
			if _, err := e.encodeSeparator(); err != nil {
				return fmt.Errorf("error encoding separator : %w", err)
			}
			if err := e.encodeIndicies(path); err != nil {
				return fmt.Errorf("error encoding indicies %q: %w", path, err)
			}
			if i < len(paths)-1 {
				if _, err := e.encodeSeparator(); err != nil {
					return fmt.Errorf("error encoding separator : %w", err)
				}
			}
		}
	}
	_, err := e.encodeTrailer()
	if err != nil {
		return fmt.Errorf("error encoding trailer: %w", err)
	}
	return nil
}
*/

func (e *Encoder) encodeLiteral(d *Dictionary, prefix []int, literal interface{}) {
	if i := d.GetIndex(literal); i != -1 {
		e.encodeIndicies(e.headerPath, append(prefix, i))
	} else {
		i := d.AddKey(literal)
		e.encodeKey(literal)
		e.encodeIndicies(e.headerPath, append(prefix, i))
	}
}

func (e *Encoder) encodeObject(d *Dictionary, prefix []int, obj interface{}) error {

	if _, ok := obj.(string); ok {
		e.encodeLiteral(d, prefix, obj)
	}

	if _, ok := obj.(int); ok {
		e.encodeLiteral(d, prefix, obj)
	}

	if _, ok := obj.(float64); ok {
		e.encodeLiteral(d, prefix, obj)
	}

	if m, ok := obj.(map[string]string); ok {
		for k, v := range m {
			if i := d.GetIndex(k); i != -1 {
				e.encodeObject(d, append(prefix, i), v)
			} else {
				i := d.AddKey(k)
				e.encodeKey(k)
				e.encodeObject(d, append(prefix, i), v)
			}
		}
		return nil
	}

	if m, ok := obj.(map[string]interface{}); ok {
		for k, v := range m {
			if i := d.GetIndex(k); i != -1 {
				e.encodeObject(d, append(prefix, i), v)
			} else {
				i := d.AddKey(k)
				e.encodeKey(k)
				e.encodeObject(d, append(prefix, i), v)
			}
		}
		return nil
	}

	if slc, ok := obj.([]interface{}); ok {
		i := d.AddKey(make([]interface{}, 0)) // just add an empty slice to the dictionary to reserve the index number
		e.encodeKey(slc)                      // this will write K:A:<length of slice>
		for position, v := range slc {
			p := d.AddKey(position)
			e.encodeKey(position)
			e.encodeObject(d, append(prefix, i, p), v)
		}
	}

	return nil
}

func (e *Encoder) Encode(v interface{}) error {
	if e.count > 0 {
		e.encodeBoundaryMarker()
	}
	d := NewDictionary()
	e.encodeObject(d, make([]int, 0), v)
	e.count += 1
	return nil
}

// Flush flushes the underlying writer, if it has a Flush method.
// This writer itself does no buffering.
func (e *Encoder) Flush() error {
	if flusher, ok := e.writer.(interface{ Flush() error }); ok {
		err := flusher.Flush()
		if err != nil {
			return fmt.Errorf("error flushing underlying writer: %w", err)
		}
	}
	return nil
}

// Close closes the underlying writer, if it has a Close method.
func (e *Encoder) Close() error {
	if closer, ok := e.writer.(interface{ Close() error }); ok {
		err := closer.Close()
		if err != nil {
			return fmt.Errorf("error closing underlying writer: %w", err)
		}
	}
	return nil
}
