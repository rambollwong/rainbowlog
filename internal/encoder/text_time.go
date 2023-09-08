package encoder

import (
	"time"
)

func (j TextEncoder) Duration(dst []byte, unit time.Duration, useInt bool, d time.Duration) []byte {
	if useInt {
		return j.Int64(dst, int64(d/unit))
	}
	return j.Float64(dst, float64(d)/float64(unit))
}

func (j TextEncoder) Time(dst []byte, format string, t time.Time) []byte {
	switch format {
	case TimeFormatUnix:
		return j.Int64(dst, t.Unix())
	case TimeFormatUnixMs:
		return j.Int64(dst, t.UnixNano()/1000000)
	case TimeFormatUnixMicro:
		return j.Int64(dst, t.UnixNano()/1000)
	default:
		dst = t.AppendFormat(dst, format)
		return dst
	}
}

func (j TextEncoder) Durations(dst []byte, unit time.Duration, useInt bool, d ...time.Duration) []byte {
	if len(d) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Duration(dst, unit, useInt, d[0])
	if len(d) > 1 {
		for _, v := range d {
			dst = j.Comma(dst)
			dst = j.Duration(dst, unit, useInt, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}

func (j TextEncoder) Times(dst []byte, format string, val ...time.Time) []byte {
	if len(val) == 0 {
		dst = j.ArrayStart(dst)
		dst = j.ArrayEnd(dst)
		return dst
	}
	dst = j.ArrayStart(dst)
	dst = j.Time(dst, format, val[0])
	if len(val) > 1 {
		for _, v := range val[1:] {
			dst = j.Comma(dst)
			dst = j.Time(dst, format, v)
		}
	}
	dst = j.ArrayEnd(dst)
	return dst
}
