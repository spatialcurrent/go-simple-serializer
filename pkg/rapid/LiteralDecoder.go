package rapid

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"
)

type LiteralDecoder struct {
	reader        io.ByteReader
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

func NewLiteralDecoder(r io.ByteReader) *LiteralDecoder {
	//fmt.Fprintf(os.Stderr, "string: % x\n", DefaultTypeString)
	//fmt.Fprintf(os.Stderr, "map[string]string: % x\n", []byte{DefaultTypeMap|DefaultTypeString, DefaultTypeString})
	//fmt.Fprintf(os.Stderr, "map[string]interface{}: % x\n", []byte{DefaultTypeMap|DefaultTypeString, DefaultTypeInterface})
	return &LiteralDecoder{
		reader:        r,
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

func (d *LiteralDecoder) Reset(r io.ByteReader) {
	d.reader = r
}

func (d *LiteralDecoder) readByte() (byte, error) {
	return d.reader.ReadByte()
}

func (d *LiteralDecoder) readBytes(n int) ([]byte, error) {
	b := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		d, err := d.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return b, io.EOF
			}
			return b, fmt.Errorf("error reading byte %d: %w", i, err)
		}
		b = append(b, d)
	}
	return b, nil
}

func (d *LiteralDecoder) readInt64() (int64, error) {
	b, err := d.readBytes(8)
	if err != nil {
		return int64(0), fmt.Errorf("error reading %d bytes: %w", 8, err)
	}
	return int64(d.order.Uint64(b)), nil
}

func (d *LiteralDecoder) readFloat64() (float64, error) {
	b, err := d.readBytes(8)
	if err != nil {
		return 0.0, fmt.Errorf("error reading float64: %w", err)
	}
	return math.Float64frombits(d.order.Uint64(b)), nil
}

func (d *LiteralDecoder) readBool() (bool, error) {
	x, err := d.reader.ReadByte()
	if err != nil {
		return false, fmt.Errorf("error reading bool: %w", err)
	}
	return x&0x1 == 0x1, nil
}

func (d *LiteralDecoder) readString() (string, error) {
	b := make([]byte, 0)
	for {
		d, err := d.reader.ReadByte()
		if err != nil {
			return "", fmt.Errorf("error reading string byte: %w", err)
		}
		if d == byte(0) {
			break
		}
		b = append(b, d)
	}
	return strings.ReplaceAll(string(b), "\\n", "\n"), nil
}

func (d *LiteralDecoder) readStrings(n int) ([]string, error) {
	slc := make([]string, 0, n)
	for i := 0; i < n; i++ {
		str, err := d.readString()
		if err != nil {
			return slc, fmt.Errorf("error reading element %d: %w", i, err)
		}
		slc = append(slc, str)
	}
	return slc, nil
}

func (d *LiteralDecoder) readFloat64s(n int) ([]float64, error) {
	slc := make([]float64, 0, n)
	for i := 0; i < n; i++ {
		f, err := d.readFloat64()
		if err != nil {
			return slc, fmt.Errorf("error reading floats: error reading float %d: %w", i, err)
		}
		slc = append(slc, f)
	}
	return slc, nil
}

func (d *LiteralDecoder) readMapStringString(n int) (map[string]string, error) {
	m := map[string]string{}
	keys, err := d.readStrings(n)
	if err != nil {
		return m, fmt.Errorf("error reading keys: %w", err)
	}
	values, err := d.readStrings(n)
	if err != nil {
		return m, fmt.Errorf("error reading values: %w", err)
	}
	for i := 0; i < n; i++ {
		m[keys[i]] = values[i]
	}
	return m, nil
}

func (d *LiteralDecoder) readMapStringInterface(n int) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	keys, err := d.readStrings(n)
	if err != nil {
		return m, fmt.Errorf("error reading keys: %w", err)
	}
	values, err := d.readInterfaces(n)
	if err != nil {
		return m, fmt.Errorf("error reading values: %w", err)
	}
	for i := 0; i < n; i++ {
		m[keys[i]] = values[i]
	}
	return m, nil
}

