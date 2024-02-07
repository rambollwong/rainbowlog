package rainbowlog

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rambollwong/rainbowlog/level"
)

var (
	_         Record = (*LogRecord)(nil)
	nilRecord Record = &NilRecord{}
)

type Record interface {
	WithLabels(label ...string) Record
	WithLevel(lv level.Level) Record
	WithDoneFunc(f func(msg string)) Record
	WithCallerSkip(skip int) Record
	UseIntDur() Record
	Reset()
	Discard() Record
	Done()
	Msg(msg string) Record
	Msgf(format string, v ...interface{}) Record
	Err(err error) Record
	Str(key string, val string) Record
	Strs(key string, vals ...string) Record
	Stringer(key string, val fmt.Stringer) Record
	Stringers(key string, vals ...fmt.Stringer) Record
	Bytes(key string, val []byte) Record
	Hex(key string, val []byte) Record
	Int(key string, val int) Record
	Ints(key string, vals ...int) Record
	Int8(key string, val int8) Record
	Int8s(key string, vals ...int8) Record
	Int16(key string, val int16) Record
	Int16s(key string, vals ...int16) Record
	Int32(key string, val int32) Record
	Int32s(key string, vals ...int32) Record
	Int64(key string, val int64) Record
	Int64s(key string, vals ...int64) Record
	Uint(key string, val uint) Record
	Uints(key string, vals ...uint) Record
	Uint8(key string, val uint8) Record
	Uint8s(key string, vals ...uint8) Record
	Uint16(key string, val uint16) Record
	Uint16s(key string, vals ...uint16) Record
	Uint32(key string, val uint32) Record
	Uint32s(key string, vals ...uint32) Record
	Uint64(key string, val uint64) Record
	Uint64s(key string, vals ...uint64) Record
	Float32(key string, val float32) Record
	Float32s(key string, vals ...float32) Record
	Float64(key string, val float64) Record
	Float64s(key string, vals ...float64) Record
	Time(key string, format string, val time.Time) Record
	Times(key string, format string, vals ...time.Time) Record
	Dur(key string, unit time.Duration, val time.Duration) Record
	Durs(key string, unit time.Duration, vals ...time.Duration) Record
	Any(key string, i any) Record
	IPAddr(key string, ip net.IP) Record
	IPPrefix(key string, pfx net.IPNet) Record
	MACAddr(key string, ha net.HardwareAddr) Record
}

// LogRecord represents a log Record.
// It is finalized by the Done method.
// Done method also writes the encoded bytes to the writer.
type LogRecord struct {
	mu            sync.Mutex
	recordPackers []recordPacker
	level         level.Level
	label         string
	msg           string
	useIntDur     bool
	stack         bool
	doneFunc      func(msg string)

	logger *Logger
}

// WithLabels sets the labels for the Logger using one or more strings.
// The labels are joined using the "|" separator.
func (r *LogRecord) WithLabels(label ...string) Record {
	r.label = strings.Join(label, "|")
	return r
}

// WithLevel sets the level as the META_LEVEL field.
func (r *LogRecord) WithLevel(lv level.Level) Record {
	r.level = lv
	return r
}

func (r *LogRecord) WithDoneFunc(f func(msg string)) Record {
	r.doneFunc = f
	return r
}

// WithCallerSkip adds skip frames when calling caller.
func (r *LogRecord) WithCallerSkip(skip int) Record {
	for _, rp := range r.recordPackers {
		rp.CallerSkip(skip)
	}
	return r
}

// UseIntDur enables the calculation of Duration to retain only integer digits.
func (r *LogRecord) UseIntDur() Record {
	r.useIntDur = true
	return r
}

func (r *LogRecord) strikeOrNot() bool {
	if r.level >= r.logger.level {
		return false
	}
	return true
}

// Reset the Record instance.
// This should be called before the Record is used again.
func (r *LogRecord) Reset() {
	for _, rp := range r.recordPackers {
		rp.Reset()
	}
	r.level = level.Disabled
	r.label = ""
	r.msg = ""
	r.doneFunc = nil
}

// Discard disables the Record that it won't be printed.
func (r *LogRecord) Discard() Record {
	r.level = level.Disabled
	return r
}

func (r *LogRecord) fatalOrPanic() {
	switch r.level {
	case level.Fatal:
		os.Exit(1)
	case level.Panic:
		if r.msg == "" {
			panic("panic")
		}
		panic(r.msg)
	default:
	}
}

// Done finish appending data and writing Record.
//
// NOTICE: once this method is called, the *Record should be disposed.
// Calling Done twice can have unexpected result.
func (r *LogRecord) Done() {
	for _, hook := range r.logger.hooks {
		hook.RunHook(r, r.level, r.msg)
	}
	if r.doneFunc != nil {
		defer r.doneFunc(r.msg)
	}

	// recycling
	defer r.logger.recordPool.Put(r)
	defer r.fatalOrPanic()

	if r.strikeOrNot() {
		return
	}

	for _, rp := range r.recordPackers {
		rp.Done()
	}
}

