package rapid

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
)

type LiteralEncoder struct {
	writer        io.Writer
	typeBool      byte
	typeInt8      byte
	typeInt64     byte
	typeFloat64   byte
	typeString    byte
	typeArray     byte
	typeMap       byte
	typeInterface byte
	order         binary.ByteOrder
}

func NewLiteralEncoder(w io.Writer) *LiteralEncoder {
	return &LiteralEncoder{
		writer:        w,
		typeBool:      DefaultTypeBool,
		typeInt8:      DefaultTypeInt8,
		typeInt64:     DefaultTypeInt64,
		typeFloat64:   DefaultTypeFloat64,
		typeString:    DefaultTypeString,
		typeArray:     DefaultTypeArray,
		typeMap:       DefaultTypeMap,
		typeInterface: DefaultTypeInterface,
		order:         binary.LittleEndian,
	}
}

func (e *LiteralEncoder) Reset(w io.Writer) {
	e.writer = w
}

func (e *LiteralEncoder) max(v []int) int {
	max := 0
	for _, x := range v {
		if x > max {
			max = x
		}
	}
	return max
}

func (e *LiteralEncoder) writeByte(v byte) error {
	_, err := e.writer.Write([]byte{v})
	if err != nil {
		return fmt.Errorf("error encoding byte value %v: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) writeBytes(v ...byte) error {
	_, err := e.writer.Write(v)
	if err != nil {
		return fmt.Errorf("error encoding bytes %x: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) writeStrings(v []string) error {
	for i, str := range v {
		str = strings.ReplaceAll(str, "\n", "\\n")
		//
		b := make([]byte, 1+len(str))
		//
		copy(b[0:len(str)], []byte(str))
		//
		b[len(b)-1] = byte(0)
		//
		if _, err := e.writer.Write(b); err != nil {
			return fmt.Errorf("error encoding string %d of slice with value %v: %w", i, v, err)
		}
	}
	return nil
}

func (e *LiteralEncoder) writeInterfaces(v []interface{}) error {
	for i, x := range v {
		err := e.Encode(x)
		if err != nil {
			return fmt.Errorf("error encoding element %d in %#v: %w", i, v, err)
		}
	}
	return nil
}

func (e *LiteralEncoder) encodeBool(v bool) error {
	if v {
		err := e.writeBytes(e.typeBool, 1)
		if err != nil {
			return fmt.Errorf("error encoding bool value %v: %w", v, err)
		}
		return nil
	}
	err := e.writeBytes(e.typeBool, 0)
	if err != nil {
		return fmt.Errorf("error encoding bool value %v: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeInt8(v int8) error {
	_, err := e.writer.Write([]byte{e.typeInt8, uint8(v)})
	if err != nil {
		return fmt.Errorf("error encoding int8 value %d: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeInt64(v int64) error {
	b := make([]byte, 1+8)
	b[0] = e.typeInt64
	e.order.PutUint64(b[1:], uint64(v))
	_, err := e.writer.Write(b)
	if err != nil {
		return fmt.Errorf("error encoding uint64 value %d: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeInt(v int) error {
	if v < 256 {
		return e.encodeInt8(int8(v))
	}
	return e.encodeInt64(int64(v))
}

func (e *LiteralEncoder) encodeFloat64(v float64) error {
	b := make([]byte, 1+8)
	b[0] = e.typeFloat64
	e.order.PutUint64(b[1:], math.Float64bits(v))
	_, err := e.writer.Write(b)
	if err != nil {
		return fmt.Errorf("error encoding float64 value %v: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeString(v string) error {
	if strings.Contains(v, string([]byte{byte(0)})) {
		return fmt.Errorf("error encoding string with null byte: %q", v)
	}
	v = strings.ReplaceAll(v, "\n", "\\n")
	//
	b := make([]byte, 2+len(v))
	//
	b[0] = e.typeString
	//
	copy(b[1:len(v)+1], []byte(v))
	//
	b[len(b)-1] = byte(0)
	//
	if _, err := e.writer.Write(b); err != nil {
		return fmt.Errorf("error encoding string value %v: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeSliceInterface(v []interface{}) error {
	if err := e.writeByte(e.typeArray | e.typeInterface); err != nil {
		return fmt.Errorf("error encoding string value %v: %w", v, err)
	}
	if err := e.encodeInt(len(v)); err != nil {
		return fmt.Errorf("error encoding array size %d: %w", len(v), err)
	}
	if err := e.writeInterfaces(v); err != nil {
		return fmt.Errorf("error encoding slice of interfaces %#v: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeSliceBool(v []bool) error {
	if err := e.writeByte(e.typeArray | e.typeBool); err != nil {
		return fmt.Errorf("error encoding bool array value %v: %w", v, err)
	}
	if err := e.encodeInt(len(v)); err != nil {
		return fmt.Errorf("error encoding bool array size %d: %w", len(v), err)
	}
	b := make([]byte, 0, len(v))
	for _, x := range v {
		if x {
			b = append(b, byte(1))
		} else {
			b = append(b, byte(0))
		}
	}
	if _, err := e.writer.Write(b); err != nil {
		return fmt.Errorf("error encoding bool array value %#v: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeSliceInt(v []int) error {

	if max := e.max(v); max < 256 {
		if err := e.writeByte(e.typeArray | e.typeInt8); err != nil {
			return fmt.Errorf("error encoding int8 array value %v: %w", v, err)
		}
		if err := e.encodeInt(len(v)); err != nil {
			return fmt.Errorf("error encoding int8 array size %d: %w", len(v), err)
		}
		b := make([]byte, 0, len(v))
		for _, x := range v {
			b = append(b, byte(x))
		}
		if _, err := e.writer.Write(b); err != nil {
			return fmt.Errorf("error encoding int8 array value %d: %w", v, err)
		}
		return nil
	}

	if err := e.writeByte(e.typeArray | e.typeInt64); err != nil {
		return fmt.Errorf("error encoding int64 array value %v: %w", v, err)
	}
	if err := e.encodeInt(len(v)); err != nil {
		return fmt.Errorf("error encoding int64 array size %d: %w", len(v), err)
	}
	b := make([]byte, 8*len(v))
	for i, x := range v {
		e.order.PutUint64(b[8*i:], uint64(x))
	}
	if _, err := e.writer.Write(b); err != nil {
		return fmt.Errorf("error encoding int64 array value %d: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeSliceFloat64(v []float64) error {
	if err := e.writeByte(e.typeArray | e.typeFloat64); err != nil {
		return fmt.Errorf("error encoding float64 array value %v: %w", v, err)
	}
	if err := e.encodeInt(len(v)); err != nil {
		return fmt.Errorf("error encoding float64 array size %d: %w", len(v), err)
	}
	b := make([]byte, 8*len(v))
	for i, x := range v {
		e.order.PutUint64(b[8*i:], math.Float64bits(x))
	}
	if _, err := e.writer.Write(b); err != nil {
		return fmt.Errorf("error encoding float64 array value %#v: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeSliceString(v []string) error {
	if err := e.writeByte(e.typeArray | e.typeString); err != nil {
		return fmt.Errorf("error encoding string value %v: %w", v, err)
	}
	if err := e.encodeInt(len(v)); err != nil {
		return fmt.Errorf("error encoding array size %d: %w", len(v), err)
	}
	if err := e.writeStrings(v); err != nil {
		return fmt.Errorf("error writing strings %q: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeMapStringString(m map[string]string) error {
	if err := e.writeBytes(e.typeMap|e.typeString, e.typeString); err != nil {
		return fmt.Errorf("error encoding string value %#v: %w", m, err)
	}

	if err := e.encodeInt(len(m)); err != nil {
		return fmt.Errorf("error encoding map size %d: %w", len(m), err)
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	values := make([]string, 0, len(keys))
	for _, k := range keys {
		values = append(values, m[k])
	}

	if err := e.writeStrings(keys); err != nil {
		return fmt.Errorf("error writing string keys for %#v: %w", m, err)
	}

	if err := e.writeStrings(values); err != nil {
		return fmt.Errorf("error writing string values %#v: %w", m, err)
	}

	return nil
}

func (e *LiteralEncoder) encodeMapStringInt(m map[string]int) error {

	if len(m) == 0 {
		if err := e.writeBytes(e.typeMap|e.typeString, e.typeInt8); err != nil {
			return fmt.Errorf("error encoding string value %v: %w", m, err)
		}
		if err := e.encodeInt(0); err != nil {
			return fmt.Errorf("error encoding map size %d: %w", len(m), err)
		}
		return nil
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	max := m[keys[0]]
	for i := 1; i < len(keys); i++ {
		if x := m[keys[i]]; x > max {
			max = x
		}
	}

	if max < 256 {
		values := make([]byte, 0, len(keys))
		for i := 1; i < len(keys); i++ {
			values = append(values, byte(m[keys[i]]))
		}
		if err := e.writeBytes(e.typeMap|e.typeString, e.typeInt8); err != nil {
			return fmt.Errorf("error encoding string value %v: %w", m, err)
		}
		if err := e.encodeInt(len(m)); err != nil {
			return fmt.Errorf("error encoding map size %d: %w", len(m), err)
		}
		if err := e.writeStrings(keys); err != nil {
			return fmt.Errorf("error writing strings %q: %w", keys, err)
		}
		if _, err := e.writer.Write(values); err != nil {
			return fmt.Errorf("error writing int8 array value %#v: %w", values, err)
		}
		return nil
	}

	values := make([]byte, 8*len(keys))
	for i := 1; i < len(keys); i++ {
		e.order.PutUint64(values[8*i:], uint64(m[keys[i]]))
	}
	if err := e.writeBytes(e.typeMap|e.typeString, e.typeInt64); err != nil {
		return fmt.Errorf("error encoding string value %v: %w", m, err)
	}
	if err := e.encodeInt(len(m)); err != nil {
		return fmt.Errorf("error encoding map size %d: %w", len(m), err)
	}
	if err := e.writeStrings(keys); err != nil {
		return fmt.Errorf("error writing strings %q: %w", m, err)
	}
	if _, err := e.writer.Write(values); err != nil {
		return fmt.Errorf("error encoding int8 array value %#v: %w", values, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeMapStringInterface(m map[string]interface{}) error {
	if len(m) == 0 {
		if err := e.writeBytes(e.typeMap|e.typeString, e.typeInterface); err != nil {
			return fmt.Errorf("error encoding map[string]interface{} value %v: %w", m, err)
		}
		if err := e.encodeInt(0); err != nil {
			return fmt.Errorf("error encoding map size %d: %w", len(m), err)
		}
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if err := e.writeBytes(e.typeMap|e.typeString, e.typeInterface); err != nil {
		return fmt.Errorf("error encoding header for map[string]interface value %#v: %w", m, err)
	}
	if err := e.encodeInt(len(m)); err != nil {
		return fmt.Errorf("error encoding map size %d: %w", len(m), err)
	}
	if err := e.writeStrings(keys); err != nil {
		return fmt.Errorf("error writing strings %q: %w", m, err)
	}
	for i := 0; i < len(keys); i++ {
		if err := e.Encode(m[keys[i]]); err != nil {
			return fmt.Errorf("error encoding key %q in map %#v: %w", keys[i], m, err)
		}
	}
	return nil
}

func (e *LiteralEncoder) encodeMapStringMapStringInterface(m map[string]map[string]interface{}) error {
	if len(m) == 0 {
		if err := e.writeBytes(e.typeMap|e.typeString, e.typeInterface); err != nil {
			return fmt.Errorf("error encoding map[string]interface{} value %v: %w", m, err)
		}
		if err := e.encodeInt(0); err != nil {
			return fmt.Errorf("error encoding map size %d: %w", len(m), err)
		}
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if err := e.writeBytes(e.typeMap|e.typeString, e.typeInterface); err != nil {
		return fmt.Errorf("error encoding header for map[string]interface value %#v: %w", m, err)
	}
	if err := e.encodeInt(len(m)); err != nil {
		return fmt.Errorf("error encoding map size %d: %w", len(m), err)
	}
	if err := e.writeStrings(keys); err != nil {
		return fmt.Errorf("error writing strings %q: %w", m, err)
	}
	for i := 0; i < len(keys); i++ {
		if err := e.Encode(m[keys[i]]); err != nil {
			return fmt.Errorf("error encoding key %q in map %#v: %w", keys[i], m, err)
		}
	}
	return nil
}

func (e *LiteralEncoder) encodeMapStringSliceFloat64(m map[string][]float64) error {
	if len(m) == 0 {
		if err := e.writeBytes(e.typeMap|e.typeString, e.typeInterface); err != nil {
			return fmt.Errorf("error encoding map[string]interface{} value %v: %w", m, err)
		}
		if err := e.encodeInt(0); err != nil {
			return fmt.Errorf("error encoding map size %d: %w", len(m), err)
		}
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if err := e.writeBytes(e.typeMap|e.typeString, e.typeInterface); err != nil {
		return fmt.Errorf("error encoding header for map[string]interface value %#v: %w", m, err)
	}
	if err := e.encodeInt(len(m)); err != nil {
		return fmt.Errorf("error encoding map size %d: %w", len(m), err)
	}
	if err := e.writeStrings(keys); err != nil {
		return fmt.Errorf("error writing strings %q: %w", m, err)
	}
	for i := 0; i < len(keys); i++ {
		if err := e.Encode(m[keys[i]]); err != nil {
			return fmt.Errorf("error encoding key %q in map %#v: %w", keys[i], m, err)
		}
	}
	return nil
}

func (e *LiteralEncoder) Encode(v interface{}) error {
	if b, ok := v.(bool); ok {
		return e.encodeBool(b)
	}
	if i, ok := v.(int); ok {
		return e.encodeInt(i)
	}
	if f, ok := v.(float64); ok {
		return e.encodeFloat64(f)
	}
	if str, ok := v.(string); ok {
		return e.encodeString(str)
	}
	if slc, ok := v.([]string); ok {
		return e.encodeSliceString(slc)
	}
	if slc, ok := v.([]bool); ok {
		return e.encodeSliceBool(slc)
	}
	if slc, ok := v.([]int); ok {
		return e.encodeSliceInt(slc)
	}
	if slc, ok := v.([]float64); ok {
		return e.encodeSliceFloat64(slc)
	}
	if slc, ok := v.([]interface{}); ok {
		return e.encodeSliceInterface(slc)
	}
	if m, ok := v.(map[string]string); ok {
		return e.encodeMapStringString(m)
	}
	if m, ok := v.(map[string]int); ok {
		return e.encodeMapStringInt(m)
	}
	if m, ok := v.(map[string]interface{}); ok {
		return e.encodeMapStringInterface(m)
	}
	if m, ok := v.(map[string]map[string]interface{}); ok {
		return e.encodeMapStringMapStringInterface(m)
	}
	if m, ok := v.(map[string][]float64); ok {
		return e.encodeMapStringSliceFloat64(m)
	}
	return fmt.Errorf("error encoding value %#v", v)
}

// Flush flushes the underlying writer, if it has a Flush method.
// This writer itself does no buffering.
func (e *LiteralEncoder) Flush() error {
	if flusher, ok := e.writer.(interface{ Flush() error }); ok {
		err := flusher.Flush()
		if err != nil {
			return fmt.Errorf("error flushing underlying writer: %w", err)
		}
	}
	return nil
}

// Close closes the underlying writer, if it has a Close method.
func (e *LiteralEncoder) Close() error {
	if closer, ok := e.writer.(interface{ Close() error }); ok {
		err := closer.Close()
		if err != nil {
			return fmt.Errorf("error closing underlying writer: %w", err)
		}
	}
	return nil
}