func (d *LiteralDecoder) readInterfaces(n int) ([]interface{}, error) {
	slc := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		xh, err := d.readByte()
		if err != nil {
			return slc, fmt.Errorf("error reading byte: %w", err)
		}
		if xh&d.typeMap == d.typeMap {
			if xh&d.typeString == d.typeString {
				xh2, err := d.readByte()
				if err != nil {
					return slc, fmt.Errorf("error reading byte: %w", err)
				}
				if xh2&d.typeString == d.typeString {
					size := 0
					err := d.decodeInt(&size)
					if err != nil {
						return slc, fmt.Errorf("error decoding size of map of string to string: %w", err)
					}
					m, err := d.readMapStringString(size)
					if err != nil {
						return slc, fmt.Errorf("error reading map of string to string: %w", err)
					}
					slc = append(slc, m)
					continue
				}
				if xh2&d.typeInterface == d.typeInterface {
					size := 0
					err := d.decodeInt(&size)
					if err != nil {
						return slc, fmt.Errorf("error decoding size of map string to interface: %w", err)
					}
					m, err := d.readMapStringInterface(size)
					if err != nil {
						return slc, fmt.Errorf("error reading map of string to interface: %w", err)
					}
					slc = append(slc, m)
					continue
				}
			}
			return slc, fmt.Errorf("unsupported type %x", xh)
		}
		if xh&d.typeArray == d.typeArray {
			if xh&d.typeString == d.typeString {
				size := 0
				err := d.decodeInt(&size)
				if err != nil {
					return slc, fmt.Errorf("error decoding size of map of string to string: %w", err)
				}
				values, err := d.readStrings(size)
				if err != nil {
					return slc, fmt.Errorf("error reading array of strings with size %d: %w", size, err)
				}
				slc = append(slc, values)
				continue
			}
			if xh&d.typeFloat64 == d.typeFloat64 {
				size := 0
				err := d.decodeInt(&size)
				if err != nil {
					return slc, fmt.Errorf("error decoding size of map of string to string: %w", err)
				}
				values, err := d.readFloat64s(size)
				if err != nil {
					return slc, fmt.Errorf("error reading array of floats with size %d: %w", size, err)
				}
				slc = append(slc, values)
				continue
			}
			if xh&d.typeInterface == d.typeInterface {
				size := 0
				err := d.decodeInt(&size)
				if err != nil {
					return slc, fmt.Errorf("error decoding size of map string to interface: %w", err)
				}
				values, err := d.readInterfaces(size)
				if err != nil {
					return slc, fmt.Errorf("error reading map of string to interface: %w", err)
				}
				slc = append(slc, values)
				continue
			}
			return slc, fmt.Errorf("unsupported type %x", xh)
		}
		if xh&d.typeBool == d.typeBool {
			xd, err := d.readBool()
			if err != nil {
				return slc, fmt.Errorf("error reading bool: %w", err)
			}
			slc = append(slc, xd)
			continue
		}
		if xh&d.typeInt8 == d.typeInt8 {
			xd, err := d.readByte()
			if err != nil {
				return slc, fmt.Errorf("error reading int8: %w", err)
			}
			slc = append(slc, int(int8(xd)))
			continue
		}
		if xh&d.typeInt64 == d.typeInt64 {
			xd, err := d.readInt64()
			if err != nil {
				return slc, fmt.Errorf("error reading int64: %w", err)
			}
			slc = append(slc, int(xd))
			continue
		}
		if xh&d.typeFloat64 == d.typeFloat64 {
			xd, err := d.readFloat64()
			if err != nil {
				return slc, fmt.Errorf("error reading float64: %w", err)
			}
			slc = append(slc, xd)
			continue
		}
		if xh&d.typeString == d.typeString {
			xs, err := d.readString()
			if err != nil {
				return slc, fmt.Errorf("error reading string: %w", err)
			}
			slc = append(slc, xs)
			continue
		}
		return slc, fmt.Errorf("unsupported type %x", xh)
	}
	return slc, nil
}

func (d *LiteralDecoder) decodeBool(v *bool) error {
	h, err := d.reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h != d.typeBool {
		return fmt.Errorf("unexpected type code, found %x, expecting bool type code %x", h, d.typeBool)
	}
	x, err := d.reader.ReadByte()
	if err != nil {
		return fmt.Errorf("error decoding bool: %w", err)
	}
	*v = (x & 0x1) == 0x1
	return nil
}

func (d *LiteralDecoder) decodeInt(v *int) error {
	h, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding int header: %w", err)
	}
	switch h {
	case d.typeInt8:
		d, err := d.readByte()
		if err != nil {
			return fmt.Errorf("error reading byte: %w", err)
		}
		*v = int(uint8(d))
		return nil
	case d.typeInt64:
		d, err := d.readInt64()
		if err != nil {
			return fmt.Errorf("error reading int64: %w", err)
		}
		*v = int(d)
		return nil
	}
	return fmt.Errorf("unexpected type code, found %d (0x%x), expecting int type code", h, h)
}

func (d *LiteralDecoder) decodeFloat64(v *float64) error {
	h, err := d.reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h != d.typeFloat64 {
		return fmt.Errorf("unexpected type code, found %x, expecting float64 type code %x", h, d.typeFloat64)
	}
	x, err := d.readFloat64()
	if err != nil {
		return fmt.Errorf("error reading float64: %w", err)
	}
	*v = x
	return nil
}

