package rainbowlog

import (
	"io"
	"strings"

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

// WithLabels sets the labels for the Logger using one or more strings.
// The labels are joined using the "|" separator and assigned to the Logger's label property.
func WithLabels(labels ...string) Option {
	return func(logger *Logger) {
		logger.label = strings.Join(labels, "|")
	}
}

// WithStack enables recording stack or not.
func WithStack(enable bool) Option {
	return func(logger *Logger) {
		logger.stack = enable
	}
}

// WithMetaKeys sets meta key field names for each Record.
// Meta keys data will be placed at the front of the Record.
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
// If Record.Err() has been invoked, ErrorMarshalFunc will be invoked when printing logs.
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
			src := config.SizeRollingFileConfig
			fileSizeLimit, err := util.ParseToBytesSize(src.FileSizeLimit, 1024)
			if err != nil {
				panic("wrong file size limit: " + src.FileSizeLimit)
			}
			writer, err := filewriter.NewSizeRollingFileWriter(
				src.LogFilePath,
				src.LogFileBaseName,
				src.MaxBackups,
				fileSizeLimit,
			)
			if err != nil {
				panic("error new size rolling file writer: " + err.Error())
			}
			encoder := GlobalEncoderParseFunc(src.Encoder)
			var w LevelWriter
			if src.UseBufferedWriter {
				bufferSize, err := util.ParseToBytesSize(src.WriterBufferSize, 1024)
				if err != nil {
					panic("wrong buffer size: " + src.WriterBufferSize)
				}
				w = BufferedLevelWriter(writer, int(bufferSize))
			} else {
				w = LevelWriterAdapter(writer)
			}
			logger.writerEncoders = append(logger.writerEncoders, WriterEncoderPair{
				writer: w,
				enc:    encoder,
			})
		}
		if config.TimeRollingFileConfig.Enable {
			trc := config.TimeRollingFileConfig
			writer, err := filewriter.NewTimeRollingFileWriter(
				trc.LogFilePath,
				trc.LogFileBaseName,
				trc.MaxBackups,
				trc.RollingPeriod,
			)
			if err != nil {
				panic("error new time rolling file writer: " + err.Error())
			}
			encoder := GlobalEncoderParseFunc(trc.Encoder)
			var w LevelWriter
			if trc.UseBufferedWriter {
				bufferSize, err := util.ParseToBytesSize(trc.WriterBufferSize, 1024)
				if err != nil {
					panic("wrong buffer size: " + trc.WriterBufferSize)
				}
				w = BufferedLevelWriter(writer, int(bufferSize))
			} else {
				w = LevelWriterAdapter(writer)
			}
			logger.writerEncoders = append(logger.writerEncoders, WriterEncoderPair{
				writer: w,
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
