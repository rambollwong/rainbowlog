package rainbowlog

import (
	"fmt"
	"net"
	"runtime"
	"time"
)

type recordPacker interface {
	Reset()
	CallerSkip(skip int)
	Msg(msg string)
	Err(err error)
	Done()

	Str(key, val string)
	Strs(key string, vals ...string)
	Stringer(key string, val fmt.Stringer)
	Stringers(key string, vals ...fmt.Stringer)
	Bytes(key string, val []byte)
	Hex(key string, val []byte)
	Int(key string, val int)
	Ints(key string, vals ...int)
	Int8(key string, val int8)
	Int8s(key string, vals ...int8)
	Int16(key string, val int16)
	Int16s(key string, vals ...int16)
	Int32(key string, val int32)
	Int32s(key string, vals ...int32)
	Int64(key string, val int64)
	Int64s(key string, vals ...int64)
	Uint(key string, val uint)
	Uints(key string, vals ...uint)
	Uint8(key string, val uint8)
	Uint8s(key string, vals ...uint8)
	Uint16(key string, val uint16)
	Uint16s(key string, vals ...uint16)
	Uint32(key string, val uint32)
	Uint32s(key string, vals ...uint32)
	Uint64(key string, val uint64)
	Uint64s(key string, vals ...uint64)
	Float32(key string, val float32)
	Float32s(key string, vals ...float32)
	Float64(key string, val float64)
	Float64s(key string, vals ...float64)
	Time(key, fmt string, val time.Time)
	Times(key, fmt string, vals ...time.Time)
	Dur(key string, unit, val time.Duration)
	Durs(key string, unit time.Duration, vals ...time.Duration)
	Interface(key string, i interface{})
	IPAddr(key string, ip net.IP)
	IPPrefix(key string, pfx net.IPNet)
	MACAddr(key string, ha net.HardwareAddr)
}

var _ recordPacker = (*RecordPackerForWriter)(nil)

type RecordPackerForWriter struct {
	record               *Record
	meta                 *[]byte
	raw                  *[]byte
	callerSkipFrameCount int
	writerEncoderPair    *WriterEncoderPair
}

func (j *RecordPackerForWriter) Reset() {
	// reset raw
	if j.raw != nil {
		bytesPool.Put(j.raw)
	}
	j.raw = bytesPool.Get()
	// reset meta
	if j.meta != nil {
		bytesPool.Put(j.meta)
	}
	j.meta = bytesPool.Get()
	// reset caller skip frame count
	j.callerSkipFrameCount = 0
}

func (j *RecordPackerForWriter) CallerSkip(skip int) {
	j.callerSkipFrameCount += skip
}

func (j *RecordPackerForWriter) Msg(msg string) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, MsgFieldName)
	*j.raw = j.writerEncoderPair.enc.String(*j.raw, msg)
}

func (j *RecordPackerForWriter) Err(err error) {
	if err == nil || err.Error() == "" {
		return
	}
	j.Str(ErrFieldName, j.record.logger.errorMarshalFunc(err))
	if j.record.stack && GlobalErrorStackMarshalFunc != nil {
		switch esm := j.record.logger.errorStackMarshalFunc(err).(type) {
		case nil:
		case error:
			j.Str(ErrStackFieldName, esm.Error())
		case string:
			j.Str(ErrStackFieldName, esm)
		default:
			j.Interface(ErrStackFieldName, esm)
		}
	}
}

func (j *RecordPackerForWriter) createMeta() {
	j.meta = bytesPool.Get()
	// begin
	*j.meta = j.writerEncoderPair.enc.BeginMarker(*j.meta)
	if j.record.logger.metaKeys == nil {
		return
	}
	// range meta keys
	for _, s := range j.record.logger.metaKeys.Keys() {
		switch s {
		case MetaTimeFieldName:
			*j.meta = j.writerEncoderPair.enc.Key(*j.meta, MetaTimeFieldName)
			*j.meta = j.writerEncoderPair.enc.Time(*j.meta, j.record.logger.timeFormat, TimestampFunc())
		case MetaLabelFieldName:
			if j.record.label == "" {
				continue
			}
			*j.meta = j.writerEncoderPair.enc.Key(*j.meta, MetaLabelFieldName)
			*j.meta = j.writerEncoderPair.enc.String(*j.meta, j.record.label)
		case MetaLevelFieldName:
			*j.meta = j.writerEncoderPair.enc.Key(*j.meta, MetaLevelFieldName)
			*j.meta = j.writerEncoderPair.enc.String(*j.meta, j.record.logger.levelFieldMarshalFunc(j.record.level))
		case MetaCallerFieldName:
			skip := j.callerSkipFrameCount + CallerSkipFrameCount + innerCallerSkipFrameCount
			_, file, line, ok := runtime.Caller(skip)
			if !ok {
				continue
			}
			*j.meta = j.writerEncoderPair.enc.Key(*j.meta, MetaCallerFieldName)
			*j.meta = j.writerEncoderPair.enc.String(*j.meta, j.record.logger.callerMarshalFunc(file, line))
		}
	}
}

