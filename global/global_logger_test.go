package global

import (
	"testing"

	"github.com/rambollwong/rainbowlog/log"
	"github.com/rambollwong/rainbowlog/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLazyInit(t *testing.T) {
	Logger = nil
	lazyInit()
	require.NotNil(t, Logger, "lazyInit should create a Logger when nil")
	prev := Logger
	lazyInit()
	assert.Equal(t, prev, Logger, "lazyInit should be idempotent")
}

func TestUseDefault(t *testing.T) {
	UseDefault()
	require.NotNil(t, Logger)
}

func TestUseRainbowDefault(t *testing.T) {
	UseRainbowDefault()
	require.NotNil(t, Logger)
}

func TestUseCustomOptions(t *testing.T) {
	UseCustomOptions(log.WithLevel(level.Warn))
	require.NotNil(t, Logger)
}

func TestConvenienceFunctions(t *testing.T) {
	UseDefault()

	for name, fn := range map[string]func() log.Record{
		"Debug": Debug, "Info": Info, "Warn": Warn, "Error": Error, "Record": Record,
	} {
		t.Run(name, func(t *testing.T) {
			rec := fn()
			assert.NotNil(t, rec)
		})
	}
	t.Run("Level", func(t *testing.T) {
		rec := Level(level.Info)
		assert.NotNil(t, rec)
	})
}

func TestLazyInitConvenienceFunction(t *testing.T) {
	Logger = nil
	rec := Info()
	assert.NotNil(t, rec)
	assert.NotNil(t, Logger, "Info() should trigger lazyInit")
}

func TestPanicFunctions(t *testing.T) {
	// These should not be called in tests (they exit/panic),
	// but we can verify they return non-nil records.
	UseDefault()

	t.Run("Fatal", func(t *testing.T) {
		rec := Fatal()
		assert.NotNil(t, rec)
		rec.Discard() // discard to avoid exit
	})

	t.Run("Panic", func(t *testing.T) {
		rec := Panic()
		assert.NotNil(t, rec)
		rec.Discard()
	})
}
