package rainbowlog

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
	"time"

	"github.com/rambollwong/rainbowlog/level"
)

var _ recordPacker = (*ConsolePacker)(nil)

type ConsolePacker struct {
	RecordPackerForWriter
	consoleColor bool
}

func (j *ConsolePacker) printCall(consoleColorsKey string, call func(j *ConsolePacker)) {
	var cs []int
	// key color printing
	if j.consoleColor {
		cs = j.record.logger.metaKeys.ConsoleColors(consoleColorsKey)
		if cs != nil {
			for _, c := range cs {
				*j.raw = colorStart(*j.raw, c)
			}
		}
	}
	// field print
	call(j)
	// end key color printing
	if cs != nil {
		for i := 0; i < len(cs); i++ {
			*j.raw = colorEnd(*j.raw)
		}
	}
}

func (j *ConsolePacker) Msg(msg string) {
	j.printCall(MsgFieldName, func(j *ConsolePacker) {
		// key & value field print
		j.RecordPackerForWriter.Msg(msg)
	})
}

func (j *ConsolePacker) Err(err error) {
	if err == nil || err.Error() == "" {
		return
	}
	j.printCall(ErrFieldName, func(j *ConsolePacker) {
		j.RecordPackerForWriter.Err(err)
	})
}

func (j *ConsolePacker) printMetaLevelStart(dst *[]byte) int {
	if j.consoleColor {
		switch j.record.level {
		case level.Debug:
			cs := j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelDebugFieldName)
			if cs != nil {
				for _, c := range cs {
					*dst = colorStart(*dst, c)
				}
			}
			return len(cs)
		case level.Info:
			cs := j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelInfoFieldName)
			if cs != nil {
				for _, c := range cs {
					*dst = colorStart(*dst, c)
				}
			}
			return len(cs)
		case level.Warn:
			cs := j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelWarnFieldName)
			if cs != nil {
				for _, c := range cs {
					*dst = colorStart(*dst, c)
				}
			}
			return len(cs)
		case level.Error:
			cs := j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelErrorFieldName)
			if cs != nil {
				for _, c := range cs {
					*dst = colorStart(*dst, c)
				}
			}
			return len(cs)
		case level.Panic:
			cs := j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelPanicFieldName)
			if cs != nil {
				for _, c := range cs {
					*dst = colorStart(*dst, c)
				}
			}
			return len(cs)
		case level.Fatal:
			cs := j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelFatalFieldName)
			if cs != nil {
				for _, c := range cs {
					*dst = colorStart(*dst, c)
				}
			}
			return len(cs)
		case level.Trace:
			cs := j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelTraceFieldName)
			if cs != nil {
				for _, c := range cs {
					*dst = colorStart(*dst, c)
				}
			}
			return len(cs)
		default:
			return 0
		}
	}
	return 0
}

func (j *ConsolePacker) printFirstMetaStart(dst *[]byte) int {
	if j.consoleColor {
		var cs []int
		switch j.record.level {
		case level.Debug:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelDebugLineFieldName)
		case level.Info:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelInfoLineFieldName)
		case level.Warn:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelWarnLineFieldName)
		case level.Error:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelErrorLineFieldName)
		case level.Panic:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelPanicLineFieldName)
		case level.Fatal:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelFatalLineFieldName)
		case level.Trace:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelTraceLineFieldName)
		default:
			return 0
		}
		if cs != nil {
			for _, c := range cs {
				*dst = colorStart(*dst, c)
			}
			*dst = colorStart(*dst, ColorBlack)
			return len(cs) + 1
		}
	}
	return 0
}

func (j *ConsolePacker) printRainbowStart(dst *[]byte, metaIdx int, metaKey string) int {
	if !j.consoleColor {
		return 0
	}
	var cs []int
	var endI int
	if metaIdx == 0 {
		// First meta printing with rainbow
		endI = j.printFirstMetaStart(dst)
	} else {
		// Others meta printing
		cs = j.record.logger.metaKeys.ConsoleColors(metaKey)
		if cs != nil {
			for _, c := range cs {
				*dst = colorStart(*dst, c)
			}
			endI = len(cs)
		}
	}
	return endI
}

