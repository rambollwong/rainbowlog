// Package log provides the definition and initialization operation of a global rainbow logger.
package log

import (
	"os"
	"path/filepath"

	"github.com/rambollwong/rainbowlog"
	"github.com/rambollwong/rainbowlog/level"
)

var (
	Logger *rainbowlog.Logger

	DefaultConfigFilePath = "./"
	DefaultConfigFileName = "rainbowlog.yaml"
)

// UseDefault will initial the global Logger by using default options.
//
//	This is equivalent to using the following options:
//		rainbowlog.WithLevel(level.Debug)
//		rainbowlog.WithLabel("")
//		rainbowlog.WithStack(false)
//		rainbowlog.WithMetaKeys(
//			rainbowlog.MetaTimeFieldName,
//			rainbowlog.MetaLevelFieldName,
//			rainbowlog.MetaLabelFieldName,
//			rainbowlog.MetaCallerFieldName,
//			rainbowlog.MsgFieldName,
//		)
//		rainbowlog.WithConsolePrint(false)
//		rainbowlog.WithLevelFieldMarshalFunc(rainbowlog.GlobalLevelFieldMarshalFunc)
//		rainbowlog.WithCallerMarshalFunc(rainbowlog.GlobalCallerMarshalFunc)
//		rainbowlog.WithErrorMarshalFunc(rainbowlog.GlobalErrorMarshalFunc)
//		rainbowlog.WithErrorStackMarshalFunc(rainbowlog.GlobalErrorStackMarshalFunc)
//		rainbowlog.WithTimeFormat(rainbowlog.GlobalTimeFormat)
//		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stderr)
func UseDefault() {
	Logger = rainbowlog.New(rainbowlog.WithDefault(), rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stderr))
}

// UseRainbowDefault will initial the global Logger by using default options and enable rainbow console printing.
//
//	This is equivalent to using the following options:
//		rainbowlog.WithLevel(level.Debug)
//		rainbowlog.WithLabel("")
//		rainbowlog.WithStack(false)
//		rainbowlog.WithMetaKeys(
//			rainbowlog.MetaTimeFieldName,
//			rainbowlog.MetaLevelFieldName,
//			rainbowlog.MetaLabelFieldName,
//			rainbowlog.MetaCallerFieldName,
//			rainbowlog.MsgFieldName,
//		)
//		rainbowlog.WithConsolePrint(true)
//		rainbowlog.WithRainbowConsole(true)
//		rainbowlog.WithLevelFieldMarshalFunc(rainbowlog.GlobalLevelFieldMarshalFunc)
//		rainbowlog.WithCallerMarshalFunc(rainbowlog.GlobalCallerMarshalFunc)
//		rainbowlog.WithErrorMarshalFunc(rainbowlog.GlobalErrorMarshalFunc)
//		rainbowlog.WithErrorStackMarshalFunc(rainbowlog.GlobalErrorStackMarshalFunc)
//		rainbowlog.WithTimeFormat(rainbowlog.GlobalTimeFormat)
func UseRainbowDefault() {
	Logger = rainbowlog.New(
		rainbowlog.WithDefault(),
		rainbowlog.WithConsolePrint(true),
		rainbowlog.WithRainbowConsole(true),
	)
}

// UseCustomOptions allows user to initial the global Logger by custom options.
func UseCustomOptions(opts ...rainbowlog.Option) {
	Logger = rainbowlog.New(opts...)
}

// UseDefaultConfigFile allows user to initial the global Logger by global config file and custom options.
func UseDefaultConfigFile(opts ...rainbowlog.Option) {
	Logger = rainbowlog.New(append(
		[]rainbowlog.Option{
			rainbowlog.WithDefault(),
			rainbowlog.WithConfigFile(filepath.Join(DefaultConfigFilePath, DefaultConfigFileName)),
		}, opts...)...,
	)
}

// Record create a new rainbowlog.Record with basic.
func Record() rainbowlog.Record {
	return Logger.Record()
}

// Level create a new rainbowlog.Record with the logger level given.
func Level(le level.Level) rainbowlog.Record {
	return Logger.Level(le)
}

// Debug create a new rainbowlog.Record with debug level setting.
func Debug() rainbowlog.Record {
	return Logger.Debug()
}

// Info create a new rainbowlog.Record with info level setting.
func Info() rainbowlog.Record {
	return Logger.Info()
}

// Warn create a new rainbowlog.Record with warn level setting.
func Warn() rainbowlog.Record {
	return Logger.Warn()
}

// Error create a new rainbowlog.Record with error level setting.
func Error() rainbowlog.Record {
	return Logger.Error()
}

// Fatal create a new rainbowlog.Record with fatal level setting.
func Fatal() rainbowlog.Record {
	return Logger.Fatal()
}

// Panic create a new rainbowlog.Record with panic level setting.
func Panic() rainbowlog.Record {
	return Logger.Panic()
}
