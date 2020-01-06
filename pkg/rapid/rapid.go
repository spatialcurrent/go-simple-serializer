package rapid

const (
	DefaultTypeBool      = 1
	DefaultTypeInt8      = 2
	DefaultTypeInt64     = 4
	DefaultTypeFloat64   = 8
	DefaultTypeString    = 16
	DefaultTypeArray     = 32
	DefaultTypeMap       = 64
	DefaultTypeInterface = 128
)

func flatten(v [][]byte) []byte {
	out := make([]byte, 0)
	for _, b := range v {
		out = append(out, b...)
	}
	return out
}
