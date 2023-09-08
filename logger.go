package rainbowlog

import (
	"github.com/RambollWong/rainbowlog/internal/encoder"
	"github.com/RambollWong/rainbowlog/level"
	"os"
)

// WriterEncoderPair wraps a writer and an encoder.
type WriterEncoderPair struct {
	writer LevelWriter
	enc    Encoder
}

// Logger is the rainbow-log structure.
// This is the main entrance of rainbow log.
type Logger struct {
	level                 level.Level
	label                 string
	writerEncoders        []WriterEncoderPair
	hooks                 []Hook
	stack                 bool
	metaKeys              *metaKeys
	consolePrint          bool
	consoleColor          bool
	levelFieldMarshalFunc LevelFieldMarshalFunc
	callerMarshalFunc     CallerMarshalFunc
	errorMarshalFunc      ErrorMarshalFunc
	errorStackMarshalFunc ErrorStackMarshalFunc
	timeFormat            string

	// each Logger instance has an independent *Record pool.
	recordPool *recordPool
}

func createLogger() *Logger {
	return &Logger{
		level:                 level.Debug,
		label:                 "",
		writerEncoders:        nil,
		stack:                 false,
		metaKeys:              emptyMetaKes(),
		consolePrint:          false,
		consoleColor:          false,
		levelFieldMarshalFunc: GlobalLevelFieldMarshalFunc,
		callerMarshalFunc:     GlobalCallerMarshalFunc,
		errorMarshalFunc:      GlobalErrorMarshalFunc,
		errorStackMarshalFunc: GlobalErrorStackMarshalFunc,
		timeFormat:            GlobalTimeFormat,
		recordPool:            nil,
	}
}

// New creates a new *Logger with options optional.
func New(opts ...Option) *Logger {
	logger := createLogger()
	// apply options
	for _, opt := range opts {
		opt(logger)
	}
	// init logger
	logger.initLogger()
	return logger
}

// SubLogger creates a new *Logger that inherit from parent logger.
// Optional options for sub logger also be supported.
func (l *Logger) SubLogger(opts ...Option) *Logger {
	logger := &Logger{
		level:          l.level,
		label:          l.label,
		writerEncoders: l.writerEncoders,
		stack:          l.stack,
		metaKeys:       l.metaKeys.Clone(),
		consolePrint:   l.consolePrint,
		consoleColor:   l.consoleColor,
		recordPool:     nil,
	}
	// apply options
	for _, opt := range opts {
		opt(logger)
	}
	// init logger
	logger.initLogger()
	return logger
}

func (l *Logger) createRecord() *Record {
	r := &Record{
		recordPackers: nil,
		level:         level.Disabled,
		label:         l.label,
		stack:         l.stack,
		useIntDur:     GlobalDurationValueUseInt,
		doneFunc:      nil,
		logger:        l,
	}
	if l.consolePrint {
		r.recordPackers = make([]recordPacker, len(l.writerEncoders)+1)
		r.recordPackers[len(l.writerEncoders)] = &ConsolePacker{
			RecordPackerForWriter: RecordPackerForWriter{
				record:               r,
				meta:                 bytesPool.Get(),
				raw:                  bytesPool.Get(),
				callerSkipFrameCount: 0,
				writerEncoderPair: &WriterEncoderPair{
					enc:    NewTextEncoder(l.metaKeys.Keys()...),
					writer: levelWriterAdapter{Writer: os.Stdout},
				},
			},
			consoleColor: l.consoleColor,
		}
	} else {
		r.recordPackers = make([]recordPacker, len(l.writerEncoders))
	}
	for i, wep := range l.writerEncoders {
		var enc Encoder
		switch tmp := wep.enc.(type) {
		case encoder.TextEncoder, *encoder.TextEncoder:
			enc = NewTextEncoder(l.metaKeys.Keys()...)
		default:
			enc = tmp
		}
		r.recordPackers[i] = &RecordPackerForWriter{
			record:               r,
			meta:                 bytesPool.Get(),
			raw:                  bytesPool.Get(),
			callerSkipFrameCount: 0,
			writerEncoderPair:    &WriterEncoderPair{enc: enc, writer: wep.writer},
		}
	}
	return r
}

func (l *Logger) initLogger() {
	// init record pool
	l.recordPool = newRecordPool(l.createRecord)
}

// Record create a new *Record with basic.
func (l *Logger) Record() *Record {
	if l.level == level.Disabled {
		return nil
	}
	r := l.recordPool.Get()
	r.label = l.label
	return r
}

// Level create a new *Record with the logger level given.
func (l *Logger) Level(le level.Level) *Record {
	if l.level == level.Disabled {
		return nil
	}
	r := l.Record()
	r.level = le
	return r
}

// Debug create a new *Record with debug level setting.
func (l *Logger) Debug() *Record {
	return l.Level(level.Debug)
}

// Info create a new *Record with info level setting.
func (l *Logger) Info() *Record {
	return l.Level(level.Info)
}

// Warn create a new *Record with warn level setting.
func (l *Logger) Warn() *Record {
	return l.Level(level.Warn)
}

// Error create a new *Record with error level setting.
func (l *Logger) Error() *Record {
	return l.Level(level.Error)
}

// Fatal create a new *Record with fatal level setting.
func (l *Logger) Fatal() *Record {
	return l.Level(level.Fatal)
}

// Panic create a new *Record with panic level setting.
func (l *Logger) Panic() *Record {
	return l.Level(level.Panic)
}
