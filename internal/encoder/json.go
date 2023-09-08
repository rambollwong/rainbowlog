package encoder

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

// JSONMarshalFunc is used to marshal interface to JSON encoded byte slice.
// Making it package level instead of embedded in Encoder brings
// some extra efforts at importing, but avoids value copy when the functions
// of Encoder being invoked.
var JSONMarshalFunc func(v interface{}) ([]byte, error) = json.Marshal

type JsonEncoder struct{}

func (j JsonEncoder) MetaEnd(dst []byte) []byte {
	return dst
}

func (j JsonEncoder) Key(dst []byte, key string) []byte {
	if len(dst) > 0 && dst[len(dst)-1] != '{' {
		dst = j.Delim(dst)
	}
	dst = j.String(dst, key)
	dst = append(dst, ':')
	return dst
}

// BlankSpace append ' ' to dst.
func (j JsonEncoder) BlankSpace(dst []byte) []byte {
	return append(dst, ' ')
}

// Comma append ',' to dst.
func (j JsonEncoder) Comma(dst []byte) []byte {
	return append(dst, ',')
}

func (j JsonEncoder) Delim(dst []byte) []byte {
	if len(dst) > 0 {
		dst = append(dst, ',')
	}
	return dst
}

func (j JsonEncoder) ArrayDelim(dst []byte) []byte {
	if len(dst) > 0 {
		dst = append(dst, ',')
	}
	return dst
}

func (j JsonEncoder) ArrayEnd(dst []byte) []byte {
	dst = append(dst, ']')
	return dst
}

func (j JsonEncoder) ArrayStart(dst []byte) []byte {
	return append(dst, '[')
}

func (j JsonEncoder) BeginMarker(dst []byte) []byte {
	return append(dst, '{')
}

func (j JsonEncoder) EndMarker(dst []byte) []byte {
	return append(dst, '}')
}

func (j JsonEncoder) IPAddr(dst []byte, ip net.IP) []byte {
	return j.String(dst, ip.String())
}

func (j JsonEncoder) IPPrefix(dst []byte, pfx net.IPNet) []byte {
	return j.String(dst, pfx.String())
}

func (j JsonEncoder) Interface(dst []byte, i interface{}) []byte {
	marshaled, err := JSONMarshalFunc(i)
	if err != nil {
		return j.String(dst, fmt.Sprintf("marshaling error: %v", err))
	}
	return append(dst, marshaled...)
}

func (j JsonEncoder) LineBreak(dst []byte) []byte {
	return append(dst, '\n')
}

func (j JsonEncoder) MACAddr(dst []byte, ha net.HardwareAddr) []byte {
	return j.String(dst, ha.String())
}

func (j JsonEncoder) Nil(dst []byte) []byte {
	return append(dst, "null"...)
}

func (j JsonEncoder) ObjectData(dst []byte, o []byte) []byte {
	if o[0] == '{' {
		if len(dst) > 1 {
			dst = j.Delim(dst)
		}
		o = o[1:]
	} else if len(dst) > 1 {
		dst = j.Delim(dst)
	}
	return append(dst, o...)
}

func (JsonEncoder) Bool(dst []byte, val bool) []byte {
	return strconv.AppendBool(dst, val)
}

func (j JsonEncoder) Bools(dst []byte, val ...bool) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Bool(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Delim(dst)
			dst = j.Bool(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}
