package rapid

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
)

type LiteralDecoder struct {
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
	//
	reader io.ByteReader
	//
	order binary.ByteOrder
}

func NewLiteralDecoder(r io.ByteReader) *LiteralDecoder {
	return &LiteralDecoder{
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
		reader: r,
		//
		order: binary.LittleEndian,
	}
}

func (d *LiteralDecoder) Reset(r io.ByteReader) {
	d.reader = r
}

func (d *LiteralDecoder) isInt(b byte) bool {
	return ((b == d.typeUInt2) ||
		(b == d.typeUInt3) ||
		(b == d.typeUInt4) ||
		(b == d.typeUInt8) ||
		(b == d.typeUInt16) ||
		(b == d.typeInt32) ||
		(b == d.typeInt64))
}

func (d *LiteralDecoder) isIntArray(b byte) bool {
	return ((b == d.typeArrayUInt2) ||
		(b == d.typeArrayUInt4) ||
		(b == d.typeArrayUInt8) ||
		(b == d.typeArrayUInt16) ||
		(b == d.typeArrayInt32) ||
		(b == d.typeArrayInt64))
}

func (d *LiteralDecoder) isArray(b byte) bool {
	return d.isIntArray(b) || ((b == d.typeArrayInterface) ||
		(b == d.typeArrayBool) ||
		(b == d.typeArrayFloat64) ||
		(b == d.typeArrayFloat64x2) ||
		(b == d.typeArrayString))
}

