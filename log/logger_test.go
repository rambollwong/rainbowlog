package log

import (
	"testing"

	"github.com/rambollwong/rainbowlog/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================== Bug regression tests ==============================

// Regression: pooled record must retain stack flag after Recycle→Get cycle
func TestRegression_PooledRecordRetainsStackFlag(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithStack(true),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	// First Get from pool (createRecord)
	rec1 := logger.Record().(*LogRecord)
	assert.True(t, rec1.stack, "first record from pool should have stack=true")
	assert.False(t, rec1.useIntDur, "default useIntDur should be false")

	// Simulate Done → Put → Reset
	logger.recordPool.Put(rec1)

	// Second Get from pool (reused)
	rec2 := logger.Record().(*LogRecord)
	assert.True(t, rec2.stack, "reused record lost stack flag — Record() must re-apply logger.stack after Reset")
}

// Regression: errorMarshalFunc nil fallback
func TestRegression_ErrorMarshalFuncNilFallback(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	logger.errorMarshalFunc = nil

	// Should not panic — defaults to GlobalErrorMarshalFunc
	assert.NotPanics(t, func() {
		logger.Error().Err(assert.AnError).Msg("test").Done()
	}, "nil errorMarshalFunc should fall back to GlobalErrorMarshalFunc, not panic")
}

// Regression: useIntDur flag applied on reuse
func TestRegression_UseIntDurAppliedOnReuse(t *testing.T) {
	GlobalDurationValueUseInt = true
	defer func() { GlobalDurationValueUseInt = false }()

	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	rec1 := logger.Record().(*LogRecord)
	assert.True(t, rec1.useIntDur, "first record useIntDur should be true")

	logger.recordPool.Put(rec1)

	rec2 := logger.Record().(*LogRecord)
	assert.True(t, rec2.useIntDur, "reused record lost useIntDur flag")
}

// ============================== LogRecord tests ==============================

func TestLogRecord_Msg(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	rec := logger.Info().Msg("hello world")
	assert.NotNil(t, rec)
	// Msg should not panic
}

func TestLogRecord_Err(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	assert.NotPanics(t, func() {
		logger.Error().Err(assert.AnError).Msg("test").Done()
	})
}

func TestLogRecord_Str(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	assert.NotPanics(t, func() {
		logger.Info().Str("key", "value").Msg("test").Done()
	})
}

func TestLogRecord_Int(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	assert.NotPanics(t, func() {
		logger.Info().Int("count", 42).Msg("test").Done()
	})
}

func TestLogRecord_Done(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	// Done should return record to pool
	assert.NotPanics(t, func() {
		logger.Info().Str("a", "b").Msg("done").Done()
	})
}

func TestLogRecord_Discard(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	rec := logger.Info().Discard()
	assert.Equal(t, level.Disabled, rec.(*LogRecord).level)
}

func TestLogRecord_Msgf(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	assert.NotPanics(t, func() {
		logger.Info().Msgf("hello %s", "world").Done()
	})
}

func TestLogRecord_WithCallerSkip(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	assert.NotPanics(t, func() {
		logger.Info().WithCallerSkip(1).Msg("test").Done()
	})
}

func TestLogRecord_WithLabels(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	assert.NotPanics(t, func() {
		logger.Info().WithLabels("label1", "label2").Msg("test").Done()
	})
}

func TestLogRecord_UseIntDur(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	rec := logger.Info().UseIntDur()
	assert.NotNil(t, rec)
	assert.True(t, rec.(*LogRecord).useIntDur)
}

func TestLogRecord_FatalAtDisabledLevel(t *testing.T) {
	// When logger is Disabled, Fatal should be suppressed
	lDisabled := &Logger{level: level.Disabled}
	r := &LogRecord{level: level.Fatal, logger: lDisabled}
	// Disabled means no logging and no Fatal/Panic effects
	assert.True(t, r.strikeOrNot(), "Fatal should be struck when logger is Disabled")
}

// ============================== Logger tests ==============================

func TestLogger_New(t *testing.T) {
	logger := New(
		WithLevel(level.Info),
		WithLabels("test"),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	require.NotNil(t, logger)
	assert.Equal(t, level.Info, logger.level)
	assert.Equal(t, "test", logger.label)
}

func TestLogger_LevelMethods(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)

	tests := []struct {
		name string
		fn   func() Record
		lv   level.Level
	}{
		{"Debug", logger.Debug, level.Debug},
		{"Info", logger.Info, level.Info},
		{"Warn", logger.Warn, level.Warn},
		{"Error", logger.Error, level.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := tt.fn().(*LogRecord)
			assert.Equal(t, tt.lv, rec.level)
		})
	}
}

