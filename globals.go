package rainbowlog

import (
	"fmt"
	"github.com/rambollwong/rainbowlog/level"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	LevelFieldMarshalFunc func(level.Level) string
	CallerMarshalFunc     func(file string, line int) string
	ErrorMarshalFunc      func(err error) string
	ErrorStackMarshalFunc func(err error) interface{}
)

///////////////////
/// Time Format ///
///////////////////

const (
	// TimeFormatUnix defines a time format that makes time fields to be
	// serialized as Unix timestamp integers.
	TimeFormatUnix = ""

	// TimeFormatUnixMs defines a time format that makes time fields to be
	// serialized as Unix timestamp integers in milliseconds.
	TimeFormatUnixMs = "UNIXMS"

	// TimeFormatUnixMicro defines a time format that makes time fields to be
	// serialized as Unix timestamp integers in microseconds.
	TimeFormatUnixMicro = "UNIXMICRO"

	// TimeFormatUnixNano defines a time format that makes time fields to be
	// serialized as Unix timestamp integers in nanoseconds.
	TimeFormatUnixNano = "UNIXNANO"
)

/////////////////////
/// Global params ///
/////////////////////

var (
	// MsgFieldName is the field name for core message.
	MsgFieldName = "message"
	// ErrFieldName is the field name for err.
	ErrFieldName = "error"
	// ErrStackFieldName is the field name for err stack.
	ErrStackFieldName = "stack"

	MetaTimeFieldName   = "META_TIME"
	MetaCallerFieldName = "META_CALLER"
	MetaLevelFieldName  = "META_LEVEL"
	MetaLabelFieldName  = "META_LABEL"

	// extra meta keys field name for console printing color.

	metaEndFieldName                     = "META_END"
	metaKeysColorLevelDebugFieldName     = MetaLevelFieldName + "_DEBUG"
	metaKeysColorLevelInfoFieldName      = MetaLevelFieldName + "_INFO"
	metaKeysColorLevelWarnFieldName      = MetaLevelFieldName + "_WARN"
	metaKeysColorLevelErrorFieldName     = MetaLevelFieldName + "_ERROR"
	metaKeysColorLevelFatalFieldName     = MetaLevelFieldName + "_FATAL"
	metaKeysColorLevelPanicFieldName     = MetaLevelFieldName + "_PANIC"
	metaKeysColorLevelTraceFieldName     = MetaLevelFieldName + "_TRACE"
	metaKeysColorLevelDebugLineFieldName = MetaLevelFieldName + "_DEBUG_LINE"
	metaKeysColorLevelInfoLineFieldName  = MetaLevelFieldName + "_INFO_LINE"
	metaKeysColorLevelWarnLineFieldName  = MetaLevelFieldName + "_WARN_LINE"
	metaKeysColorLevelErrorLineFieldName = MetaLevelFieldName + "_ERROR_LINE"
	metaKeysColorLevelFatalLineFieldName = MetaLevelFieldName + "_FATAL_LINE"
	metaKeysColorLevelPanicLineFieldName = MetaLevelFieldName + "_PANIC_LINE"
	metaKeysColorLevelTraceLineFieldName = MetaLevelFieldName + "_TRACE_LINE"

	GlobalDurationValueUseInt = false

	DefaultStack = false

	// ErrorHandler will be called whenever some error threw when logger working.
	// If not set, the error will be printed on the stderr.
	// This handler must be thread safe and non-blocking.
	ErrorHandler = func(err error) {
		fmt.Fprintf(os.Stderr, "rainbowlog: error found: %v\n", err)
	}

	GlobalTimeFormat = "2006-01-02 15:04:05.000"

	GlobalLevelFieldMarshalFunc LevelFieldMarshalFunc = func(l level.Level) string {
		return strings.ToUpper(l.String())
	}

	innerCallerSkipFrameCount = 3

	// CallerSkipFrameCount is the number of stack frames to skip to find the caller.
	CallerSkipFrameCount = 0

	// GlobalCallerMarshalFunc allows customization of global caller marshaling.
	GlobalCallerMarshalFunc CallerMarshalFunc = func(file string, line int) string {
		return file + ":" + strconv.Itoa(line)
	}

	// GlobalErrorMarshalFunc allows customization of global error marshaling.
	GlobalErrorMarshalFunc ErrorMarshalFunc = func(err error) string {
		if err == nil {
			return ""
		}
		return err.Error()
	}

	// GlobalErrorStackMarshalFunc extract the stack from err if any.
	GlobalErrorStackMarshalFunc ErrorStackMarshalFunc

	nowFunc = func() time.Time {
		return time.Now()
	}
)

//////////////
/// Colors ///
//////////////

const (
	ColorBold          = 1
	ColorUnderline     = 4
	ColorStrikeThrough = 9
	ColorUnderlineBold = 21
)

const (
	ColorBlack = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorPale
)

const (
	ColorBgBlack = iota + 40
	ColorBgRed
	ColorBgGreen
	ColorBgYellow
	ColorBgBlue
	ColorBgMagenta
	ColorBgCyan
	ColorBgPale
)

const (
	ColorHlBlack = iota + 90
	ColorHlRed
	ColorHlGreen
	ColorHlYellow
	ColorHlBlue
	ColorHlMagenta
	ColorHlCyan
	ColorHlPale
)

const (
	ColorHlBgBlack = iota + 100
	ColorHlBgRed
	ColorHlBgGreen
	ColorHlBgYellow
	ColorHlBgBlue
	ColorHlBgMagenta
	ColorHlBgCyan
	ColorHlBgPale
)

/////////////////////////////
/// Default logger params ///
/////////////////////////////

const (
	DefaultLevel        = level.Debug
	DefaultLabel        = ""
	DefaultConsolePrint = false
	DefaultConsoleColor = true
)

func defaultMetaKeys() *metaKeys {
	return &metaKeys{
		keys: []string{MetaTimeFieldName, MetaLevelFieldName, MetaLabelFieldName, MetaCallerFieldName, MsgFieldName},
		consoleColors: map[string][]int{
			MetaTimeFieldName:   {ColorPale},
			MetaLevelFieldName:  {},
			MetaLabelFieldName:  {ColorHlMagenta},
			MetaCallerFieldName: {ColorCyan},
			MsgFieldName:        {ColorBlue},
			ErrFieldName:        {ColorRed},

			metaKeysColorLevelDebugFieldName: {ColorMagenta},
			metaKeysColorLevelInfoFieldName:  {ColorGreen},
			metaKeysColorLevelWarnFieldName:  {ColorYellow},
			metaKeysColorLevelErrorFieldName: {ColorHlRed, ColorBold},
			metaKeysColorLevelFatalFieldName: {ColorHlRed, ColorBold},
			metaKeysColorLevelPanicFieldName: {ColorHlRed, ColorBold},
			metaKeysColorLevelTraceFieldName: {ColorHlMagenta},

			metaKeysColorLevelDebugLineFieldName: {ColorBgMagenta},
			metaKeysColorLevelInfoLineFieldName:  {ColorBgGreen},
			metaKeysColorLevelWarnLineFieldName:  {ColorBgYellow},
			metaKeysColorLevelErrorLineFieldName: {ColorHlBgRed},
			metaKeysColorLevelFatalLineFieldName: {ColorHlBgRed},
			metaKeysColorLevelPanicLineFieldName: {ColorHlBgRed},
			metaKeysColorLevelTraceLineFieldName: {ColorHlBgMagenta},

			metaEndFieldName: {ColorHlCyan},
		},
	}
}
