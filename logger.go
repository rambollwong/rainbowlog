// Package rainbowlog implements a comprehensive logging solution with multiple
// writer implementations including buffered, synchronized, and multi-writer
// capabilities. It supports structured logging with different encoders and
// offers flexible configuration options for various logging needs.
package rainbowlog

import (
	"os"

	"github.com/rambollwong/rainbowlog/internal/encoder"
	"github.com/rambollwong/rainbowlog/level"
	"golang.org/x/sync/errgroup"
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

	// pre-built encoders to avoid per-record creation in the pool's New function.
	// Index-aligned with writerEncoders; nil for writers that use a non-text encoder.
	preBuiltEncoders       []Encoder
	preBuiltConsoleEncoder Encoder
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
	weps := make([]WriterEncoderPair, len(l.writerEncoders))
	copy(weps, l.writerEncoders)
	logger := &Logger{
		level:                 l.level,
		label:                 l.label,
		writerEncoders:        weps,
		stack:                 l.stack,
		metaKeys:              l.metaKeys.Clone(),
		consolePrint:          l.consolePrint,
		consoleColor:          l.consoleColor,
		levelFieldMarshalFunc: l.levelFieldMarshalFunc,
		callerMarshalFunc:     l.callerMarshalFunc,
		errorMarshalFunc:      l.errorMarshalFunc,
		errorStackMarshalFunc: l.errorStackMarshalFunc,
		timeFormat:            l.timeFormat,
		recordPool:            nil,
	}
	// apply options
	for _, opt := range opts {
		opt(logger)
	}
	// init logger
	logger.initLogger()
	return logger
}

func (l *Logger) createRecord() Record {
	r := &LogRecord{
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
					enc:    l.preBuiltConsoleEncoder,
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
		if l.preBuiltEncoders[i] != nil {
			enc = l.preBuiltEncoders[i]
		} else {
			enc = wep.enc
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
	// Pre-build text encoders for writer-encoder pairs so that createRecord (the
	// pool's New function) can reference them without allocating per pool-object.
	l.preBuiltEncoders = make([]Encoder, len(l.writerEncoders))
	for i, wep := range l.writerEncoders {
		switch wep.enc.(type) {
		case encoder.TextEncoder, *encoder.TextEncoder:
			l.preBuiltEncoders[i] = NewTextEncoder(l.metaKeys.Keys()...)
		}
	}
	if l.consolePrint {
		l.preBuiltConsoleEncoder = NewTextEncoder(l.metaKeys.Keys()...)
	}

	// init Record pool
	l.recordPool = newRecordPool(l.createRecord)
}

type flushI interface{ Flush() error }

// Flush forces any buffered data to be written out to all writers that implement the Flush() method.
// It iterates through all the writer-encoder pairs and checks if the writer supports flushing.
// For multiLevelWriter, it recursively flushes all underlying writers that support flushing.
// Returns the first error encountered during flushing, or nil if all flush operations succeed.
func (l *Logger) Flush() (err error) {
	eg := errgroup.Group{}
	for _, writerEncoder := range l.writerEncoders {
		switch w := writerEncoder.writer.(type) {
		case levelWriterAdapter:
			fw, ok := w.Writer.(flushI)
			if ok {
				eg.Go(fw.Flush)
			}
		case *levelWriterAdapter:
			fw, ok := w.Writer.(flushI)
			if ok {
				eg.Go(fw.Flush)
			}
		case *syncWriter:
			// Skip writers that don't need explicit flushing
		case flushI:
			// Flush writers that directly implement flushI
			eg.Go(w.Flush)
		case multiLevelWriter:
			// For multiLevelWriter, flush all underlying writers that support flushing
			for _, writer := range w.writers {
				if fw, ok := writer.(flushI); ok {
					eg.Go(fw.Flush)
				}
			}
		default:
			// Skip other writers
		}
	}
	return eg.Wait()
}

// Record create a new Record with basic.
func (l *Logger) Record() Record {
	if l.level == level.Disabled {
		return nilRecord
	}
	r := l.recordPool.Get()
	r.WithLabels(l.label)
	return r
}

// Level create a new Record with the logger level given.
func (l *Logger) Level(le level.Level) Record {
	if l.level == level.Disabled {
		return nilRecord
	}
	r := l.Record()
	r.WithLevel(le)
	return r
}

// Debug create a new Record with debug level setting.
func (l *Logger) Debug() Record {
	return l.Level(level.Debug)
}

// Info create a new Record with info level setting.
func (l *Logger) Info() Record {
	return l.Level(level.Info)
}

// Warn create a new Record with warn level setting.
func (l *Logger) Warn() Record {
	return l.Level(level.Warn)
}

// Error create a new Record with error level setting.
func (l *Logger) Error() Record {
	return l.Level(level.Error)
}

// Fatal create a new Record with fatal level setting.
func (l *Logger) Fatal() Record {
	return l.Level(level.Fatal)
}

// Panic create a new Record with panic level setting.
func (l *Logger) Panic() Record {
	return l.Level(level.Panic)
}
