package rainbowlog

import (
	"io"

	"github.com/rambollwong/rainbowcat/util"
	"github.com/rambollwong/rainbowcat/writer/filewriter"
	"github.com/rambollwong/rainbowlog/config"
	"github.com/rambollwong/rainbowlog/level"
)

// Option defines logger options for custom using.
type Option func(logger *Logger)

// WithDefault sets default configurations to logger.
// If you use this option, please put it at the top of all options,
// otherwise it may overwrite the modified configuration as the default.
func WithDefault() Option {
	return func(logger *Logger) {
		logger.level = DefaultLevel
		logger.label = DefaultLabel
		logger.stack = DefaultStack
		logger.metaKeys = defaultMetaKeys()
		logger.consolePrint = DefaultConsolePrint
		logger.consoleColor = DefaultConsoleColor
		logger.levelFieldMarshalFunc = GlobalLevelFieldMarshalFunc
		logger.callerMarshalFunc = GlobalCallerMarshalFunc
		logger.errorMarshalFunc = GlobalErrorMarshalFunc
		logger.errorStackMarshalFunc = GlobalErrorStackMarshalFunc
		logger.timeFormat = GlobalTimeFormat
	}
}

// WithLevel sets logger level.
func WithLevel(lv level.Level) Option {
	return func(logger *Logger) {
		logger.level = lv
	}
}

// WithLabel sets logger label.
func WithLabel(label string) Option {
	return func(logger *Logger) {
		logger.label = label
	}
}

// WithStack enables recording stack or not.
func WithStack(enable bool) Option {
	return func(logger *Logger) {
		logger.stack = enable
	}
}

// WithMetaKeys sets meta key field names for each record.
// Meta keys data will be placed at the front of the record.
func WithMetaKeys(keys ...string) Option {
	return func(logger *Logger) {
		logger.metaKeys.SetKeys(keys)
	}
}

// WithMetaKeyColors sets an ANSI code of color for printing the meta key.
//
// NOTICE: MetaKeyColors only useful for console printing if rainbow console enabled.
// Colors defined in globals.go
func WithMetaKeyColors(key string, colors ...int) Option {
	return func(logger *Logger) {
		logger.metaKeys.SetKeyColors(key, colors)
	}
}

// WithConsolePrint sets whether enable printing to console.
func WithConsolePrint(enable bool) Option {
	return func(logger *Logger) {
		logger.consolePrint = enable
	}
}

// WithRainbowConsole sets whether enable rainbow printing on console.
// This option is useful only when ConsolePrint enabled.
func WithRainbowConsole(enable bool) Option {
	return func(logger *Logger) {
		logger.consoleColor = enable
	}
}

// AppendsEncoderWriters appends writers who use a same encoder to logger.
func AppendsEncoderWriters(encoder Encoder, writers ...io.Writer) Option {
	if encoder == nil || len(writers) == 0 {
		return func(logger *Logger) {}
	}
	wep := WriterEncoderPair{
		enc: encoder,
	}
	if len(writers) > 1 {
		wep.writer = MultiLevelWriter(writers...)
	} else {
		wep.writer = LevelWriterAdapter(writers[0])
	}
	return func(logger *Logger) {
		logger.writerEncoders = append(logger.writerEncoders, wep)
	}
}

// AppendsHooks appends hooks to logger.
func AppendsHooks(hooks ...Hook) Option {
	return func(logger *Logger) {
		logger.hooks = append(logger.hooks, hooks...)
	}
}

// WithLevelFieldMarshalFunc sets the LevelFieldMarshalFunc for logger.
// LevelFieldMarshalFunc will be invoked when printing logs,
// then the result string will be used as the value of level key field.
func WithLevelFieldMarshalFunc(levelFieldMarshalFunc LevelFieldMarshalFunc) Option {
	return func(logger *Logger) {
		logger.levelFieldMarshalFunc = levelFieldMarshalFunc
	}
}