// Msg adds the msg as the message field if not empty.
// NOTICE: This method should only be called once。
// Calling multiple times may cause unpredictable results。
func (r *LogRecord) Msg(msg string) Record {
	if r.strikeOrNot() || msg == "" {
		return r
	}
	r.msg = msg
	for _, rp := range r.recordPackers {
		rp.Msg(msg)
	}
	return r
}

// Msgf adds the formatted msg as the message field if not empty.
func (r *LogRecord) Msgf(format string, v ...interface{}) Record {
	return r.Msg(fmt.Sprintf(format, v...))
}

// Err sets the given err to the error field when err is not nil。
// NOTICE: This method should only be called once。
// Calling multiple times may cause unpredictable results。
func (r *LogRecord) Err(err error) Record {
	if r.strikeOrNot() || err == nil {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Err(err)
	}
	return r
}

func (r *LogRecord) Str(key, val string) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Str(key, val)
	}
	return r
}

func (r *LogRecord) Strs(key string, vals ...string) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Strs(key, vals...)
	}
	return r
}

func (r *LogRecord) Stringer(key string, val fmt.Stringer) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Stringer(key, val)
	}
	return r
}

func (r *LogRecord) Stringers(key string, vals ...fmt.Stringer) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Stringers(key, vals...)
	}
	return r
}

func (r *LogRecord) Bytes(key string, val []byte) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Bytes(key, val)
	}
	return r
}

func (r *LogRecord) Hex(key string, val []byte) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Hex(key, val)
	}
	return r
}

func (r *LogRecord) Int(key string, val int) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int(key, val)
	}
	return r
}

func (r *LogRecord) Ints(key string, vals ...int) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Ints(key, vals...)
	}
	return r
}

func (r *LogRecord) Int8(key string, val int8) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int8(key, val)
	}
	return r
}

func (r *LogRecord) Int8s(key string, vals ...int8) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int8s(key, vals...)
	}
	return r
}

func (r *LogRecord) Int16(key string, val int16) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int16(key, val)
	}
	return r
}

func (r *LogRecord) Int16s(key string, vals ...int16) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int16s(key, vals...)
	}
	return r
}

func (r *LogRecord) Int32(key string, val int32) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int32(key, val)
	}
	return r
}

func (r *LogRecord) Int32s(key string, vals ...int32) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int32s(key, vals...)
	}
	return r
}

func (r *LogRecord) Int64(key string, val int64) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int64(key, val)
	}
	return r
}

func (r *LogRecord) Int64s(key string, vals ...int64) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int64s(key, vals...)
	}
	return r
}

func (r *LogRecord) Uint(key string, val uint) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint(key, val)
	}
	return r
}

func (r *LogRecord) Uints(key string, vals ...uint) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uints(key, vals...)
	}
	return r
}

func (r *LogRecord) Uint8(key string, val uint8) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint8(key, val)
	}
	return r
}

func (r *LogRecord) Uint8s(key string, vals ...uint8) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint8s(key, vals...)
	}
	return r
}

func (r *LogRecord) Uint16(key string, val uint16) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint16(key, val)
	}
	return r
}

func (r *LogRecord) Uint16s(key string, vals ...uint16) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint16s(key, vals...)
	}
	return r
}

func (r *LogRecord) Uint32(key string, val uint32) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint32(key, val)
	}
	return r
}

func (r *LogRecord) Uint32s(key string, vals ...uint32) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint32s(key, vals...)
	}
	return r
}

func (r *LogRecord) Uint64(key string, val uint64) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint64(key, val)
	}
	return r
}

func (r *LogRecord) Uint64s(key string, vals ...uint64) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint64s(key, vals...)
	}
	return r
}

func (r *LogRecord) Float32(key string, val float32) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Float32(key, val)
	}
	return r
}

func (r *LogRecord) Float32s(key string, vals ...float32) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Float32s(key, vals...)
	}
	return r
}

func (r *LogRecord) Float64(key string, val float64) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Float64(key, val)
	}
	return r
}

func (r *LogRecord) Float64s(key string, vals ...float64) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Float64s(key, vals...)
	}
	return r
}

func (r *LogRecord) Time(key, format string, val time.Time) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Time(key, format, val)
	}
	return r
}

func (r *LogRecord) Times(key, format string, vals ...time.Time) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Times(key, format, vals...)
	}
	return r
}

func (r *LogRecord) Dur(key string, unit, val time.Duration) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Dur(key, unit, val)
	}
	return r
}

