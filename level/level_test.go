package level

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		lv  Level
		str string
	}{
		{Debug, DebugStr},
		{Info, InfoStr},
		{Warn, WarnStr},
		{Error, ErrorStr},
		{Fatal, FatalStr},
		{Panic, PanicStr},
		{None, NoneStr},
		{Disabled, DisabledStr},
		{Trace, TraceStr},
		{Level(100), "100"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s level string", test.str), func(t *testing.T) {
			require.Equal(t, test.str, test.lv.String())
		})
	}
}

func TestLevel_KeyFieldValue(t *testing.T) {
	tests := []struct {
		lv            Level
		keyFieldValue string
	}{
		{Debug, "DEB"},
		{Info, "INF"},
		{Warn, "WAR"},
		{Error, "ERR"},
		{Fatal, "FAT"},
		{Panic, "PAN"},
		{None, "UNK"},
		{Disabled, "DIS"},
		{Trace, "TRA"},
		{Level(100), "100"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s level string", test.keyFieldValue), func(t *testing.T) {
			require.Equal(t, test.keyFieldValue, test.lv.KeyFieldValue())
		})
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		lv  Level
		str string
	}{
		{Debug, DebugStr},
		{Info, InfoStr},
		{Warn, WarnStr},
		{Error, ErrorStr},
		{Fatal, FatalStr},
		{Panic, PanicStr},
		{None, NoneStr},
		{Disabled, DisabledStr},
		{Trace, TraceStr},
		{None, "abc"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("level from string %s", test.str), func(t *testing.T) {
			require.Equal(t, test.lv, FromString(test.str))
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		lvStr string
		lv    Level
		e     bool
	}{
		{DebugStr, Debug, false},
		{InfoStr, Info, false},
		{WarnStr, Warn, false},
		{ErrorStr, Error, false},
		{FatalStr, Fatal, false},
		{PanicStr, Panic, false},
		{NoneStr, None, false},
		{DisabledStr, Disabled, false},
		{TraceStr, Trace, false},
		{"100", Level(100), false},
		{"-129", None, true},
		{"128", None, true},
		{"abc", None, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("parse level from string %s", test.lvStr), func(t *testing.T) {
			l, err := ParseLevel(test.lvStr)
			if test.e {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.lv, l)
		})
	}
}