func TestLogger_LevelFiltered(t *testing.T) {
	logger := New(
		WithLevel(level.Warn), // only Warn and above
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	// Debug should be suppressed
	rec := logger.Debug().(*LogRecord)
	assert.True(t, rec.strikeOrNot(), "Debug below Warn should be suppressed")
}

func TestLogger_SubLogger(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithLabels("parent"),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	sub := logger.SubLogger(WithLabels("child"))

	assert.Equal(t, "child", sub.label)
	assert.Equal(t, level.Debug, sub.level)
	// Sub should have independent writerEncoders
	assert.Equal(t, len(logger.writerEncoders), len(sub.writerEncoders))
}

func TestLogger_Flush(t *testing.T) {
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	// Flush should not panic on nullWriter
	assert.NotPanics(t, func() {
		err := logger.Flush()
		assert.NoError(t, err)
	})
}

func TestLogger_GlobalDurationValueUseInt(t *testing.T) {
	original := GlobalDurationValueUseInt
	defer func() { GlobalDurationValueUseInt = original }()

	GlobalDurationValueUseInt = true
	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	rec := logger.Info().(*LogRecord)
	assert.True(t, rec.useIntDur)

	GlobalDurationValueUseInt = false
	logger2 := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
	)
	rec2 := logger2.Info().(*LogRecord)
	assert.False(t, rec2.useIntDur)
}

// ============================== LevelWriter / Writer tests ==============================

func TestLevelWriterAdapter_WrapsNonLevelWriter(t *testing.T) {
	buf := &nullWriter{}
	lw := LevelWriterAdapter(buf)
	assert.NotNil(t, lw)
	n, err := lw.Write([]byte("test"))
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
}

func TestLevelWriterAdapter_PassthroughLevelWriter(t *testing.T) {
	// Create a LevelWriter and pass it through
	original := &nullLevelWriter{}
	lw := LevelWriterAdapter(original)
	assert.Equal(t, original, lw) // should return as-is
}

type nullLevelWriter struct{ nullWriter }

func (n *nullLevelWriter) WriteLevel(_ level.Level, b []byte) (int, error) { return n.Write(b) }

// ============================== Hook tests ==============================

func TestHook_HookFunc(t *testing.T) {
	called := false
	hook := HookFunc(func(r Record, lv level.Level, message string) {
		called = true
	})

	hook.RunHook(&LogRecord{}, level.Info, "test")
	assert.True(t, called)
}

func TestLogger_AppendsHooks(t *testing.T) {
	hookCalled := false
	hook := HookFunc(func(r Record, lv level.Level, message string) {
		hookCalled = true
	})

	logger := New(
		WithLevel(level.Debug),
		WithConsolePrint(false),
		AppendsEncoderWriters(JsonEnc, &nullWriter{}),
		AppendsHooks(hook),
	)

	logger.Info().Msg("test").Done()
	assert.True(t, hookCalled, "hook should be called on Info message")
}

// nullWriter is a writer that discards all data.
type nullWriter struct{}

func (n *nullWriter) Write(b []byte) (int, error)           { return len(b), nil }
func (n *nullWriter) WriteLevel(_ level.Level, b []byte) (int, error) { return len(b), nil }
