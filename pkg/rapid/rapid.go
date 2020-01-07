package rapid

import (
	"math"
)

const (
	DefaultTypeInterface = 1 // reserved: not a real type

	DefaultTypeBool    = 2
	DefaultTypeUInt2   = 3
	DefaultTypeUInt3   = 4
	DefaultTypeUInt4   = 5
	DefaultTypeUInt8   = 6
	DefaultTypeUInt16  = 7
	DefaultTypeInt32   = 8
	DefaultTypeInt64   = 9
	DefaultTypeFloat32 = 10
	DefaultTypeFloat64 = 11
	DefaultTypeString  = 12

	DefaultTypeArrayInterface = 13
	DefaultTypeArrayBool      = 14
	DefaultTypeArrayUInt2     = 15
	DefaultTypeArrayUInt3     = 16
	DefaultTypeArrayUInt4     = 17
	DefaultTypeArrayUInt8     = 18
	DefaultTypeArrayUInt16    = 19
	DefaultTypeArrayInt32     = 20
	DefaultTypeArrayInt64     = 21
	DefaultTypeArrayFloat32   = 22
	DefaultTypeArrayFloat32x2 = 23
	DefaultTypeArrayFloat64   = 24
	DefaultTypeArrayFloat64x2 = 25
	DefaultTypeArrayString    = 26

	DefaultTypeMapString          = 27
	DefaultTypeMapStringString    = 28
	DefaultTypeMapStringInterface = 29

	MaskUpper1 = 0b10000000
	MaskUpper2 = 0b11000000
	MaskUpper3 = 0b11100000
	MaskUpper4 = 0b11110000

	MaskLower1 = 0b00000001
	MaskLower2 = 0b00000011
	MaskLower3 = 0b00000111
	MaskLower4 = 0b00001111
	MaskLower5 = 0b00011111
	MaskLower6 = 0b00111111

	MaskType = MaskLower5

	NullByte = byte(0)
)

var (
	MaxExactInt32 = math.Pow(2, 24)
	MaxExactInt64 = math.Pow(2, 53)
)

var (
	MagicNumber = []byte("RAPID")
)

func flatten(v [][]byte) []byte {
	out := make([]byte, 0)
	for _, b := range v {
		out = append(out, b...)
	}
	return out
}
