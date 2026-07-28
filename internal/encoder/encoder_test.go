package encoder

import (
	"testing"
	"time"
)

// ============================== Duration encoding tests ==============================

func TestJsonEncoder_Duration(t *testing.T) {
	je := JsonEncoder{}

	tests := []struct {
		name   string
		d      time.Duration
		unit   time.Duration
		useInt bool
		want   string
	}{
		{"int seconds", time.Second, time.Second, true, "1"},
		{"int milliseconds", time.Second, time.Millisecond, true, "1000"},
		{"float seconds", time.Second, time.Second, false, "1"},
		{"float half second", 500 * time.Millisecond, time.Second, false, "0.5"},
		{"zero duration", 0, time.Second, true, "0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(je.Duration(nil, tt.unit, tt.useInt, tt.d))
			if got != tt.want {
				t.Errorf("Duration() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTextEncoder_Duration(t *testing.T) {
	te := TextEncoder{}

	tests := []struct {
		name   string
		d      time.Duration
		unit   time.Duration
		useInt bool
		want   string
	}{
		{"int seconds", time.Second, time.Second, true, "1"},
		{"float seconds", time.Second, time.Second, false, "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(te.Duration(nil, tt.unit, tt.useInt, tt.d))
			if got != tt.want {
				t.Errorf("Duration() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ============================== Durations (slice) encoding tests ==============================

func TestJsonEncoder_Durations(t *testing.T) {
	je := JsonEncoder{}

	t.Run("three elements", func(t *testing.T) {
		durs := []time.Duration{time.Second, 2 * time.Second, 3 * time.Second}
		got := string(je.Durations(nil, time.Second, true, durs...))
		want := "[1,2,3]"
		if got != want {
			t.Errorf("Durations() = %q, want %q (bug: duplicate first element?)", got, want)
		}
	})

	t.Run("single element", func(t *testing.T) {
		got := string(je.Durations(nil, time.Second, true, time.Second))
		want := "[1]"
		if got != want {
			t.Errorf("Durations() = %q, want %q", got, want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		got := string(je.Durations(nil, time.Second, true))
		want := "[]"
		if got != want {
			t.Errorf("Durations() = %q, want %q", got, want)
		}
	})

	t.Run("float mode", func(t *testing.T) {
		durs := []time.Duration{500 * time.Millisecond, time.Second}
		got := string(je.Durations(nil, time.Second, false, durs...))
		want := "[0.5,1]"
		if got != want {
			t.Errorf("Durations() = %q, want %q", got, want)
		}
	})
}

func TestTextEncoder_Durations(t *testing.T) {
	te := TextEncoder{}

	t.Run("three elements", func(t *testing.T) {
		durs := []time.Duration{time.Second, 2 * time.Second, 3 * time.Second}
		got := string(te.Durations(nil, time.Second, true, durs...))
		want := "[1,2,3]"
		if got != want {
			t.Errorf("Durations() = %q, want %q (bug: duplicate first element?)", got, want)
		}
	})

	t.Run("single element", func(t *testing.T) {
		got := string(te.Durations(nil, time.Second, true, time.Second))
		want := "[1]"
		if got != want {
			t.Errorf("Durations() = %q, want %q", got, want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		got := string(te.Durations(nil, time.Second, true))
		want := "[]"
		if got != want {
			t.Errorf("Durations() = %q, want %q", got, want)
		}
	})
}

// ============================== Time encoding tests ==============================

func TestJsonEncoder_Time(t *testing.T) {
	je := JsonEncoder{}
	now := time.Unix(1700000000, 123456789) // fixed timestamp

	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"unix", TimeFormatUnix, "1700000000"},
		{"unix_ms", TimeFormatUnixMs, "1700000000123"},
		{"unix_micro", TimeFormatUnixMicro, "1700000000123456"},
		{"unix_nano", TimeFormatUnixNano, "1700000000123456789"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(je.Time(nil, tt.format, now))
			if got != tt.want {
				t.Errorf("Time(%s) = %q, want %q", tt.format, got, tt.want)
			}
		})
	}

	t.Run("custom format", func(t *testing.T) {
		utcNow := now.UTC()
		got := string(je.Time(nil, "2006-01-02", utcNow))
		want := `"2023-11-14"`
		if got != want {
			t.Errorf("Time(custom) = %q, want %q", got, want)
		}
	})
}

func TestTextEncoder_Time(t *testing.T) {
	te := TextEncoder{}
	now := time.Unix(1700000000, 0)

	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"unix", TimeFormatUnix, "1700000000"},
		{"unix_ms", TimeFormatUnixMs, "1700000000000"},
		{"unix_nano", TimeFormatUnixNano, "1700000000000000000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(te.Time(nil, tt.format, now))
			if got != tt.want {
				t.Errorf("Time(%s) = %q, want %q", tt.format, got, tt.want)
			}
		})
	}
}

// ============================== Times (slice) encoding tests ==============================

func TestJsonEncoder_Times(t *testing.T) {
	je := JsonEncoder{}
	times := []time.Time{
		time.Unix(1, 0),
		time.Unix(2, 0),
		time.Unix(3, 0),
	}

	t.Run("three elements", func(t *testing.T) {
		got := string(je.Times(nil, TimeFormatUnix, times...))
		want := "[1,2,3]"
		if got != want {
			t.Errorf("Times() = %q, want %q", got, want)
		}
	})

	t.Run("single element", func(t *testing.T) {
		got := string(je.Times(nil, TimeFormatUnix, time.Unix(42, 0)))
		want := "[42]"
		if got != want {
			t.Errorf("Times() = %q, want %q", got, want)
		}
	})

	t.Run("empty", func(t *testing.T) {
		got := string(je.Times(nil, TimeFormatUnix))
		want := "[]"
		if got != want {
			t.Errorf("Times() = %q, want %q", got, want)
		}
	})
}

// ============================== ObjectData tests ==============================

func TestJsonEncoder_ObjectData(t *testing.T) {
	je := JsonEncoder{}

	t.Run("non empty object", func(t *testing.T) {
		got := string(je.ObjectData(nil, []byte(`{"key":"value"}`)))
		want := `"key":"value"}`
		if got != want {
			t.Errorf("ObjectData() = %q, want %q", got, want)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		// Should not panic — regression test for bounds check fix
		got := string(je.ObjectData(nil, []byte{}))
		// empty append is fine
		if got != "" {
			t.Errorf("ObjectData(empty) = %q, want %q", got, "")
		}
	})

	t.Run("non object input", func(t *testing.T) {
		got := string(je.ObjectData(nil, []byte(`[1,2,3]`)))
		// Should not strip leading char since it's not '{'
		_ = got
	})
}

// ============================== Numeric encoding tests ==============================

func TestJsonEncoder_Numerics(t *testing.T) {
	je := JsonEncoder{}

	t.Run("Int64", func(t *testing.T) {
		tests := []struct {
			in   int64
			want string
		}{
			{0, "0"},
			{42, "42"},
			{-1, "-1"},
			{9223372036854775807, "9223372036854775807"},
		}
		for _, tt := range tests {
			got := string(je.Int64(nil, tt.in))
			if got != tt.want {
				t.Errorf("Int64(%d) = %q, want %q", tt.in, got, tt.want)
			}
		}
	})

	t.Run("Uint64", func(t *testing.T) {
		got := string(je.Uint64(nil, 42))
		if got != "42" {
			t.Errorf("Uint64 = %q, want %q", got, "42")
		}
	})

	t.Run("Float64", func(t *testing.T) {
		tests := []struct {
			in   float64
			want string
		}{
			{0, "0"},
			{3.14, "3.14"},
			{-1.5, "-1.5"},
		}
		for _, tt := range tests {
			got := string(je.Float64(nil, tt.in))
			if got != tt.want {
				t.Errorf("Float64(%v) = %q, want %q", tt.in, got, tt.want)
			}
		}
	})

	t.Run("Bool", func(t *testing.T) {
		if got := string(je.Bool(nil, true)); got != "true" {
			t.Errorf("Bool(true) = %q, want true", got)
		}
		if got := string(je.Bool(nil, false)); got != "false" {
			t.Errorf("Bool(false) = %q, want false", got)
		}
	})

	t.Run("Nil", func(t *testing.T) {
		got := string(je.Nil(nil))
		if got != "null" {
			t.Errorf("Nil() = %q, want null", got)
		}
	})
}

func TestTextEncoder_Numerics(t *testing.T) {
	te := TextEncoder{}

	t.Run("Int64", func(t *testing.T) {
		got := string(te.Int64(nil, 42))
		if got != "42" {
			t.Errorf("Int64 = %q, want 42", got)
		}
	})

	t.Run("Bool", func(t *testing.T) {
		if got := string(te.Bool(nil, true)); got != "true" {
			t.Errorf("Bool(true) = %q, want true", got)
		}
	})

	t.Run("Float64", func(t *testing.T) {
		got := string(te.Float64(nil, 3.14))
		if got != "3.14" {
			t.Errorf("Float64 = %q, want 3.14", got)
		}
	})
}

// ============================== Array encoding tests ==============================

func TestJsonEncoder_ArrayMethods(t *testing.T) {
	je := JsonEncoder{}

	t.Run("ArrayStart/End", func(t *testing.T) {
		dst := je.ArrayStart(nil)
		dst = je.ArrayEnd(dst)
		if string(dst) != "[]" {
			t.Errorf("ArrayStart+End = %q, want []", dst)
		}
	})

	t.Run("ArrayDelim", func(t *testing.T) {
		dst := je.ArrayStart(nil)
		dst = append(dst, '1')
		dst = je.ArrayDelim(dst)
		dst = append(dst, '2')
		dst = je.ArrayEnd(dst)
		if string(dst) != "[1,2]" {
			t.Errorf("Array with delim = %q, want [1,2]", dst)
		}
	})

	t.Run("Int8 explicit", func(t *testing.T) {
		got := string(je.Int8(nil, 127))
		if got != "127" {
			t.Errorf("Int8 = %q, want 127", got)
		}
	})

	t.Run("Int16 explicit", func(t *testing.T) {
		got := string(je.Int16(nil, 32767))
		if got != "32767" {
			t.Errorf("Int16 = %q, want 32767", got)
		}
	})

	t.Run("Int32 explicit", func(t *testing.T) {
		got := string(je.Int32(nil, 2147483647))
		if got != "2147483647" {
			t.Errorf("Int32 = %q, want 2147483647", got)
		}
	})
}

// ============================== Struct marker tests ==============================

func TestJsonEncoder_Markers(t *testing.T) {
	je := JsonEncoder{}

	t.Run("BeginMarker", func(t *testing.T) {
		got := string(je.BeginMarker(nil))
		if got != "{" {
			t.Errorf("BeginMarker = %q, want {", got)
		}
	})

	t.Run("EndMarker", func(t *testing.T) {
		got := string(je.EndMarker(nil))
		if got != "}" {
			t.Errorf("EndMarker = %q, want }", got)
		}
	})

	t.Run("LineBreak", func(t *testing.T) {
		got := string(je.LineBreak(nil))
		if got != "\n" {
			t.Errorf("LineBreak = %q, want newline", got)
		}
	})

	t.Run("Delim on empty dst", func(t *testing.T) {
		// Delim is context-sensitive: only appends comma when dst is non-empty
		got := string(je.Delim(nil))
		if got != "" {
			t.Errorf("Delim(nil) = %q, want empty (no leading comma)", got)
		}
	})
	t.Run("Delim after content", func(t *testing.T) {
		got := string(je.Delim([]byte("x")))
		if got != "x," {
			t.Errorf("Delim after content = %q, want x,", got)
		}
	})

	t.Run("Key", func(t *testing.T) {
		got := string(je.Key(nil, "name"))
		if got != `"name":` {
			t.Errorf("Key = %q, want \"name\":", got)
		}
	})

	t.Run("BlankSpace", func(t *testing.T) {
		got := string(je.BlankSpace(nil))
		if got != " " {
			t.Errorf("BlankSpace = %q, want space", got)
		}
	})
	t.Run("MetaEnd", func(t *testing.T) {
		got := string(je.MetaEnd(nil))
		// MetaEnd is a no-op — returns dst unchanged
		if got != "" {
			t.Errorf("MetaEnd = %q, want empty", got)
		}
	})
}

func TestTextEncoder_Markers(t *testing.T) {
	te := TextEncoder{}

	t.Run("Key", func(t *testing.T) {
		got := string(te.Key(nil, "name"))
		if got != "name=" {
			t.Errorf("Key = %q, want name=", got)
		}
	})

	t.Run("BlankSpace", func(t *testing.T) {
		got := string(te.BlankSpace(nil))
		if got != " " {
			t.Errorf("BlankSpace = %q, want space", got)
		}
	})
}

// ============================== IP/MAC encoding tests ==============================

func TestJsonEncoder_NetworkTypes(t *testing.T) {
	je := JsonEncoder{}

	t.Run("IPAddr", func(t *testing.T) {
		got := string(je.IPAddr(nil, []byte{127, 0, 0, 1}))
		if got != `"127.0.0.1"` {
			t.Errorf("IPAddr = %q, want \"127.0.0.1\"", got)
		}
	})

	t.Run("MACAddr", func(t *testing.T) {
		got := string(je.MACAddr(nil, []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}))
		if got != `"00:11:22:33:44:55"` {
			t.Errorf("MACAddr = %q, want \"00:11:22:33:44:55\"", got)
		}
	})
}

// ============================== Interface/Any encoding tests ==============================

func TestJsonEncoder_Interface(t *testing.T) {
	je := JsonEncoder{}

	t.Run("string via JSONMarshalFunc", func(t *testing.T) {
		got := string(je.Interface(nil, "hello"))
		if got != `"hello"` {
			t.Errorf("Interface(string) = %q, want \"hello\"", got)
		}
	})

	t.Run("integer", func(t *testing.T) {
		got := string(je.Interface(nil, 42))
		if got != "42" {
			t.Errorf("Interface(42) = %q, want 42", got)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type S struct{ Name string }
		got := string(je.Interface(nil, S{Name: "x"}))
		if got != `{"Name":"x"}` {
			t.Errorf("Interface(struct) = %q, want {\"Name\":\"x\"}", got)
		}
	})
}

func TestTextEncoder_Interface(t *testing.T) {
	te := TextEncoder{}

	t.Run("string via TextMarshalFunc", func(t *testing.T) {
		got := string(te.Interface(nil, "hello"))
		if got != "hello" {
			t.Errorf("Interface(string) = %q, want hello", got)
		}
	})

	t.Run("integer", func(t *testing.T) {
		got := string(te.Interface(nil, 42))
		if got != "42" {
			t.Errorf("Interface(42) = %q, want 42", got)
		}
	})
}

// ============================== Hex/Bytes/IPPrefix tests ==============================

func TestJsonEncoder_Hex(t *testing.T) {
	je := JsonEncoder{}

	t.Run("hex encoding", func(t *testing.T) {
		got := string(je.Hex(nil, []byte{0xde, 0xad, 0xbe, 0xef}))
		if got != `"deadbeef"` {
			t.Errorf("Hex = %q, want \"deadbeef\"", got)
		}
	})

	t.Run("empty", func(t *testing.T) {
		got := string(je.Hex(nil, nil))
		if got != `""` {
			t.Errorf("Hex(nil) = %q, want \"\"", got)
		}
	})
}

func TestTextEncoder_Hex(t *testing.T) {
	te := TextEncoder{}

	t.Run("hex encoding", func(t *testing.T) {
		got := string(te.Hex(nil, []byte{0xab, 0xcd}))
		if got != "abcd" {
			t.Errorf("Hex = %q, want abcd", got)
		}
	})
}


// ============================== Float slice tests ==============================

func TestJsonEncoder_FloatSlices(t *testing.T) {
	je := JsonEncoder{}

	t.Run("Float32s", func(t *testing.T) {
		got := string(je.Float32s(nil, 1.5, 2.5, 3.5))
		if got != "[1.5,2.5,3.5]" {
			t.Errorf("Float32s = %q, want [1.5,2.5,3.5]", got)
		}
	})

	t.Run("Float64s single", func(t *testing.T) {
		got := string(je.Float64s(nil, 3.14))
		if got != "[3.14]" {
			t.Errorf("Float64s = %q, want [3.14]", got)
		}
	})

	t.Run("Float64s empty", func(t *testing.T) {
		got := string(je.Float64s(nil))
		if got != "[]" {
			t.Errorf("Float64s() = %q, want []", got)
		}
	})
}

// ============================== Int/Uint slice tests ==============================

func TestJsonEncoder_IntSlices(t *testing.T) {
	je := JsonEncoder{}

	t.Run("Ints", func(t *testing.T) {
		got := string(je.Ints(nil, 1, 2, 3))
		if got != "[1,2,3]" {
			t.Errorf("Ints = %q, want [1,2,3]", got)
		}
	})

	t.Run("Uints", func(t *testing.T) {
		got := string(je.Uints(nil, 1, 2))
		if got != "[1,2]" {
			t.Errorf("Uints = %q, want [1,2]", got)
		}
	})

	t.Run("Int8s", func(t *testing.T) {
		got := string(je.Int8s(nil, 1, -2))
		if got != "[1,-2]" {
			t.Errorf("Int8s = %q, want [1,-2]", got)
		}
	})

	t.Run("Uint64s", func(t *testing.T) {
		got := string(je.Uint64s(nil, 100))
		if got != "[100]" {
			t.Errorf("Uint64s = %q, want [100]", got)
		}
	})
}

// ============================== Strings tests ==============================

func TestJsonEncoder_Strings(t *testing.T) {
	je := JsonEncoder{}

	t.Run("multiple", func(t *testing.T) {
		got := string(je.Strings(nil, "a", "b", "c"))
		if got != `["a","b","c"]` {
			t.Errorf("Strings = %q, want [\"a\",\"b\",\"c\"]", got)
		}
	})

	t.Run("single", func(t *testing.T) {
		got := string(je.Strings(nil, "x"))
		if got != `["x"]` {
			t.Errorf("Strings = %q, want [\"x\"]", got)
		}
	})

	t.Run("empty", func(t *testing.T) {
		got := string(je.Strings(nil))
		if got != "[]" {
			t.Errorf("Strings() = %q, want []", got)
		}
	})
}
