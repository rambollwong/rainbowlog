package level

import (
	"fmt"
	"strconv"
	"strings"
)

// Level of log.
type Level int8

const (
	// Debug defines debug log level.
	Debug Level = iota
	// Info defines info log level.
	Info
	// Warn defines warn log level.
	Warn
	// Error defines error log level.
	Error
	// Fatal defines fatal log level.
	// If Fatal level set, os.Exit(1) will be invoked when record.Done() is called.
	Fatal
	// Panic defines panic log level.
	// If Panic level set, logger will panic out when record.Done() is called.
	Panic
	// None defines unknown log level.
	None
	// Disabled disables the logger.
	Disabled

	// Trace defines trace log level.
	Trace Level = -1

	// DisabledStr is the value used for trace-level field.
	DisabledStr = "disabled"
	// DebugStr is the value used for debug-level field.
	DebugStr = "debug"
	// InfoStr is the value used for info-level field.
	InfoStr = "info"
	// WarnStr is the value used for warn-level field.
	WarnStr = "warn"
	// ErrorStr is the value used for error-level field.
	ErrorStr = "error"
	// FatalStr is the value used for fatal-level field.
	FatalStr = "fatal"
	// PanicStr is the value used for panic-level field.
	PanicStr = "panic"
	// TraceStr is the value used for trace-level field.
	TraceStr = "trace"
	// NoneStr is the value used for unknown level field.
	NoneStr = "unknown"
)

var (
	_lsM = map[Level]string{
		Disabled: DisabledStr,
		Debug:    DebugStr,
		Info:     InfoStr,
		Warn:     WarnStr,
		Error:    ErrorStr,
		Fatal:    FatalStr,
		Panic:    PanicStr,
		Trace:    TraceStr,
		None:     NoneStr,
	}
	_lkvM = map[Level]string{
		Disabled: "DIS",
		Debug:    "DEB",
		Info:     "INF",
		Warn:     "WAR",
		Error:    "ERR",
		Fatal:    "FAT",
		Panic:    "PAN",
		Trace:    "TRA",
		None:     "UNK",
	}
	_slM = map[string]Level{
		DisabledStr: Disabled,
		DebugStr:    Debug,
		InfoStr:     Info,
		WarnStr:     Warn,
		ErrorStr:    Error,
		FatalStr:    Fatal,
		PanicStr:    Panic,
		TraceStr:    Trace,
		NoneStr:     None,
	}
)

func (l Level) String() string {
	str, ok := _lsM[l]
	if ok {
		return str
	}
	return strconv.Itoa(int(l))
}

func (l Level) KeyFieldValue() string {
	str, ok := _lkvM[l]
	if ok {
		return str
	}
	return strconv.Itoa(int(l))
}

func FromString(levelStr string) Level {
	l, ok := _slM[strings.ToLower(levelStr)]
	if ok {
		return l
	}
	return None
}

// ParseLevel converts a level string into a rainbowlog Level.
// If the string given does not match known values, return an error.
func ParseLevel(levelStr string) (Level, error) {
	l := FromString(levelStr)
	if l == None && levelStr != NoneStr {
		i, err := strconv.Atoi(levelStr)
		if err != nil {
			return None, fmt.Errorf("unknown level string: %s, defaulting to None level", levelStr)
		}
		if i > 127 || i < -128 {
			return None, fmt.Errorf("out-of-bounds level: '%d', defaulting to None level", i)
		}
		return Level(i), nil
	}
	return l, nil
}