func (j *RecordPackerForWriter) Done() {
	// end raw
	*j.raw = j.writerEncoderPair.enc.EndMarker(*j.raw)
	// create meta
	j.createMeta()
	// append meta end
	*j.meta = j.writerEncoderPair.enc.MetaEnd(*j.meta)
	// append meta and raw
	*j.raw = j.writerEncoderPair.enc.ObjectData(*j.meta, *j.raw)
	// append end
	*j.raw = j.writerEncoderPair.enc.LineBreak(*j.raw)
	// write
	_, err := j.writerEncoderPair.writer.WriteLevel(j.record.level, *j.raw)
	if err != nil && ErrorHandler != nil {
		ErrorHandler(err)
	}
}

func (j *RecordPackerForWriter) Str(key, val string) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.String(*j.raw, val)
}

func (j *RecordPackerForWriter) Strs(key string, vals ...string) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Strings(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Stringer(key string, val fmt.Stringer) {
	j.Str(key, val.String())
}

func (j *RecordPackerForWriter) Stringers(key string, vals ...fmt.Stringer) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.ArrayStart(*j.raw)
	for i, val := range vals {
		if i > 0 {
			*j.raw = j.writerEncoderPair.enc.ArrayDelim(*j.raw)
		}
		*j.raw = j.writerEncoderPair.enc.String(*j.raw, val.String())
	}
	*j.raw = j.writerEncoderPair.enc.ArrayEnd(*j.raw)
}

func (j *RecordPackerForWriter) Bytes(key string, val []byte) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Bytes(*j.raw, val)
}

func (j *RecordPackerForWriter) Hex(key string, val []byte) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Hex(*j.raw, val)
}

func (j *RecordPackerForWriter) Int(key string, val int) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Int(*j.raw, val)
}

func (j *RecordPackerForWriter) Ints(key string, vals ...int) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Ints(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Int8(key string, val int8) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Int8(*j.raw, val)
}

func (j *RecordPackerForWriter) Int8s(key string, vals ...int8) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Int8s(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Int16(key string, val int16) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Int16(*j.raw, val)
}

func (j *RecordPackerForWriter) Int16s(key string, vals ...int16) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Int16s(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Int32(key string, val int32) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Int32(*j.raw, val)
}

func (j *RecordPackerForWriter) Int32s(key string, vals ...int32) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Int32s(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Int64(key string, val int64) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Int64(*j.raw, val)
}

func (j *RecordPackerForWriter) Int64s(key string, vals ...int64) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Int64s(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Uint(key string, val uint) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Uint(*j.raw, val)
}

func (j *RecordPackerForWriter) Uints(key string, vals ...uint) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Uints(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Uint8(key string, val uint8) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Uint8(*j.raw, val)
}

func (j *RecordPackerForWriter) Uint8s(key string, vals ...uint8) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Uint8s(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Uint16(key string, val uint16) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Uint16(*j.raw, val)
}

func (j *RecordPackerForWriter) Uint16s(key string, vals ...uint16) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Uint16s(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Uint32(key string, val uint32) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Uint32(*j.raw, val)
}

func (j *RecordPackerForWriter) Uint32s(key string, vals ...uint32) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Uint32s(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Uint64(key string, val uint64) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Uint64(*j.raw, val)
}

func (j *RecordPackerForWriter) Uint64s(key string, vals ...uint64) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Uint64s(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Float32(key string, val float32) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Float32(*j.raw, val)
}

func (j *RecordPackerForWriter) Float32s(key string, vals ...float32) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Float32s(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Float64(key string, val float64) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Float64(*j.raw, val)
}

func (j *RecordPackerForWriter) Float64s(key string, vals ...float64) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Float64s(*j.raw, vals...)
}

func (j *RecordPackerForWriter) Time(key, fmt string, val time.Time) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Time(*j.raw, fmt, val)
}

func (j *RecordPackerForWriter) Times(key, fmt string, vals ...time.Time) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Times(*j.raw, fmt, vals...)
}

func (j *RecordPackerForWriter) Dur(key string, unit, val time.Duration) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Duration(*j.raw, unit, j.record.useIntDur, val)
}

func (j *RecordPackerForWriter) Durs(key string, unit time.Duration, vals ...time.Duration) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Durations(*j.raw, unit, j.record.useIntDur, vals...)
}

func (j *RecordPackerForWriter) Interface(key string, i interface{}) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.Interface(*j.raw, i)
}

func (j *RecordPackerForWriter) IPAddr(key string, ip net.IP) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.IPAddr(*j.raw, ip)
}

func (j *RecordPackerForWriter) IPPrefix(key string, pfx net.IPNet) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.IPPrefix(*j.raw, pfx)
}

func (j *RecordPackerForWriter) MACAddr(key string, ha net.HardwareAddr) {
	*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	*j.raw = j.writerEncoderPair.enc.MACAddr(*j.raw, ha)
}