func (d *LiteralDecoder) decodeString(v *string) error {
	h, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h != d.typeString {
		return fmt.Errorf("unexpected type code, found %x, expecting string type code %x", h, d.typeString)
	}
	str, err := d.readString()
	if err != nil {
		return fmt.Errorf("error reading string: %w", err)
	}
	*v = str
	return nil
}

func (d *LiteralDecoder) decodeSliceString(v *[]string) error {
	h, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if !(((h & d.typeArray) == d.typeArray) && ((h & d.typeString) == d.typeString)) {
		return fmt.Errorf("unexpected type code, found %x, expecting string array type code %x", h, d.typeArray&d.typeString)
	}
	size := 0
	if err := d.decodeInt(&size); err != nil {
		return fmt.Errorf("error decoding size of slice of strings: %w", err)
	}
	slc, err := d.readStrings(size)
	if err != nil {
		return fmt.Errorf("error reading strings: %w", err)
	}
	*v = slc
	return nil
}

func (d *LiteralDecoder) decodeSliceBool(v *[]bool) error {
	h, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding slice bool: %w", err)
	}
	if h&d.typeArray != d.typeArray || h&d.typeBool != d.typeBool {
		return fmt.Errorf("unexpected type code, found %x, expecting bool array type code", h)
	}
	size := 0
	if err := d.decodeInt(&size); err != nil {
		return fmt.Errorf("error decoding size of slice of bools: %w", err)
	}
	slc := make([]bool, 0, size)
	for i := 0; i < size; i++ {
		d, err := d.readByte()
		if err != nil {
			return fmt.Errorf("error decoding bool: %w", err)
		}
		slc = append(slc, d&0x1 == 0x1)
	}
	*v = slc
	return nil
}

func (d *LiteralDecoder) decodeSliceInt(v *[]int) error {
	h, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h&d.typeArray != d.typeArray {
		return fmt.Errorf("unexpected type code, found %x, expecting int array type code", h)
	}
	switch h & d.typeInt8 & d.typeInt64 {
	case d.typeInt8:
		size := 0
		err := d.decodeInt(&size)
		if err != nil {
			return fmt.Errorf("error decoding size of slice of int8s: %w", err)
		}
		slc := make([]int, 0, size)
		for i := 0; i < size; i++ {
			d, err := d.readByte()
			if err != nil {
				return fmt.Errorf("error decoding bool: %w", err)
			}
			slc = append(slc, int(d))
		}
		*v = slc
		return nil
	case d.typeInt64:
		size := 0
		err := d.decodeInt(&size)
		if err != nil {
			return fmt.Errorf("error decoding size of slice of int64s: %w", err)
		}
		slc := make([]int, 0, size)
		for i := 0; i < size; i++ {
			x, err := d.readInt64()
			if err != nil {
				return fmt.Errorf("error reading int64: %w", err)
			}
			slc = append(slc, int(x))
		}
		*v = slc
		return nil
	}
	return fmt.Errorf("unexpected type code, found %x, expecting int array type code", h)
}

func (d *LiteralDecoder) decodeSliceFloat64(v *[]float64) error {
	h, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h&d.typeArray != d.typeArray || h&d.typeFloat64 != d.typeFloat64 {
		return fmt.Errorf("unexpected type code, found %x, expecting float64 array type code", h)
	}
	size := 0
	if err := d.decodeInt(&size); err != nil {
		return fmt.Errorf("error decoding size of slice of float64s: %w", err)
	}
	slc, err := d.readFloat64s(size)
	if err != nil {
		return fmt.Errorf("error reading map of string to string: %w", err)
	}
	*v = slc
	return nil
}

func (d *LiteralDecoder) decodeSliceInterface(v *[]interface{}) error {
	h, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h&d.typeArray != d.typeArray || h&d.typeInterface != d.typeInterface {
		return fmt.Errorf("unexpected type code, found %x, expecting interface array type code", h)
	}
	size := 0
	if err := d.decodeInt(&size); err != nil {
		return fmt.Errorf("error decoding size of slice of interfaces: %w", err)
	}
	slc, err := d.readInterfaces(size)
	if err != nil {
		return fmt.Errorf("error reading interfaces: %w", err)
	}
	*v = slc
	return nil
}

func (d *LiteralDecoder) decodeMapStringString(v *map[string]string) error {
	h, err := d.readBytes(2)
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h[0]&d.typeMap != d.typeMap || h[0]&d.typeString != d.typeString || h[1]&d.typeString != d.typeString {
		return fmt.Errorf("unexpected type code, found %x, expecting map string to string type codes", h)
	}
	size := 0
	if err := d.decodeInt(&size); err != nil {
		return fmt.Errorf("error decoding size of map of string to string: %w", err)
	}
	m, err := d.readMapStringString(size)
	if err != nil {
		return fmt.Errorf("error reading map of string to string: %w", err)
	}
	*v = m
	return nil
}

