package encoder

import (
	"strconv"
)

func (j JsonEncoder) Float32(dst []byte, val float32) []byte {
	return j.Float64(dst, float64(val))
}

func (j JsonEncoder) Float64(dst []byte, val float64) []byte {
	return appendFloat(dst, val, 32)
}

func (j JsonEncoder) Int(dst []byte, val int) []byte {
	return j.Int64(dst, int64(val))
}

func (j JsonEncoder) Int16(dst []byte, val int16) []byte {
	return j.Int64(dst, int64(val))
}

func (j JsonEncoder) Int32(dst []byte, val int32) []byte {
	return j.Int64(dst, int64(val))
}

func (j JsonEncoder) Int64(dst []byte, val int64) []byte {
	return strconv.AppendInt(dst, val, 10)
}

func (j JsonEncoder) Int8(dst []byte, val int8) []byte {
	return j.Int64(dst, int64(val))
}

func (j JsonEncoder) Uint(dst []byte, val uint) []byte {
	return j.Uint64(dst, uint64(val))
}

func (j JsonEncoder) Uint16(dst []byte, val uint16) []byte {
	return j.Uint64(dst, uint64(val))
}

func (j JsonEncoder) Uint32(dst []byte, val uint32) []byte {
	return j.Uint64(dst, uint64(val))
}

func (j JsonEncoder) Uint64(dst []byte, val uint64) []byte {
	return strconv.AppendUint(dst, val, 10)
}

func (j JsonEncoder) Uint8(dst []byte, val uint8) []byte {
	return j.Uint64(dst, uint64(val))
}

func (j JsonEncoder) Float32s(dst []byte, val ...float32) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Float32(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Float32(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}

func (j JsonEncoder) Float64s(dst []byte, val ...float64) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Float64(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Float64(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}

func (j JsonEncoder) Ints(dst []byte, val ...int) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Int(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Int(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}

func (j JsonEncoder) Int16s(dst []byte, val ...int16) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Int16(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Int16(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}

func (j JsonEncoder) Int32s(dst []byte, val ...int32) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Int32(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Int32(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}

func (j JsonEncoder) Int64s(dst []byte, val ...int64) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Int64(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Int64(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}

func (j JsonEncoder) Int8s(dst []byte, val ...int8) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Int8(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Int8(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}

func (j JsonEncoder) Uints(dst []byte, val ...uint) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Uint(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Uint(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}

func (j JsonEncoder) Uint16s(dst []byte, val ...uint16) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Uint16(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Uint16(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}

func (j JsonEncoder) Uint32s(dst []byte, val ...uint32) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Uint32(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Uint32(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}

func (j JsonEncoder) Uint64s(dst []byte, val ...uint64) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Uint64(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Uint64(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}

func (j JsonEncoder) Uint8s(dst []byte, val ...uint8) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Uint8(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Uint8(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}
