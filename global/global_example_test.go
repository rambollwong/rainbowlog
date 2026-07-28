package global

import (
	"github.com/rambollwong/rainbowlog/level"
	"github.com/rambollwong/rainbowlog/log"
)

func ExampleUseDefault() {
	UseDefault()

	Logger.Info().Msg("Hello world!").Done()
}

func ExampleUseRainbowDefault() {
	UseRainbowDefault()

	Info().Msg("Hello world!").Done()
}

func ExampleUseCustomOptions() {
	UseCustomOptions(
		log.WithLevel(level.Warn),
	)

	Warn().Msg("custom global logger").Done()
}

func Example_convenienceLevels() {
	UseRainbowDefault()

	Debug().Msg("debug message").Done()
	Info().Msg("info message").Done()
	Warn().Msg("warn message").Done()
	Error().Msg("error message").Done()
}