func (d *LiteralDecoder) decodeMapStringInt(v *map[string]int) error {
	h, err := d.readBytes(2)
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h[0]&d.typeMap != d.typeMap || h[0]&d.typeString != d.typeString {
		return fmt.Errorf("unexpected type code, found %x, expecting map string to string type codes", h)
	}
	if h[1]&d.typeInt8 == d.typeInt8 {
		size := 0
		err := d.decodeInt(&size)
		if err != nil {
			return fmt.Errorf("error decoding size of slice of strings: %w", err)
		}
		keys, err := d.readStrings(size)
		if err != nil {
			return fmt.Errorf("error reading keys: %w", err)
		}
		values := make([]int, 0, size)
		for i := 0; i < size; i++ {
			value, err := d.readByte()
			if err != nil {
				return fmt.Errorf("error reading element %d: %w", i, err)
			}
			values = append(values, int(uint8(value)))
		}
		m := map[string]int{}
		for i := 0; i < size; i++ {
			m[keys[i]] = values[i]
		}
		*v = m
		return nil
	}
	if h[1]&d.typeInt64 == d.typeInt64 {
		size := 0
		err := d.decodeInt(&size)
		if err != nil {
			return fmt.Errorf("error decoding size of slice of strings: %w", err)
		}
		keys, err := d.readStrings(size)
		if err != nil {
			return fmt.Errorf("error reading keys: %w", err)
		}
		values := make([]int, 0, size)
		for i := 0; i < size; i++ {
			value, err := d.readInt64()
			if err != nil {
				return fmt.Errorf("error reading element %d: %w", i, err)
			}
			values = append(values, int(value))
		}
		m := map[string]int{}
		for i := 0; i < size; i++ {
			m[keys[i]] = values[i]
		}
		*v = m
		return nil
	}
	return fmt.Errorf("unexpected type code, found %x, expecting map string to string type codes", h)
}

func (d *LiteralDecoder) decodeMapStringInterface(v *map[string]interface{}) error {
	h, err := d.readBytes(2)
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding map header: %w", err)
	}
	if h[0]&d.typeMap != d.typeMap || h[0]&d.typeString != d.typeString || h[1]&d.typeInterface != d.typeInterface {
		return fmt.Errorf("unexpected type code, found %x, expecting map string to interface type codes", h)
	}
	size := 0
	if err := d.decodeInt(&size); err != nil {
		return fmt.Errorf("error decoding size of map of string to interface: %w", err)
	}
	keys, err := d.readStrings(size)
	if err != nil {
		return fmt.Errorf("error reading keys: %w", err)
	}
	values, err := d.readInterfaces(size)
	if err != nil {
		return fmt.Errorf("error reading interfaces: %w", err)
	}
	m := map[string]interface{}{}
	for i := 0; i < size; i++ {
		m[keys[i]] = values[i]
	}
	*v = m
	return nil
}

func (d *LiteralDecoder) Decode(v interface{}) error {
	if b, ok := v.(*bool); ok {
		return d.decodeBool(b)
	}
	if i, ok := v.(*int); ok {
		return d.decodeInt(i)
	}
	if f, ok := v.(*float64); ok {
		return d.decodeFloat64(f)
	}
	if str, ok := v.(*string); ok {
		return d.decodeString(str)
	}
	if slc, ok := v.(*[]string); ok {
		return d.decodeSliceString(slc)
	}
	if slc, ok := v.(*[]bool); ok {
		return d.decodeSliceBool(slc)
	}
	if slc, ok := v.(*[]int); ok {
		return d.decodeSliceInt(slc)
	}
	if slc, ok := v.(*[]float64); ok {
		return d.decodeSliceFloat64(slc)
	}
	if slc, ok := v.(*[]interface{}); ok {
		return d.decodeSliceInterface(slc)
	}
	if m, ok := v.(*map[string]string); ok {
		return d.decodeMapStringString(m)
	}
	if m, ok := v.(*map[string]int); ok {
		return d.decodeMapStringInt(m)
	}
	if m, ok := v.(*map[string]interface{}); ok {
		return d.decodeMapStringInterface(m)
	}
	return fmt.Errorf("error encoding value %#v", v)
}

// Close closes the underlying writer, if it has a Close method.
func (e *LiteralDecoder) Close() error {
	if closer, ok := e.reader.(interface{ Close() error }); ok {
		err := closer.Close()
		if err != nil {
			return fmt.Errorf("error closing underlying writer: %w", err)
		}
	}
	return nil
}
