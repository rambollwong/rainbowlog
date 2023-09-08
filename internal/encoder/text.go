package encoder

import (
	"fmt"
	"net"
	"strconv"
)

// TextMarshalFunc is used to marshal interface to text byte slice.
var TextMarshalFunc func(v interface{}) ([]byte, error) = func(v interface{}) ([]byte, error) {
	return []byte(fmt.Sprintf("%v", v)), nil
}

type TextEncoder struct {
	MetaKeys []string
}

func (j TextEncoder) isMetaKey(key string) bool {
	for _, metaKey := range j.MetaKeys {
		if metaKey == key {
			return true
		}
	}
	return false
}

func (j TextEncoder) MetaEnd(dst []byte) []byte {
	if len(dst) > 0 {
		dst = append(dst, ' ', '>')
	}
	return dst
}

func (j TextEncoder) Key(dst []byte, key string) []byte {
	if len(dst) > 1 && dst[len(dst)-1] != ' ' {
		dst = j.Delim(dst)
	}
	if j.isMetaKey(key) {
		return dst
	}
	dst = j.String(dst, key)
	dst = append(dst, '=')
	return dst
}

// BlankSpace append ' ' to dst.
func (j TextEncoder) BlankSpace(dst []byte) []byte {
	return append(dst, ' ')
}

// Comma append ',' to dst.
func (j TextEncoder) Comma(dst []byte) []byte {
	return append(dst, ',')
}

func (j TextEncoder) Delim(dst []byte) []byte {
	if len(dst) > 0 {
		dst = append(dst, ' ')
	}
	return dst
}

func (j TextEncoder) ArrayDelim(dst []byte) []byte {
	if len(dst) > 0 {
		dst = append(dst, ',')
	}
	return dst
}

func (j TextEncoder) ArrayEnd(dst []byte) []byte {
	dst = append(dst, ']')
	return dst
}

func (j TextEncoder) ArrayStart(dst []byte) []byte {
	return append(dst, '[')
}

func (j TextEncoder) BeginMarker(dst []byte) []byte {
	return dst
}

func (j TextEncoder) EndMarker(dst []byte) []byte {
	return dst
}

func (j TextEncoder) IPAddr(dst []byte, ip net.IP) []byte {
	return j.String(dst, ip.String())
}

func (j TextEncoder) IPPrefix(dst []byte, pfx net.IPNet) []byte {
	return j.String(dst, pfx.String())
}

func (j TextEncoder) Interface(dst []byte, i interface{}) []byte {
	marshaled, err := JSONMarshalFunc(i)
	if err != nil {
		return j.String(dst, fmt.Sprintf("marshaling error: %v", err))
	}
	return append(dst, marshaled...)
}

func (j TextEncoder) LineBreak(dst []byte) []byte {
	return append(dst, '\n')
}

func (j TextEncoder) MACAddr(dst []byte, ha net.HardwareAddr) []byte {
	return j.String(dst, ha.String())
}

func (j TextEncoder) Nil(dst []byte) []byte {
	return append(dst, "null"...)
}

func (j TextEncoder) ObjectData(dst []byte, o []byte) []byte {
	if len(dst) > 3 && dst[len(dst)-1] != ' ' && dst[len(dst)-2] != '|' {
		dst = j.Delim(dst)
	}
	return append(dst, o...)
}

func (TextEncoder) Bool(dst []byte, val bool) []byte {
	return strconv.AppendBool(dst, val)
}

func (j TextEncoder) Bools(dst []byte, val ...bool) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Bool(dst, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Comma(dst)
			dst = j.Bool(dst, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}