func (j *ConsolePacker) printRainbowEnd(dst *[]byte, metaIdx int, i int) {
	if j.consoleColor && i > 0 {
		if metaIdx == 0 {
			*dst = j.writerEncoderPair.enc.BlankSpace(*dst)
		}
		for t := 0; t < i; t++ {
			*dst = colorEnd(*dst)
		}
	}
}

func (j *ConsolePacker) createMeta() {
	j.meta = bytesPool.Get()
	// begin
	*j.meta = j.writerEncoderPair.enc.BeginMarker(*j.meta)
	if j.record.logger.metaKeys == nil {
		// if meta keys unset, use default
		j.record.logger.metaKeys = defaultMetaKeys()
	}
	// range meta keys
	for i, s := range j.record.logger.metaKeys.Keys() {
		switch s {
		case MetaTimeFieldName:
			endI := j.printRainbowStart(j.meta, i, s)
			*j.meta = j.writerEncoderPair.enc.Key(*j.meta, MetaTimeFieldName)
			*j.meta = j.writerEncoderPair.enc.Time(*j.meta, j.record.logger.timeFormat, time.Now())
			j.printRainbowEnd(j.meta, i, endI)
		case MetaLabelFieldName:
			if j.record.label == "" {
				continue
			}
			endI := j.printRainbowStart(j.meta, i, s)
			*j.meta = j.writerEncoderPair.enc.Key(*j.meta, MetaLabelFieldName)
			*j.meta = j.writerEncoderPair.enc.String(*j.meta, j.record.label)
			j.printRainbowEnd(j.meta, i, endI)
		case MetaLevelFieldName:
			var endI int
			if i == 0 {
				// First meta printing with rainbow
				endI = j.printFirstMetaStart(j.meta)
			} else {
				// Others meta printing
				endI = j.printMetaLevelStart(j.meta)
			}
			*j.meta = j.writerEncoderPair.enc.Key(*j.meta, MetaLevelFieldName)
			*j.meta = j.writerEncoderPair.enc.String(*j.meta, j.record.level.KeyFieldValue())
			j.printRainbowEnd(j.meta, i, endI)
		case MetaCallerFieldName:
			if j.record.logger.callerMarshalFunc == nil {
				continue
			}
			skip := j.callerSkipFrameCount + CallerSkipFrameCount + innerCallerSkipFrameCount
			_, file, line, ok := runtime.Caller(skip)
			if !ok {
				continue
			}
			endI := j.printRainbowStart(j.meta, i, s)
			*j.meta = j.writerEncoderPair.enc.Key(*j.meta, MetaCallerFieldName)
			*j.meta = j.writerEncoderPair.enc.String(*j.meta, j.record.logger.callerMarshalFunc(file, line))
			j.printRainbowEnd(j.meta, i, endI)
		}
	}
}

func (j *ConsolePacker) Done() {
	// end raw
	*j.raw = j.writerEncoderPair.enc.EndMarker(*j.raw)
	// create meta
	j.createMeta()
	// append meta end
	var endI int
	if j.consoleColor {
		cs := j.record.logger.metaKeys.ConsoleColors(metaEndFieldName)
		if cs != nil {
			for _, c := range cs {
				*j.meta = colorStart(*j.meta, c)
			}
			endI = len(cs)
		}
	}
	*j.meta = j.writerEncoderPair.enc.MetaEnd(*j.meta)
	j.printRainbowEnd(j.meta, 1, endI)
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

func (j *ConsolePacker) Str(key, val string) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.String(*j.raw, val)
}

func (j *ConsolePacker) Strs(key string, vals ...string) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Strings(*j.raw, vals...)
}

func (j *ConsolePacker) Stringer(key string, val fmt.Stringer) {
	j.Str(key, val.String())
}

func (j *ConsolePacker) Stringers(key string, vals ...fmt.Stringer) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.ArrayStart(*j.raw)
	for i, val := range vals {
		if i > 0 {
			*j.raw = j.writerEncoderPair.enc.ArrayDelim(*j.raw)
		}
		*j.raw = j.writerEncoderPair.enc.String(*j.raw, val.String())
	}
	*j.raw = j.writerEncoderPair.enc.ArrayEnd(*j.raw)
}

func (j *ConsolePacker) Bytes(key string, val []byte) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Bytes(*j.raw, val)
}