func (d *LiteralDecoder) isMap(b byte) bool {
	return (b == d.typeMapString) || (b == d.typeMapStringString) || (b == d.typeMapStringInterface)
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

func (d *LiteralDecoder) readUInt16() (uint16, error) {
	b, err := d.readBytes(2)
	if err != nil {
		return uint16(0), fmt.Errorf("error reading uint16s: %w", err)
	}
	return d.order.Uint16(b), nil
}

func (d *LiteralDecoder) readInt32() (int32, error) {
	b, err := d.readBytes(4)
	if err != nil {
		return int32(0), fmt.Errorf("error reading int32: %w", err)
	}
	return int32(d.order.Uint32(b)), nil
}

func (d *LiteralDecoder) readInt64() (int64, error) {
	b, err := d.readBytes(8)
	if err != nil {
		return int64(0), fmt.Errorf("error reading int64: %w", err)
	}
	return int64(d.order.Uint64(b)), nil
}

func (d *LiteralDecoder) readFloat32() (float32, error) {
	b, err := d.readBytes(4)
	if err != nil {
		return 0.0, fmt.Errorf("error reading float32: %w", err)
	}
	return math.Float32frombits(d.order.Uint32(b)), nil
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

func (d *LiteralDecoder) readFloat32s(n int) ([]float32, error) {
	slc := make([]float32, 0, n)
	for i := 0; i < n; i++ {
		f, err := d.readFloat32()
		if err != nil {
			return slc, fmt.Errorf("error reading float32s: error reading float32 %d: %w", i, err)
		}
		slc = append(slc, f)
	}
	return slc, nil
}

func (d *LiteralDecoder) readFloat64s(n int) ([]float64, error) {
	slc := make([]float64, 0, n)
	for i := 0; i < n; i++ {
		f, err := d.readFloat64()
		if err != nil {
			return slc, fmt.Errorf("error reading float64s: error reading float64 %d: %w", i, err)
		}
		slc = append(slc, f)
	}
	return slc, nil
}

func (d *LiteralDecoder) readMapStringString(n int) (map[string]string, error) {
	keys, err := d.readStrings(n)
	if err != nil {
		return map[string]string{}, fmt.Errorf("error reading keys: %w", err)
	}
	values, err := d.readStrings(n)
	if err != nil {
		return map[string]string{}, fmt.Errorf("error reading values: %w", err)
	}
	m := make(map[string]string, len(keys))
	for i := 0; i < n; i++ {
		m[keys[i]] = values[i]
	}
	return m, nil
}

func (d *LiteralDecoder) readMapStringInterface(n int) (map[string]interface{}, error) {
	keys, err := d.readStrings(n)
	if err != nil {
		return map[string]interface{}{}, fmt.Errorf("error reading keys: %w", err)
	}
	values, err := d.readInterfaces(n)
	if err != nil {
		return map[string]interface{}{}, fmt.Errorf("error reading values: %w", err)
	}
	m := make(map[string]interface{}, len(keys))
	for i := 0; i < n; i++ {
		m[keys[i]] = values[i]
	}
	return m, nil
}

func (d *LiteralDecoder) readInterface(xh byte) (interface{}, error) {
	if xh&MaskType == d.typeBool {
		xd, err := d.readBool()
		if err != nil {
			return nil, fmt.Errorf("error reading bool: %w", err)
		}
		return xd, nil
	}
	if xh&MaskType == d.typeUInt2 {
		return int((xh & MaskUpper2) >> 6), nil
	}
	if xh&MaskType == d.typeUInt3 {
		return int((xh & MaskUpper3) >> 5), nil
	}
	if xh == d.typeUInt8 {
		xd, err := d.readByte()
		if err != nil {
			return nil, fmt.Errorf("error reading int8: %w", err)
		}
		return int(int8(xd)), nil
	}
	if xh == d.typeUInt16 {
		xd, err := d.readUInt16()
		if err != nil {
			return nil, fmt.Errorf("error reading uint16: %w", err)
		}
		return int(xd), nil
	}
	if xh == d.typeInt32 {
		xd, err := d.readInt32()
		if err != nil {
			return nil, fmt.Errorf("error reading int32: %w", err)
		}
		return int(xd), nil
	}
	if xh == d.typeInt64 {
		xd, err := d.readInt64()
		if err != nil {
			return nil, fmt.Errorf("error reading int64: %w", err)
		}
		return int(xd), nil
	}
	if xh == d.typeFloat32 {
		xd, err := d.readFloat32()
		if err != nil {
			return nil, fmt.Errorf("error reading float32: %w", err)
		}
		return xd, nil
	}
	if xh == d.typeFloat64 {
		xd, err := d.readFloat64()
		if err != nil {
			return nil, fmt.Errorf("error reading float64: %w", err)
		}
		return xd, nil
	}
	if xh == d.typeString {
		xs, err := d.readString()
		if err != nil {
			return nil, fmt.Errorf("error reading string: %w", err)
		}
		return xs, nil
	}
	return nil, fmt.Errorf("unsupported type %d", xh)
}

func (d *LiteralDecoder) readInterfaces(n int) ([]interface{}, error) {
	slc := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		xh, err := d.readByte()
		if err != nil {
			return slc, fmt.Errorf("error reading byte: %w", err)
		}
		if d.isMap(xh & MaskType) {
			if xh&MaskType == d.typeMapStringString {
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
			if xh&MaskType == d.typeMapStringInterface {
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
			if xh == d.typeMapString {
				xh2, err := d.readByte()
				if err != nil {
					return slc, fmt.Errorf("error reading byte: %w", err)
				}
				if xh2 == d.typeString {
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
				if xh2 == d.typeInterface {
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
			return slc, fmt.Errorf("unsupported type %d", xh)
		}
		if d.isArray(xh & MaskType) {
			if xh == d.typeArrayString {
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
			if xh == d.typeArrayFloat32 {
				size := 0
				err := d.decodeInt(&size)
				if err != nil {
					return slc, fmt.Errorf("error decoding size of map of string to string: %w", err)
				}
				values, err := d.readFloat32s(size)
				if err != nil {
					return slc, fmt.Errorf("error reading array of floats with size %d: %w", size, err)
				}
				slc = append(slc, values)
				continue
			}
			if xh == d.typeArrayFloat32x2 {
				values, err := d.readFloat32s(2)
				if err != nil {
					return slc, fmt.Errorf("error reading array of float32s with size 2: %w", err)
				}
				slc = append(slc, values)
				continue
			}
			if xh == d.typeArrayFloat64 {
				size := 0
				err := d.decodeInt(&size)
				if err != nil {
					return slc, fmt.Errorf("error decoding size of map of string to string: %w", err)
				}
				values, err := d.readFloat64s(size)
				if err != nil {
					return slc, fmt.Errorf("error reading array of float64s with size %d: %w", size, err)
				}
				slc = append(slc, values)
				continue
			}
			if xh == d.typeArrayFloat64x2 {
				values, err := d.readFloat64s(2)
				if err != nil {
					return slc, fmt.Errorf("error reading array of float64s with size 2: %w", err)
				}
				slc = append(slc, values)
				continue
			}
			if xh == d.typeArrayInterface {
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
			return slc, fmt.Errorf("unsupported type %d", xh)
		}
		xv, err := d.readInterface(xh)
		if err != nil {
			return slc, err
		}
		slc = append(slc, xv)
	}
	return slc, nil
}

func (d *LiteralDecoder) decodeBool(v *bool) error {
	h, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h != d.typeBool {
		return fmt.Errorf("unexpected type code, found %x, expecting bool type code %x", h, d.typeBool)
	}
	*v = int(((h&MaskUpper1)>>7)&MaskLower1) == 1
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
	switch h & MaskType {
	case d.typeUInt2:
		*v = int(((h & MaskUpper2) >> 6) & MaskLower2)
		return nil
	case d.typeUInt3:
		*v = int(((h & MaskUpper3) >> 5) & MaskLower3)
		return nil
	case d.typeUInt4:
		return errors.New("error decoding int4: not enough space in byte to encode int4: must be an error")
	case d.typeUInt8:
		d, err := d.readByte()
		if err != nil {
			return fmt.Errorf("error reading byte: %w", err)
		}
		*v = int(uint8(d))
		return nil
	case d.typeUInt16:
		d, err := d.readUInt16()
		if err != nil {
			return fmt.Errorf("error reading int64: %w", err)
		}
		*v = int(d)
	case d.typeInt32:
		d, err := d.readInt32()
		if err != nil {
			return fmt.Errorf("error reading int64: %w", err)
		}
		*v = int(d)
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

func (d *LiteralDecoder) decodeFloat32(v *float32) error {
	h, err := d.reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding header: %w", err)
	}
	if h != d.typeFloat32 {
		return fmt.Errorf("unexpected type code, found %x, expecting float32 type code %x", h, d.typeFloat32)
	}
	x, err := d.readFloat32()
	if err != nil {
		return fmt.Errorf("error reading float32: %w", err)
	}
	*v = x
	return nil
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
		return fmt.Errorf("error decoding slice of string: %w", err)
	}
	if h != d.typeArrayString {
		return fmt.Errorf("unexpected type code, found %x, expecting string array type code %x", h, d.typeArrayString)
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
		return fmt.Errorf("error decoding slice of bool: %w", err)
	}
	if h != d.typeArrayBool {
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
	if !d.isIntArray(h) {
		return fmt.Errorf("unexpected type code, found %d, expecting int array type code", h)
	}
	switch h {
	case d.typeArrayUInt2:
		size := 0
		err := d.decodeInt(&size)
		if err != nil {
			return fmt.Errorf("error decoding size of slice of int2s: %w", err)
		}
		slc := make([]int, 0, size)
		for i := 0; i < size/4; i++ {
			d, err := d.readByte()
			if err != nil {
				return fmt.Errorf("error decoding int4: %w", err)
			}
			slc = append(slc, int(d&MaskLower2))
			slc = append(slc, int((d>>2)&MaskLower2))
			slc = append(slc, int((d>>4)&MaskLower2))
			slc = append(slc, int((d>>6)&MaskLower2))
		}
		*v = slc
		return nil
	case d.typeArrayUInt4:
		size := 0
		err := d.decodeInt(&size)
		if err != nil {
			return fmt.Errorf("error decoding size of slice of int4s: %w", err)
		}
		slc := make([]int, 0, size)
		for i := 0; i < size/2; i++ {
			d, err := d.readByte()
			if err != nil {
				return fmt.Errorf("error decoding int4: %w", err)
			}
			slc = append(slc, int(d&MaskLower4))
			slc = append(slc, int((d>>4)&MaskLower4))
		}
		*v = slc
		return nil
	case d.typeArrayUInt16:
		size := 0
		err := d.decodeInt(&size)
		if err != nil {
			return fmt.Errorf("error decoding size of slice of int16s: %w", err)
		}
		slc := make([]int, 0, size)
		for i := 0; i < size; i++ {
			x, err := d.readUInt16()
			if err != nil {
				return fmt.Errorf("error reading int16: %w", err)
			}
			slc = append(slc, int(x))
		}
		*v = slc
		return nil
	case d.typeArrayInt32:
		size := 0
		err := d.decodeInt(&size)
		if err != nil {
			return fmt.Errorf("error decoding size of slice of int32s: %w", err)
		}
		slc := make([]int, 0, size)
		for i := 0; i < size; i++ {
			x, err := d.readInt32()
			if err != nil {
				return fmt.Errorf("error reading int32: %w", err)
			}
			slc = append(slc, int(x))
		}
		*v = slc
		return nil
	case d.typeArrayInt64:
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

func (d *LiteralDecoder) decodeSliceFloat32(v *[]float32) error {
	h, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h == d.typeArrayFloat32x2 {
		slc, err := d.readFloat32s(2)
		if err != nil {
			return fmt.Errorf("error reading map of string to string: %w", err)
		}
		*v = slc
		return nil
	}
	if h == d.typeArrayFloat32 {
		size := 0
		if err := d.decodeInt(&size); err != nil {
			return fmt.Errorf("error decoding size of slice of float32s: %w", err)
		}
		slc, err := d.readFloat32s(size)
		if err != nil {
			return fmt.Errorf("error reading map of string to string: %w", err)
		}
		*v = slc
		return nil
	}
	return fmt.Errorf("unexpected type code, found %x, expecting float32 array type code", h)
}

func (d *LiteralDecoder) decodeSliceFloat64(v *[]float64) error {
	h, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h == d.typeArrayFloat64x2 {
		slc, err := d.readFloat64s(2)
		if err != nil {
			return fmt.Errorf("error reading map of string to string: %w", err)
		}
		*v = slc
		return nil
	}
	if h == d.typeArrayFloat64 {
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
	return fmt.Errorf("unexpected type code, found %x, expecting float64 array type code", h)
}

func (d *LiteralDecoder) decodeSliceInterface(v *[]interface{}) error {
	h, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h != d.typeArrayInterface {
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
	h0, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h0 == d.typeMapStringString {
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
	h1, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h0 != d.typeMapString || h1 != d.typeString {
		return fmt.Errorf("unexpected type code, found 0x%x 0x%x, expecting map string to string type codes", h0, h1)
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
	if h[0] != d.typeMapString || !d.isInt(h[1]&MaskType) {
		return fmt.Errorf("unexpected type code, found %x, expecting map string to string type codes", h)
	}
	if h[1] == d.typeUInt8 {
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
	if h[1] == d.typeUInt16 {
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
			value, err := d.readUInt16()
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
	if h[1] == d.typeInt32 {
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
			value, err := d.readInt32()
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
	if h[1] == d.typeInt64 {
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
	h0, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h0 == d.typeMapStringInterface {
		size := 0
		if err := d.decodeInt(&size); err != nil {
			return fmt.Errorf("error decoding size of map of string to interface: %w", err)
		}
		m, err := d.readMapStringInterface(size)
		if err != nil {
			return fmt.Errorf("error reading map of string to interface: %w", err)
		}
		*v = m
		return nil
	}
	h1, err := d.readByte()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("error decoding bool: %w", err)
	}
	if h0 != d.typeMapString || h1 != d.typeInterface {
		return fmt.Errorf("unexpected type code, found 0x%x 0x%x, expecting map string to interface type codes", h0, h1)
	}
	size := 0
	if err := d.decodeInt(&size); err != nil {
		return fmt.Errorf("error decoding size of map of string to interface: %w", err)
	}
	m, err := d.readMapStringInterface(size)
	if err != nil {
		return fmt.Errorf("error reading for map of string to interface: %w", err)
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
	if f, ok := v.(*float32); ok {
		return d.decodeFloat32(f)
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
	if slc, ok := v.(*[]float32); ok {
		return d.decodeSliceFloat32(slc)
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
