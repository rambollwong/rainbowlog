package rainbowlog

import (
	"github.com/ramboll-max/rainbowlog/level"
	"runtime"
	"strconv"
	"time"
)

var _ recordPacker = (*ConsolePacker)(nil)

type ConsolePacker struct {
	RecordPackerForWriter
	consoleColor bool
}

func (j *ConsolePacker) Msg(msg string) {
	var cs []int
	// key color printing
	if j.consoleColor {
		cs = j.record.logger.metaKeys.ConsoleColors(MsgFieldName)
		if cs != nil {
			for _, c := range cs {
				*j.raw = colorStart(*j.raw, c)
			}
		}
	}
	// key & value field print
	j.RecordPackerForWriter.Msg(msg)
	// end key color printing
	if cs != nil {
		for i := 0; i < len(cs); i++ {
			*j.raw = colorEnd(*j.raw)
		}
	}
}

func (j *ConsolePacker) Err(err error) {
	if err == nil || err.Error() == "" {
		return
	}
	var cs []int
	// key color printing
	if j.consoleColor {
		cs = j.record.logger.metaKeys.ConsoleColors(ErrFieldName)
		if cs != nil {
			for _, c := range cs {
				*j.raw = colorStart(*j.raw, c)
			}
		}
	}

	j.RecordPackerForWriter.Err(err)

	// end key color printing
	if cs != nil {
		for i := 0; i < len(cs); i++ {
			*j.raw = colorEnd(*j.raw)
		}
	}
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
			break
		case level.Info:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelInfoLineFieldName)
			break
		case level.Warn:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelWarnLineFieldName)
			break
		case level.Error:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelErrorLineFieldName)
			break
		case level.Panic:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelPanicLineFieldName)
			break
		case level.Fatal:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelFatalLineFieldName)
			break
		case level.Trace:
			cs = j.record.logger.metaKeys.ConsoleColors(metaKeysColorLevelTraceLineFieldName)
			break
		default:
			return 0
		}
		if cs != nil {
			//starter := bytesPool.Get()
			//for _, c := range cs {
			//	*starter = colorStart(*starter, c)
			//}
			//*starter = colorStart(*starter, ColorBlack)
			//*dst = append(*starter, *dst...)
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

func colorStart(dst []byte, color int) []byte {
	s := "\x1b[" + strconv.Itoa(color) + "m"
	return append(dst, s...)
}

func colorEnd(dst []byte) []byte {
	return append(dst, "\x1b[0m"...)
}