func (r *LogRecord) Durs(key string, unit time.Duration, vals ...time.Duration) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Durs(key, unit, vals...)
	}
	return r
}

func (r *LogRecord) Any(key string, i any) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Any(key, i)
	}
	return r
}

func (r *LogRecord) IPAddr(key string, ip net.IP) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.IPAddr(key, ip)
	}
	return r
}

func (r *LogRecord) IPPrefix(key string, pfx net.IPNet) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.IPPrefix(key, pfx)
	}
	return r
}

func (r *LogRecord) MACAddr(key string, ha net.HardwareAddr) Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.MACAddr(key, ha)
	}
	return r
}

type NilRecord struct {
}

func (n *NilRecord) WithLabels(label ...string) Record {
	return n
}

func (n *NilRecord) WithLevel(lv level.Level) Record {
	return n
}

func (n *NilRecord) WithDoneFunc(f func(msg string)) Record {
	return n
}

func (n *NilRecord) WithCallerSkip(skip int) Record {
	return n
}

func (n *NilRecord) UseIntDur() Record {
	return n
}

func (n *NilRecord) Reset() {
	return
}

func (n *NilRecord) Discard() Record {
	return n
}

func (n *NilRecord) Done() {
	return
}

func (n *NilRecord) Msg(msg string) Record {
	return n
}

func (n *NilRecord) Msgf(format string, v ...interface{}) Record {
	return n
}

func (n *NilRecord) Err(err error) Record {
	return n
}

func (n *NilRecord) Str(key string, val string) Record {
	return n
}

func (n *NilRecord) Strs(key string, vals ...string) Record {
	return n
}

func (n *NilRecord) Stringer(key string, val fmt.Stringer) Record {
	return n
}

func (n *NilRecord) Stringers(key string, vals ...fmt.Stringer) Record {
	return n
}

func (n *NilRecord) Bytes(key string, val []byte) Record {
	return n
}

func (n *NilRecord) Hex(key string, val []byte) Record {
	return n
}

func (n *NilRecord) Int(key string, val int) Record {
	return n
}

func (n *NilRecord) Ints(key string, vals ...int) Record {
	return n
}

func (n *NilRecord) Int8(key string, val int8) Record {
	return n
}

func (n *NilRecord) Int8s(key string, vals ...int8) Record {
	return n
}

func (n *NilRecord) Int16(key string, val int16) Record {
	return n
}

func (n *NilRecord) Int16s(key string, vals ...int16) Record {
	return n
}

func (n *NilRecord) Int32(key string, val int32) Record {
	return n
}

func (n *NilRecord) Int32s(key string, vals ...int32) Record {
	return n
}

func (n *NilRecord) Int64(key string, val int64) Record {
	return n
}

func (n *NilRecord) Int64s(key string, vals ...int64) Record {
	return n
}

func (n *NilRecord) Uint(key string, val uint) Record {
	return n
}

func (n *NilRecord) Uints(key string, vals ...uint) Record {
	return n
}

func (n *NilRecord) Uint8(key string, val uint8) Record {
	return n
}

func (n *NilRecord) Uint8s(key string, vals ...uint8) Record {
	return n
}

func (n *NilRecord) Uint16(key string, val uint16) Record {
	return n
}

func (n *NilRecord) Uint16s(key string, vals ...uint16) Record {
	return n
}

func (n *NilRecord) Uint32(key string, val uint32) Record {
	return n
}

func (n *NilRecord) Uint32s(key string, vals ...uint32) Record {
	return n
}

func (n *NilRecord) Uint64(key string, val uint64) Record {
	return n
}

func (n *NilRecord) Uint64s(key string, vals ...uint64) Record {
	return n
}

func (n *NilRecord) Float32(key string, val float32) Record {
	return n
}

func (n *NilRecord) Float32s(key string, vals ...float32) Record {
	return n
}

func (n *NilRecord) Float64(key string, val float64) Record {
	return n
}

func (n *NilRecord) Float64s(key string, vals ...float64) Record {
	return n
}

func (n *NilRecord) Time(key string, format string, val time.Time) Record {
	return n
}

func (n *NilRecord) Times(key string, format string, vals ...time.Time) Record {
	return n
}

func (n *NilRecord) Dur(key string, unit time.Duration, val time.Duration) Record {
	return n
}

func (n *NilRecord) Durs(key string, unit time.Duration, vals ...time.Duration) Record {
	return n
}

func (n *NilRecord) Any(key string, i any) Record {
	return n
}

func (n *NilRecord) IPAddr(key string, ip net.IP) Record {
	return n
}

func (n *NilRecord) IPPrefix(key string, pfx net.IPNet) Record {
	return n
}

func (n *NilRecord) MACAddr(key string, ha net.HardwareAddr) Record {
	return n
}
