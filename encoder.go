package rainbowlog

import (
	"github.com/rambollwong/rainbowlog/internal/encoder"
	"net"
	"time"
)

var JsonEnc Encoder = encoder.JsonEncoder{}
var TextEnc Encoder = encoder.TextEncoder{
	MetaKeys: defaultMetaKeys().Keys(),
}
var consoleEnc Encoder = encoder.TextEncoder{
	MetaKeys: defaultMetaKeys().Keys(),
}

var NewTextEncoder = func(metaKeys ...string) Encoder {
	return encoder.TextEncoder{MetaKeys: metaKeys}
}

type Encoder interface {
	encoderWithArray
	BeginMarker(dst []byte) []byte
	BlankSpace(dst []byte) []byte
	Bool(dst []byte, val bool) []byte
	Bytes(dst []byte, s []byte) []byte
	Comma(dst []byte) []byte
	Delim(dst []byte) []byte
	Duration(dst []byte, unit time.Duration, useInt bool, d time.Duration) []byte
	EndMarker(dst []byte) []byte
	Float32(dst []byte, val float32) []byte
	Float64(dst []byte, val float64) []byte
	Hex(dst []byte, s []byte) []byte
	IPAddr(dst []byte, ip net.IP) []byte
	IPPrefix(dst []byte, pfx net.IPNet) []byte
	Int(dst []byte, val int) []byte
	Int16(dst []byte, val int16) []byte
	Int32(dst []byte, val int32) []byte
	Int64(dst []byte, val int64) []byte
	Int8(dst []byte, val int8) []byte
	Interface(dst []byte, i interface{}) []byte
	Key(dst []byte, key string) []byte
	LineBreak(dst []byte) []byte
	MACAddr(dst []byte, ha net.HardwareAddr) []byte
	Nil(dst []byte) []byte
	ObjectData(dst []byte, o []byte) []byte
	String(dst []byte, s string) []byte
	Time(dst []byte, format string, t time.Time) []byte
	Uint(dst []byte, val uint) []byte
	Uint16(dst []byte, val uint16) []byte
	Uint32(dst []byte, val uint32) []byte
	Uint64(dst []byte, val uint64) []byte
	Uint8(dst []byte, val uint8) []byte
	MetaEnd(dst []byte) []byte
}

type encoderWithArray interface {
	ArrayDelim(dst []byte) []byte
	ArrayEnd(dst []byte) []byte
	ArrayStart(dst []byte) []byte
	Bools(dst []byte, val ...bool) []byte
	Durations(dst []byte, unit time.Duration, useInt bool, d ...time.Duration) []byte
	Float32s(dst []byte, val ...float32) []byte
	Float64s(dst []byte, val ...float64) []byte
	Ints(dst []byte, val ...int) []byte
	Int16s(dst []byte, val ...int16) []byte
	Int32s(dst []byte, val ...int32) []byte
	Int64s(dst []byte, val ...int64) []byte
	Int8s(dst []byte, val ...int8) []byte
	Strings(dst []byte, s ...string) []byte
	Times(dst []byte, format string, t ...time.Time) []byte
	Uints(dst []byte, val ...uint) []byte
	Uint16s(dst []byte, val ...uint16) []byte
	Uint32s(dst []byte, val ...uint32) []byte
	Uint64s(dst []byte, val ...uint64) []byte
	Uint8s(dst []byte, val ...uint8) []byte
}
