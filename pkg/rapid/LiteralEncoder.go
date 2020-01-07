package rapid

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strings"
)

type LiteralEncoder struct {
	typeInterface          byte
	typeBool               byte
	typeUInt2              byte
	typeUInt3              byte
	typeUInt4              byte
	typeUInt8              byte
	typeUInt16             byte
	typeInt32              byte
	typeInt64              byte
	typeFloat32            byte
	typeFloat64            byte
	typeString             byte
	typeArrayInterface     byte
	typeArrayBool          byte
	typeArrayUInt2         byte
	typeArrayUInt4         byte
	typeArrayUInt8         byte
	typeArrayUInt16        byte
	typeArrayInt32         byte
	typeArrayInt64         byte
	typeArrayFloat32       byte
	typeArrayFloat32x2     byte
	typeArrayFloat64       byte
	typeArrayFloat64x2     byte
	typeArrayString        byte
	typeMapString          byte
	typeMapStringString    byte
	typeMapStringInterface byte
	writer                 io.Writer
	order                  binary.ByteOrder
	sortMapKeys            bool
	convertFloatToInt      bool
	size                   int
}

func NewLiteralEncoder(w io.Writer) *LiteralEncoder {
	return &LiteralEncoder{
		typeInterface:          DefaultTypeInterface,
		typeBool:               DefaultTypeBool,
		typeUInt2:              DefaultTypeUInt2,
		typeUInt3:              DefaultTypeUInt3,
		typeUInt4:              DefaultTypeUInt4,
		typeUInt8:              DefaultTypeUInt8,
		typeUInt16:             DefaultTypeUInt16,
		typeInt32:              DefaultTypeInt32,
		typeInt64:              DefaultTypeInt64,
		typeFloat32:            DefaultTypeFloat32,
		typeFloat64:            DefaultTypeFloat64,
		typeString:             DefaultTypeString,
		typeArrayInterface:     DefaultTypeArrayInterface,
		typeArrayBool:          DefaultTypeArrayBool,
		typeArrayUInt2:         DefaultTypeArrayUInt2,
		typeArrayUInt4:         DefaultTypeArrayUInt4,
		typeArrayUInt8:         DefaultTypeArrayUInt8,
		typeArrayUInt16:        DefaultTypeArrayUInt16,
		typeArrayInt32:         DefaultTypeArrayInt32,
		typeArrayInt64:         DefaultTypeArrayInt64,
		typeArrayFloat32:       DefaultTypeArrayFloat32,
		typeArrayFloat32x2:     DefaultTypeArrayFloat32x2,
		typeArrayFloat64:       DefaultTypeArrayFloat64,
		typeArrayFloat64x2:     DefaultTypeArrayFloat64x2,
		typeArrayString:        DefaultTypeArrayString,
		typeMapString:          DefaultTypeMapString,
		typeMapStringString:    DefaultTypeMapStringString,
		typeMapStringInterface: DefaultTypeMapStringInterface,
		//
		writer: w,
		//
		order: binary.LittleEndian,
		//
		sortMapKeys:       true,
		convertFloatToInt: true,
		//
		size: 0,
	}
}

func (e *LiteralEncoder) isInt(f float64) bool {
	return f == math.Trunc(f)
}

func (e *LiteralEncoder) sortStrings(keys []string) {
	//sort.Strings(keys)

	sort.Slice(keys, func(i, j int) bool { return len(keys[i]) < len(keys[j]) })
}

func (e *LiteralEncoder) Size() int {
	return e.size
}

func (e *LiteralEncoder) Reset(w io.Writer) {
	e.writer = w
	e.size = 0
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
	n, err := e.writer.Write([]byte{v})
	if err != nil {
		return fmt.Errorf("error encoding byte value %v: %w", v, err)
	}
	e.size += n
	return nil
}