// WithCallerMarshalFunc sets the CallerMarshalFunc for logger.
// If MetaCallerFieldName has been set by WithMetaKeys option
// (if WithMetaKeys is not appended, MetaCallerFieldName is set by default),
// CallerMarshalFunc will be invoked when printing logs.
func WithCallerMarshalFunc(callerMarshalFunc CallerMarshalFunc) Option {
	return func(logger *Logger) {
		logger.callerMarshalFunc = callerMarshalFunc
	}
}

// WithErrorMarshalFunc sets the ErrorMarshalFunc for logger.
// If record.Err() has been invoked, ErrorMarshalFunc will be invoked when printing logs.
// The marshal result string will be used as the value of error key field.
func WithErrorMarshalFunc(errorMarshalFunc ErrorMarshalFunc) Option {
	return func(logger *Logger) {
		logger.errorMarshalFunc = errorMarshalFunc
	}
}

// WithErrorStackMarshalFunc sets the ErrorStackMarshalFunc for logger.
func WithErrorStackMarshalFunc(errorStackMarshalFunc ErrorStackMarshalFunc) Option {
	return func(logger *Logger) {
		logger.errorStackMarshalFunc = errorStackMarshalFunc
	}
}

// WithTimeFormat sets the format of time for logger printing.
func WithTimeFormat(timeFormat string) Option {
	return func(logger *Logger) {
		logger.timeFormat = timeFormat
	}
}

// WithConfig sets the Logger's properties according to the provided configuration parameters.
// If the Enable field in the configuration is false, the logger level will be set Disabled.
// Normally, the WithDefault() option should be set before calling this option.
func WithConfig(config config.LoggerConfig) Option {
	if !config.Enable {
		return WithLevel(level.Disabled)
	}
	return func(logger *Logger) {
		lv := level.FromString(config.Level)
		if lv == level.None {
			panic("wrong logger level: " + config.Level)
		}
		logger.level = lv
		logger.label = config.Label
		logger.stack = config.Stack
		logger.consolePrint = config.EnableConsolePrinting
		logger.consoleColor = config.EnableRainbowConsole
		if config.TimeFormat != "" {
			logger.timeFormat = config.TimeFormat
		}
		if config.SizeRollingFileConfig.Enable {
			fileSizeLimit, err := util.ParseToBytesSize(config.SizeRollingFileConfig.FileSizeLimit, 1024)
			if err != nil {
				panic("wrong file size limit: " + config.SizeRollingFileConfig.FileSizeLimit)
			}
			writer, err := filewriter.NewSizeRollingFileWriter(
				config.SizeRollingFileConfig.LogFilePath,
				config.SizeRollingFileConfig.LogFileBaseName,
				config.SizeRollingFileConfig.MaxBackups,
				fileSizeLimit,
			)
			if err != nil {
				panic("error new size rolling file writer: " + err.Error())
			}
			encoder := GlobalEncoderParseFunc(config.SizeRollingFileConfig.Encoder)
			logger.writerEncoders = append(logger.writerEncoders, WriterEncoderPair{
				writer: LevelWriterAdapter(writer),
				enc:    encoder,
			})
		}
		if config.TimeRollingFileConfig.Enable {
			writer, err := filewriter.NewTimeRollingFileWriter(
				config.TimeRollingFileConfig.LogFilePath,
				config.TimeRollingFileConfig.LogFileBaseName,
				config.TimeRollingFileConfig.MaxBackups,
				config.TimeRollingFileConfig.RollingPeriod,
			)
			if err != nil {
				panic("error new time rolling file writer: " + err.Error())
			}
			encoder := GlobalEncoderParseFunc(config.SizeRollingFileConfig.Encoder)
			logger.writerEncoders = append(logger.writerEncoders, WriterEncoderPair{
				writer: LevelWriterAdapter(writer),
				enc:    encoder,
			})
		}
	}
}

// WithConfigFile loads the log configuration from the specified configuration file,
// sets the properties of the Logger according to the configuration parameters.
// If an error occurs while loading the configuration file, trigger a panic.
func WithConfigFile(configFile string) Option {
	cfg, err := config.LoadLoggerConfigFromFile(configFile)
	if err != nil {
		panic("error load logger config from file: " + configFile)
	}
	return WithConfig(*cfg)
}
