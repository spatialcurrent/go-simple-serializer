// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package rapid

import (
	//"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spatialcurrent/go-simple-serializer/pkg/escaper"
)

type Encoder struct {
	writer                io.Writer
	dictionary            *Dictionary
	escaper               *escaper.Escaper
	headerPath            string
	headerKey             string
	typeString            string
	typeInt               string
	typeFloat             string
	typeBool              string
	typeArray             string
	separator             string
	trailer               string
	boundaryMarker        string
	resetMarker           string
	sortKeys              bool
	writeKeysBreadthFirst bool
	buffer                bool
	dictionaryLimit       int
	charset []rune
	keysWritten           int
	count                 int
}

func NewEncoder(w io.Writer) *Encoder {
	charset := make([]rune, 0)
	for i := 0; i < 256; i++ {
		if r := rune(i); r != ':' && r != '\'' && r != '\n' {
			charset = append(charset, r)
		}
	}
	return &Encoder{
		writer:                w,
		dictionary:            NewDictionary(),
		escaper:               escaper.New().Prefix("\\").Sub("\\", ":"),
		headerPath:            "P",
		headerKey:             "K",
		typeString:            "S",
		typeInt:               "I",
		typeFloat:             "F",
		typeBool:              "B",
		typeArray:             "A",
		separator:             ":",
		trailer:               "\n",
		boundaryMarker:        "---\n",
		resetMarker:           "===\n",
		sortKeys:              true,
		writeKeysBreadthFirst: true,
		buffer:                true,
		dictionaryLimit:       10,
		keysWritten:           0,
		count:                 0,
		charset: charset,
	}
}

func (e *Encoder) encodeHeader(header string) (int, error) {
	return fmt.Fprint(e.writer, header)
}

func (e *Encoder) encodeCount(count int) (int, error) {
	return fmt.Fprintf(e.writer, "%d", count)
}

func (e *Encoder) encodeIndex(index int) (int, error) {
	//return fmt.Fprintf(e.writer, "%d", index)
	if index >= len(e.charset) {
		return 0, fmt.Errorf("cannot encode index greater than %d", len(e.charset))
	}
	return fmt.Fprintf(e.writer, "%s", string(e.charset[index]))
}