func (e *LiteralEncoder) writeBytes(v ...byte) error {
	n, err := e.writer.Write(v)
	if err != nil {
		return fmt.Errorf("error encoding bytes %x: %w", v, err)
	}
	e.size += n
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
		if err := e.writeBytes(b...); err != nil {
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
		err := e.writeBytes(e.typeBool | MaskUpper1)
		if err != nil {
			return fmt.Errorf("error encoding int2 value %d: %w", v, err)
		}
		return nil
	}
	err := e.writeBytes(e.typeBool)
	if err != nil {
		return fmt.Errorf("error encoding int2 value %d: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeUInt2(v int) error {
	err := e.writeBytes(e.typeUInt2 | (byte(v) << 6))
	if err != nil {
		return fmt.Errorf("error encoding int2 value %d: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeUInt3(v int) error {
	//fmt.Fprintf(os.Stderr, "Encoding uint3: %#v %b\n", v, (byte(v) << 5))
	err := e.writeBytes(e.typeUInt3 | (byte(v) << 5))
	if err != nil {
		return fmt.Errorf("error encoding int3 value %d: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeUInt8(v uint8) error {
	err := e.writeBytes(e.typeUInt8, byte(v))
	if err != nil {
		return fmt.Errorf("error encoding uint64 value %d: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeUInt16(v uint16) error {
	b := make([]byte, 1+2)
	b[0] = e.typeUInt16
	e.order.PutUint16(b[1:], uint16(v))
	err := e.writeBytes(b...)
	if err != nil {
		return fmt.Errorf("error encoding int32 value %d: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeInt32(v int32) error {
	b := make([]byte, 1+4)
	b[0] = e.typeInt32
	e.order.PutUint32(b[1:], uint32(v))
	err := e.writeBytes(b...)
	if err != nil {
		return fmt.Errorf("error encoding int32 value %d: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeInt64(v int64) error {
	b := make([]byte, 1+8)
	b[0] = e.typeInt64
	e.order.PutUint64(b[1:], uint64(v))
	err := e.writeBytes(b...)
	if err != nil {
		return fmt.Errorf("error encoding uint64 value %d: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeInt(v int) error {
	if v >= 0 && v < 4 {
		return e.encodeUInt2(v)
	}
	if v >= 0 && v < 8 {
		return e.encodeUInt3(v)
	}
	if v >= 0 && v < 256 {
		return e.encodeUInt8(uint8(v))
	}
	if v >= 0 && v < 65536 {
		return e.encodeUInt16(uint16(v))
	}
	if v >= -2147483648 && v < 2147483648 {
		return e.encodeInt32(int32(v))
	}
	return e.encodeInt64(int64(v))
}

func (e *LiteralEncoder) encodeFloat32(v float32) error {
	fmt.Fprintf(os.Stderr, "encodeFloat32(%#v)\n", v)
	b := make([]byte, 1+4)
	b[0] = e.typeFloat32
	e.order.PutUint32(b[1:], math.Float32bits(v))
	err := e.writeBytes(b...)
	if err != nil {
		return fmt.Errorf("error encoding float32 value %v: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeFloat64(v float64) error {
	if e.convertFloatToInt && e.isInt(v) {
		if v >= 0 && v < MaxExactInt64 {
			if v < MaxExactInt32 {
				i := int(int32(v))
				if i < 4 {
					return e.encodeUInt2(i)
				}
				if i < 8 {
					return e.encodeUInt3(i)
				}
				if i < 256 {
					return e.encodeUInt8(uint8(i))
				}
				if i < 65536 {
					return e.encodeUInt16(uint16(i))
				}
				//fmt.Fprintf(os.Stderr, "convert float to int32: %d\n", int32(v))
				return e.encodeInt32(int32(i))
			}
			//fmt.Fprintf(os.Stderr, "convert float to int64: %d\n", int64(v))
			return e.encodeInt64(int64(v))
		}
	}
	b := make([]byte, 1+8)
	b[0] = e.typeFloat64
	e.order.PutUint64(b[1:], math.Float64bits(v))
	err := e.writeBytes(b...)
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
	b[len(b)-1] = NullByte
	//
	if err := e.writeBytes(b...); err != nil {
		return fmt.Errorf("error encoding string value %v: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeSliceInterface(v []interface{}) error {
	if err := e.writeByte(e.typeArrayInterface); err != nil {
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
	if err := e.writeByte(e.typeArrayBool); err != nil {
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
	if err := e.writeBytes(b...); err != nil {
		return fmt.Errorf("error encoding bool array value %#v: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeSliceInt(v []int) error {

	if len(v) == 0 {
		if err := e.writeByte(e.typeArrayInt64); err != nil {
			return fmt.Errorf("error encoding int64 array value %v: %w", v, err)
		}
		if err := e.encodeInt(0); err != nil {
			return fmt.Errorf("error encoding int64 array size %d: %w", len(v), err)
		}
	}

	min := v[0]
	max := v[0]
	for i := 1; i < len(v); i++ {
		if x := v[i]; x > max {
			max = x
		}
		if x := v[i]; x < min {
			min = x
		}
	}

	if min >= 0 && max < 4 {
		if err := e.writeByte(e.typeArrayUInt2); err != nil {
			return fmt.Errorf("error encoding int8 array value %v: %w", v, err)
		}
		if err := e.encodeInt(len(v)); err != nil {
			return fmt.Errorf("error encoding int8 array size %d: %w", len(v), err)
		}
		b := make([]byte, 0, len(v)/4)
		for i := 0; i < len(v)/4; i++ {
			b = append(b, byte(v[i*4])|(byte(v[(i*4)+1])<<2)|(byte(v[(i*4)+2])<<4)|(byte(v[(i*4)+3])<<6))
		}
		if err := e.writeBytes(b...); err != nil {
			return fmt.Errorf("error encoding int8 array value %d: %w", v, err)
		}
		return nil
	}

	if min >= 0 && max < 16 {
		if err := e.writeByte(e.typeArrayUInt4); err != nil {
			return fmt.Errorf("error encoding int8 array value %v: %w", v, err)
		}
		if err := e.encodeInt(len(v)); err != nil {
			return fmt.Errorf("error encoding int8 array size %d: %w", len(v), err)
		}
		b := make([]byte, 0, len(v)/2)
		for i := 0; i < len(v)/2; i++ {
			b = append(b, byte(v[i*2])|(byte(v[(i*2)+1])>>4))
		}
		if err := e.writeBytes(b...); err != nil {
			return fmt.Errorf("error encoding int8 array value %d: %w", v, err)
		}
		return nil
	}

	if min >= 0 && max < 256 {
		if err := e.writeByte(e.typeArrayUInt8); err != nil {
			return fmt.Errorf("error encoding int8 array value %v: %w", v, err)
		}
		if err := e.encodeInt(len(v)); err != nil {
			return fmt.Errorf("error encoding int8 array size %d: %w", len(v), err)
		}
		b := make([]byte, 0, len(v))
		for _, x := range v {
			b = append(b, byte(x))
		}
		if err := e.writeBytes(b...); err != nil {
			return fmt.Errorf("error encoding int8 array value %d: %w", v, err)
		}
		return nil
	}

	if err := e.writeByte(e.typeArrayInt64); err != nil {
		return fmt.Errorf("error encoding int64 array value %v: %w", v, err)
	}
	if err := e.encodeInt(len(v)); err != nil {
		return fmt.Errorf("error encoding int64 array size %d: %w", len(v), err)
	}
	b := make([]byte, 8*len(v))
	for i, x := range v {
		e.order.PutUint64(b[8*i:], uint64(x))
	}
	if err := e.writeBytes(b...); err != nil {
		return fmt.Errorf("error encoding int64 array value %d: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeSliceFloat32(v []float32) error {
	size := len(v)
	if size == 2 {
		if err := e.writeByte(e.typeArrayFloat32x2); err != nil {
			return fmt.Errorf("error encoding float32 array value %v: %w", v, err)
		}
		b := make([]byte, 8)
		e.order.PutUint32(b[0:], math.Float32bits(v[0]))
		e.order.PutUint32(b[4:], math.Float32bits(v[1]))
		if err := e.writeBytes(b...); err != nil {
			return fmt.Errorf("error encoding float32 array value %#v: %w", v, err)
		}
		return nil
	}
	if err := e.writeByte(e.typeArrayFloat32); err != nil {
		return fmt.Errorf("error encoding float32 array value %v: %w", v, err)
	}
	if err := e.encodeInt(len(v)); err != nil {
		return fmt.Errorf("error encoding float32 array size %d: %w", len(v), err)
	}
	b := make([]byte, 4*len(v))
	for i, x := range v {
		e.order.PutUint32(b[4*i:], math.Float32bits(x))
	}
	if err := e.writeBytes(b...); err != nil {
		return fmt.Errorf("error encoding float32 array value %#v: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeSliceFloat64(v []float64) error {
	size := len(v)
	if size == 2 {
		if err := e.writeByte(e.typeArrayFloat64x2); err != nil {
			return fmt.Errorf("error encoding float64 array value %v: %w", v, err)
		}
		b := make([]byte, 16)
		e.order.PutUint64(b[0:], math.Float64bits(v[0]))
		e.order.PutUint64(b[8:], math.Float64bits(v[1]))
		if err := e.writeBytes(b...); err != nil {
			return fmt.Errorf("error encoding float64 array value %#v: %w", v, err)
		}
		return nil
	}
	if err := e.writeByte(e.typeArrayFloat64); err != nil {
		return fmt.Errorf("error encoding float64 array value %v: %w", v, err)
	}
	if err := e.encodeInt(len(v)); err != nil {
		return fmt.Errorf("error encoding float64 array size %d: %w", len(v), err)
	}
	b := make([]byte, 8*len(v))
	for i, x := range v {
		e.order.PutUint64(b[8*i:], math.Float64bits(x))
	}
	if err := e.writeBytes(b...); err != nil {
		return fmt.Errorf("error encoding float64 array value %#v: %w", v, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeSliceString(v []string) error {
	if err := e.writeByte(e.typeArrayString); err != nil {
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
	if err := e.writeBytes(e.typeMapStringString); err != nil {
		return fmt.Errorf("error encoding string value %#v: %w", m, err)
	}

	if err := e.encodeInt(len(m)); err != nil {
		return fmt.Errorf("error encoding map size %d: %w", len(m), err)
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	if e.sortMapKeys {
		e.sortStrings(keys)
	}

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
		if err := e.writeBytes(e.typeMapString, e.typeUInt8); err != nil {
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

	if e.sortMapKeys {
		e.sortStrings(keys)
	}

	min := m[keys[0]]
	max := m[keys[0]]
	for i := 1; i < len(keys); i++ {
		if x := m[keys[i]]; x > max {
			max = x
		}
		if x := m[keys[i]]; x < min {
			min = x
		}
	}

	if min >= 0 && max < 256 {
		values := make([]byte, 0, len(keys))
		for i := 1; i < len(keys); i++ {
			values = append(values, byte(m[keys[i]]))
		}
		if err := e.writeBytes(e.typeMapString, e.typeUInt8); err != nil {
			return fmt.Errorf("error encoding string value %v: %w", m, err)
		}
		if err := e.encodeInt(len(m)); err != nil {
			return fmt.Errorf("error encoding map size %d: %w", len(m), err)
		}
		if err := e.writeStrings(keys); err != nil {
			return fmt.Errorf("error writing strings %q: %w", keys, err)
		}
		if err := e.writeBytes(values...); err != nil {
			return fmt.Errorf("error writing int8 array value %#v: %w", values, err)
		}
		return nil
	}

	values := make([]byte, 8*len(keys))
	for i := 1; i < len(keys); i++ {
		e.order.PutUint64(values[8*i:], uint64(m[keys[i]]))
	}
	if err := e.writeBytes(e.typeMapString, e.typeInt64); err != nil {
		return fmt.Errorf("error encoding string value %v: %w", m, err)
	}
	if err := e.encodeInt(len(m)); err != nil {
		return fmt.Errorf("error encoding map size %d: %w", len(m), err)
	}
	if err := e.writeStrings(keys); err != nil {
		return fmt.Errorf("error writing strings %q: %w", m, err)
	}
	if err := e.writeBytes(values...); err != nil {
		return fmt.Errorf("error encoding int8 array value %#v: %w", values, err)
	}
	return nil
}

func (e *LiteralEncoder) encodeMapStringInterface(m map[string]interface{}) error {
	if len(m) == 0 {
		if err := e.writeBytes(e.typeMapStringInterface); err != nil {
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

	if e.sortMapKeys {
		e.sortStrings(keys)
	}

	if err := e.writeBytes(e.typeMapStringInterface); err != nil {
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
		if err := e.writeBytes(e.typeMapStringInterface); err != nil {
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
	if e.sortMapKeys {
		e.sortStrings(keys)
	}
	if err := e.writeBytes(e.typeMapStringInterface); err != nil {
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
		if err := e.writeBytes(e.typeMapString, e.typeInterface); err != nil {
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
	if e.sortMapKeys {
		e.sortStrings(keys)
	}
	if err := e.writeBytes(e.typeMapString, e.typeInterface); err != nil {
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
	if f, ok := v.(float32); ok {
		return e.encodeFloat32(f)
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
	if slc, ok := v.([]float32); ok {
		return e.encodeSliceFloat32(slc)
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