func (j *ConsolePacker) Hex(key string, val []byte) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Hex(*j.raw, val)
}

func (j *ConsolePacker) Int(key string, val int) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Int(*j.raw, val)
}

func (j *ConsolePacker) Ints(key string, vals ...int) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Ints(*j.raw, vals...)
}

func (j *ConsolePacker) Int8(key string, val int8) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Int8(*j.raw, val)
}

func (j *ConsolePacker) Int8s(key string, vals ...int8) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Int8s(*j.raw, vals...)
}

func (j *ConsolePacker) Int16(key string, val int16) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Int16(*j.raw, val)
}

func (j *ConsolePacker) Int16s(key string, vals ...int16) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Int16s(*j.raw, vals...)
}

func (j *ConsolePacker) Int32(key string, val int32) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Int32(*j.raw, val)
}

func (j *ConsolePacker) Int32s(key string, vals ...int32) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Int32s(*j.raw, vals...)
}

func (j *ConsolePacker) Int64(key string, val int64) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Int64(*j.raw, val)
}

func (j *ConsolePacker) Int64s(key string, vals ...int64) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Int64s(*j.raw, vals...)
}

func (j *ConsolePacker) Uint(key string, val uint) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Uint(*j.raw, val)
}

func (j *ConsolePacker) Uints(key string, vals ...uint) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Uints(*j.raw, vals...)
}

func (j *ConsolePacker) Uint8(key string, val uint8) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Uint8(*j.raw, val)
}

func (j *ConsolePacker) Uint8s(key string, vals ...uint8) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Uint8s(*j.raw, vals...)
}

func (j *ConsolePacker) Uint16(key string, val uint16) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Uint16(*j.raw, val)
}

func (j *ConsolePacker) Uint16s(key string, vals ...uint16) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Uint16s(*j.raw, vals...)
}

func (j *ConsolePacker) Uint32(key string, val uint32) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Uint32(*j.raw, val)
}

func (j *ConsolePacker) Uint32s(key string, vals ...uint32) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Uint32s(*j.raw, vals...)
}

func (j *ConsolePacker) Uint64(key string, val uint64) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Uint64(*j.raw, val)
}

func (j *ConsolePacker) Uint64s(key string, vals ...uint64) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Uint64s(*j.raw, vals...)
}

func (j *ConsolePacker) Float32(key string, val float32) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Float32(*j.raw, val)
}

func (j *ConsolePacker) Float32s(key string, vals ...float32) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Float32s(*j.raw, vals...)
}

func (j *ConsolePacker) Float64(key string, val float64) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Float64(*j.raw, val)
}

func (j *ConsolePacker) Float64s(key string, vals ...float64) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Float64s(*j.raw, vals...)
}

func (j *ConsolePacker) Time(key, fmt string, val time.Time) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Time(*j.raw, fmt, val)
}

func (j *ConsolePacker) Times(key, fmt string, vals ...time.Time) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Times(*j.raw, fmt, vals...)
}

func (j *ConsolePacker) Dur(key string, unit, val time.Duration) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Duration(*j.raw, unit, j.record.useIntDur, val)
}

func (j *ConsolePacker) Durs(key string, unit time.Duration, vals ...time.Duration) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Durations(*j.raw, unit, j.record.useIntDur, vals...)
}

func (j *ConsolePacker) Any(key string, i any) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.Interface(*j.raw, i)
}

func (j *ConsolePacker) IPAddr(key string, ip net.IP) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.IPAddr(*j.raw, ip)
}

func (j *ConsolePacker) IPPrefix(key string, pfx net.IPNet) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.IPPrefix(*j.raw, pfx)
}

func (j *ConsolePacker) MACAddr(key string, ha net.HardwareAddr) {
	j.printCall(keysColorName, func(j *ConsolePacker) {
		*j.raw = j.writerEncoderPair.enc.Key(*j.raw, key)
	})
	*j.raw = j.writerEncoderPair.enc.MACAddr(*j.raw, ha)
}

func colorStart(dst []byte, color int) []byte {
	s := "\x1b[" + strconv.Itoa(color) + "m"
	return append(dst, s...)
}

func colorEnd(dst []byte) []byte {
	return append(dst, "\x1b[0m"...)
}