func (e *Encoder) encodeValue(value interface{}, escape bool) (int, error) {
	if str, ok := value.(string); ok {
		if escape {
			return fmt.Fprintf(e.writer, "%s%s%s", e.typeString, e.separator, e.escaper.Escape(str))
		}
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

func (e *Encoder) encodeResetMarker() (int, error) {
	return fmt.Fprint(e.writer, e.resetMarker)
}

func (e *Encoder) encodeKey(key interface{}, escape bool) error {
	if _, err := e.encodeHeader(e.headerKey); err != nil {
		return fmt.Errorf("error encoding header %q: %w", e.headerKey, err)
	}
	/*if _, err := e.encodeSeparator(); err != nil {
		return fmt.Errorf("error encoding separator : %w", err)
	}*/
	if _, err := e.encodeValue(key, escape); err != nil {
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
		for _, v := range indicies {
			if _, err := e.encodeIndex(v); err != nil {
				return fmt.Errorf("error encoding index %d: %w", v, err)
			}
			/*if i < len(indicies)-1 {
				if _, err := e.encodeSeparator(); err != nil {
					return fmt.Errorf("error encoding separator : %w", err)
				}
			}*/
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
		e.encodeKey(literal, false)
		e.encodeIndicies(e.headerPath, append(prefix, i))
	}
}

func (e *Encoder) keysFromMapStringString(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	if e.sortKeys {
		sort.Strings(keys)
	}
	return keys
}

func (e *Encoder) keysFromMapStringInterface(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	if e.sortKeys {
		sort.Strings(keys)
	}
	return keys
}

func (e *Encoder) encodeObject(d *Dictionary, prefix []int, obj interface{}) error {

	if _, ok := obj.(bool); ok {
		e.encodeLiteral(d, prefix, obj)
	}

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
		for _, k := range e.keysFromMapStringString(m) {
			if i := d.GetIndex(k); i != -1 {
				e.encodeObject(d, append(prefix, i), m[k])
			} else {
				i := d.AddKey(k)
				e.encodeKey(k, false)
				e.encodeObject(d, append(prefix, i), m[k])
			}
		}
		return nil
	}

	if m, ok := obj.(map[string]interface{}); ok {
		for _, k := range e.keysFromMapStringInterface(m) {
			if i := d.GetIndex(k); i != -1 {
				e.encodeObject(d, append(prefix, i), m[k])
			} else {
				i := d.AddKey(k)
				e.encodeKey(k, false)
				e.encodeObject(d, append(prefix, i), m[k])
			}
		}
		return nil
	}

	if slc, ok := obj.([]interface{}); ok {
		i := d.AddKey(make([]interface{}, 0)) // just add an empty slice to the dictionary to reserve the index number
		e.encodeKey(slc, false)               // this will write K:A:<length of slice>
		for position, v := range slc {
			p := d.AddKey(position)
			e.encodeKey(position, false)
			e.encodeObject(d, append(prefix, i, p), v)
		}
	}

	return nil
}

func (e *Encoder) getKeyType(value interface{}) string {
	if _, ok := value.(bool); ok {
		return e.typeBool
	}
	if _, ok := value.(string); ok {
		return e.typeString
	}
	if _, ok := value.(int); ok {
		return e.typeInt
	}
	if _, ok := value.(float64); ok {
		return e.typeFloat
	}
	if _, ok := value.(float32); ok {
		return e.typeFloat
	}
	if _, ok := value.([]interface{}); ok {
		return e.typeArray
	}
	return ""
}

func (e *Encoder) encodeKeyValue(value interface{}, escape bool) (int, error) {
	if str, ok := value.(string); ok {
		if escape {
			return fmt.Fprint(e.writer, strings.ReplaceAll(e.escaper.Escape(str), "\n", "\\n"))
		}
		return fmt.Fprint(e.writer, strings.ReplaceAll(str, "\n", "\\n"))
	}
	if b, ok := value.(bool); ok {
		if b {
			return fmt.Fprint(e.writer, "1")
		}
		return fmt.Fprint(e.writer, "0")
	}
	if i, ok := value.(int); ok {
		return fmt.Fprintf(e.writer, "%d", i)
	}
	if f, ok := value.(float64); ok {
		return fmt.Fprintf(e.writer, "%f", f)
	}
	if f, ok := value.(float32); ok {
		return fmt.Fprintf(e.writer, "%f", f)
	}
	if slc, ok := value.([]interface{}); ok {
		return fmt.Fprintf(e.writer, "%d", len(slc))
	}
	return 0, fmt.Errorf("could not encode value %#v", value)
}

func (e *Encoder) encodeKeys(keys []interface{}) error {
	if len(keys) > 0 {
		headerWritten := false
		currentType := ""
		currentPosition := 0
		for i, key := range keys {

			nextType := e.getKeyType(key)

			// if there is a type change
			if i > 0 && (nextType != currentType || nextType == e.typeArray) {
				if _, err := e.encodeTrailer(); err != nil {
					return fmt.Errorf("error encoding trailer: %w", err)
				}
				currentPosition = 0
				headerWritten = false
			}

			// if the header has not been written or there is a type change
			if (!headerWritten) || nextType != currentType || nextType == e.typeArray {
				if _, err := fmt.Fprintf(e.writer, "%s%s%s", e.headerKey, nextType, e.separator); err != nil {
					return fmt.Errorf("error encoding key %q: %w", key, err)
				}
				headerWritten = true
			}

			if currentPosition > 0 {
				if _, err := fmt.Fprintf(e.writer, "%s", e.separator); err != nil {
					return fmt.Errorf("error encoding separator %q: %w", key, err)
				}
			}

			if _, err := e.encodeKeyValue(key, true); err != nil {
				return fmt.Errorf("error encoding key value %#v: %w", key, err)
			}

			currentPosition += 1
			currentType = nextType
		}
		// write final trailer
		if _, err := e.encodeTrailer(); err != nil {
			return fmt.Errorf("error encoding trailer: %w", err)
		}
	}
	return nil
}

func (e *Encoder) encodePaths(paths [][]int) error {
	for _, path := range paths {
		err := e.encodeIndicies(e.headerPath, path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) parseObject(d *Dictionary, prefix []int, obj interface{}) [][]int {

	if _, ok := obj.(bool); ok {
		if i := d.GetIndex(obj); i != -1 {
			return [][]int{append(prefix, i)}
		}
		return [][]int{append(prefix, d.AddKey(obj))}
	}

	if _, ok := obj.(string); ok {
		if i := d.GetIndex(obj); i != -1 {
			return [][]int{append(prefix, i)}
		}
		return [][]int{append(prefix, d.AddKey(obj))}
	}

	if _, ok := obj.(int); ok {
		if i := d.GetIndex(obj); i != -1 {
			return [][]int{append(prefix, i)}
		}
		return [][]int{append(prefix, d.AddKey(obj))}
	}

	if _, ok := obj.(float64); ok {
		if i := d.GetIndex(obj); i != -1 {
			return [][]int{append(prefix, i)}
		}
		return [][]int{append(prefix, d.AddKey(obj))}
	}

	if m, ok := obj.(map[string]string); ok {
		paths := make([][]int, 0, len(m))
		if e.writeKeysBreadthFirst {
			keys := e.keysFromMapStringString(m)
			// parse keys
			for _, k := range keys {
				if !d.HasKey(k) {
					d.AddKey(k)
				}
			}
			// parse values
			for _, k := range keys {
				paths = append(paths, e.parseObject(d, append(prefix, d.GetIndex(k)), m[k])...)
			}
		} else {
			for _, k := range e.keysFromMapStringString(m) {
				if i := d.GetIndex(k); i != -1 {
					paths = append(paths, e.parseObject(d, append(prefix, i), m[k])...)
				} else {
					paths = append(paths, e.parseObject(d, append(prefix, d.AddKey(k)), m[k])...)
				}
			}
		}
		return paths
	}

	if m, ok := obj.(map[string]interface{}); ok {
		paths := make([][]int, 0, len(m))
		if e.writeKeysBreadthFirst {
			keys := e.keysFromMapStringInterface(m)
			// parse keys
			for _, k := range keys {
				if !d.HasKey(k) {
					d.AddKey(k)
				}
			}
			// parse values
			for _, k := range keys {
				paths = append(paths, e.parseObject(d, append(prefix, d.GetIndex(k)), m[k])...)
			}
		} else {
			for _, k := range e.keysFromMapStringInterface(m) {
				if i := d.GetIndex(k); i != -1 {
					paths = append(paths, e.parseObject(d, append(prefix, i), m[k])...)
				} else {
					paths = append(paths, e.parseObject(d, append(prefix, d.AddKey(k)), m[k])...)
				}
			}
		}
		return paths
	}

	if slc, ok := obj.([]interface{}); ok {
		paths := make([][]int, 0, len(slc))
		i := d.AddKey(slc) // if using buffer we need to add the full array, so it can be encoded later

		if e.writeKeysBreadthFirst {
			for position, _ := range slc {
				if i := d.GetIndex(position); i == -1 {
					d.AddKey(position)
				}
			}
			for position, v := range slc {
				paths = append(paths, e.parseObject(d, append(prefix, i, d.GetIndex(position)), v)...)
			}
		} else {
			for position, v := range slc {
				paths = append(paths, e.parseObject(d, append(prefix, i, d.AddKey(position)), v)...)
			}
		}
		return paths
	}

	return make([][]int, 0)
}

func (e *Encoder) Encode(v interface{}) error {
	if e.count%e.dictionaryLimit == 0 {
		if e.count > 0 {
			e.encodeResetMarker()
		}
		// reuse the existing dictionary
		e.dictionary.Reset()
		// reset number of keys written
		e.keysWritten = 0
	} else {
		if e.count > 0 {
			e.encodeBoundaryMarker()
		}
	}
	if e.buffer {
		paths := e.parseObject(e.dictionary, make([]int, 0), v)
		if newKeys := e.dictionary.Keys()[e.keysWritten:]; len(newKeys) > 0 {
			err := e.encodeKeys(newKeys)
			if err != nil {
				return err
			}
		}
		if err := e.encodePaths(paths); err != nil {
			return fmt.Errorf("error encoding paths for %#v: %w", v, err)
		}
		// set the number of keys written to the dictionary size
		e.keysWritten = len(e.dictionary.Keys())
	} else {
		err := e.encodeObject(e.dictionary, make([]int, 0), v)
		if err != nil {
			return err
		}
	}
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
