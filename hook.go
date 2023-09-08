package rainbowlog

import "github.com/RambollWong/rainbowlog/level"

// Hook defines an interface to a record hook.
type Hook interface {
	// RunHook runs the hook with the event.
	RunHook(r *Record, level level.Level, message string)
}

// HookFunc is an adaptor to allow the use of an ordinary function as a Hook.
type HookFunc func(r *Record, level level.Level, message string)

// RunHook implements the Hook interface.
func (hf HookFunc) RunHook(r *Record, level level.Level, message string) {
	hf(r, level, message)
}

// LevelHook applies a different hook for each level.
type LevelHook struct {
	NoneLevelHook, TraceHook, DebugHook, InfoHook, WarnHook, ErrorHook, FatalHook, PanicHook Hook
}

// RunHook implements the Hook interface.
func (lh LevelHook) RunHook(r *Record, lv level.Level, message string) {
	switch lv {
	case level.Debug:
		if lh.DebugHook != nil {
			lh.DebugHook.RunHook(r, lv, message)
		}
	case level.Info:
		if lh.InfoHook != nil {
			lh.InfoHook.RunHook(r, lv, message)
		}
	case level.Warn:
		if lh.WarnHook != nil {
			lh.WarnHook.RunHook(r, lv, message)
		}
	case level.Error:
		if lh.ErrorHook != nil {
			lh.ErrorHook.RunHook(r, lv, message)
		}
	case level.Trace:
		if lh.TraceHook != nil {
			lh.TraceHook.RunHook(r, lv, message)
		}
	case level.Fatal:
		if lh.FatalHook != nil {
			lh.FatalHook.RunHook(r, lv, message)
		}
	case level.Panic:
		if lh.PanicHook != nil {
			lh.PanicHook.RunHook(r, lv, message)
		}
	case level.None:
		if lh.NoneLevelHook != nil {
			lh.NoneLevelHook.RunHook(r, lv, message)
		}
	}
}

func NewLevelHook() LevelHook {
	return LevelHook{}
}
