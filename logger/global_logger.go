// Package logger provides the definition and initialization operation of a global rainbow logger.
package logger

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/rambollwong/rainbowlog/level"
	"github.com/rambollwong/rainbowlog/log"
)

var (
	mu     sync.Mutex
	Logger *log.Logger

	DefaultConfigFilePath = "./"
	DefaultConfigFileName = "log.yaml"
)

func lazyInit() {
	if Logger == nil {
		mu.Lock()
		defer mu.Unlock()
		if Logger == nil {
			UseRainbowDefault()
		}
	}
}

// UseDefault will initial the global Logger by using default options.
//
//	This is equivalent to using the following options:
//		log.WithLevel(level.Debug)
//		log.WithLabel("")
//		log.WithStack(false)
//		log.WithMetaKeys(
//			log.MetaTimeFieldName,
//			log.MetaLevelFieldName,
//			log.MetaLabelFieldName,
//			log.MetaCallerFieldName,
//			log.MsgFieldName,
//		)
//		log.WithConsolePrint(false)
//		log.WithLevelFieldMarshalFunc(log.GlobalLevelFieldMarshalFunc)
//		log.WithCallerMarshalFunc(log.GlobalCallerMarshalFunc)
//		log.WithErrorMarshalFunc(log.GlobalErrorMarshalFunc)
//		log.WithErrorStackMarshalFunc(log.GlobalErrorStackMarshalFunc)
//		log.WithTimeFormat(log.GlobalTimeFormat)
//		log.AppendsEncoderWriters(log.JsonEnc, os.Stderr)
func UseDefault() {
	Logger = log.New(log.WithDefault(), log.AppendsEncoderWriters(log.JsonEnc, os.Stderr))
}

// UseRainbowDefault will initial the global Logger by using default options and enable rainbow console printing.
//
//	This is equivalent to using the following options:
//		log.WithLevel(level.Debug)
//		log.WithLabel("")
//		log.WithStack(false)
//		log.WithMetaKeys(
//			log.MetaTimeFieldName,
//			log.MetaLevelFieldName,
//			log.MetaLabelFieldName,
//			log.MetaCallerFieldName,
//			log.MsgFieldName,
//		)
//		log.WithConsolePrint(true)
//		log.WithRainbowConsole(true)
//		log.WithLevelFieldMarshalFunc(log.GlobalLevelFieldMarshalFunc)
//		log.WithCallerMarshalFunc(log.GlobalCallerMarshalFunc)
//		log.WithErrorMarshalFunc(log.GlobalErrorMarshalFunc)
//		log.WithErrorStackMarshalFunc(log.GlobalErrorStackMarshalFunc)
//		log.WithTimeFormat(log.GlobalTimeFormat)
func UseRainbowDefault() {
	Logger = log.New(
		log.WithDefault(),
		log.WithConsolePrint(true),
		log.WithRainbowConsole(true),
	)
}

// UseCustomOptions allows user to initial the global Logger by custom options.
func UseCustomOptions(opts ...log.Option) {
	Logger = log.New(opts...)
}

// UseDefaultConfigFile allows user to initial the global Logger by global config file and custom options.
func UseDefaultConfigFile(opts ...log.Option) {
	Logger = log.New(append(
		[]log.Option{
			log.WithDefault(),
			log.WithConfigFile(filepath.Join(DefaultConfigFilePath, DefaultConfigFileName)),
		}, opts...)...,
	)
}

// Record create a new log.Record with basic.
func Record() log.Record {
	lazyInit()
	return Logger.Record()
}

// Level create a new log.Record with the logger level given.
func Level(le level.Level) log.Record {
	lazyInit()
	return Logger.Level(le)
}

// Debug create a new log.Record with debug level setting.
func Debug() log.Record {
	lazyInit()
	return Logger.Debug()
}

// Info create a new log.Record with info level setting.
func Info() log.Record {
	lazyInit()
	return Logger.Info()
}

// Warn create a new log.Record with warn level setting.
func Warn() log.Record {
	lazyInit()
	return Logger.Warn()
}

// Error create a new log.Record with error level setting.
func Error() log.Record {
	lazyInit()
	return Logger.Error()
}

// Fatal create a new log.Record with fatal level setting.
func Fatal() log.Record {
	lazyInit()
	return Logger.Fatal()
}

// Panic create a new log.Record with panic level setting.
func Panic() log.Record {
	lazyInit()
	return Logger.Panic()
}
