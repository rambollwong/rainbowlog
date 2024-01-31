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

// Record represents a log record.
// It is finalized by the Done method.
// Done method also writes the encoded bytes to the writer.
type Record struct {
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
func (r *Record) WithLabels(label ...string) *Record {
	r.label = strings.Join(label, "|")
	return r
}

// WithLevel sets the level as the META_LEVEL field.
func (r *Record) WithLevel(lv level.Level) *Record {
	r.level = lv
	return r
}

func (r *Record) WithDoneFunc(f func(msg string)) *Record {
	r.doneFunc = f
	return r
}

// WithCallerSkip adds skip frames when calling caller.
func (r *Record) WithCallerSkip(skip int) *Record {
	for _, rp := range r.recordPackers {
		rp.CallerSkip(skip)
	}
	return r
}

// UseIntDur enables the calculation of Duration to retain only integer digits.
func (r *Record) UseIntDur() *Record {
	r.useIntDur = true
	return r
}

func (r *Record) strikeOrNot() bool {
	if r.level >= r.logger.level {
		return false
	}
	return true
}

// Reset the record instance.
// This should be called before the record is used again.
func (r *Record) Reset() {
	for _, rp := range r.recordPackers {
		rp.Reset()
	}
	r.level = level.Disabled
	r.label = ""
	r.msg = ""
	r.doneFunc = nil
}

// Discard disables the record that it won't be printed.
func (r *Record) Discard() *Record {
	r.level = level.Disabled
	return r
}

func (r *Record) fatalOrPanic() {
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

// Done finish appending data and writing record.
//
// NOTICE: once this method is called, the *Record should be disposed.
// Calling Done twice can have unexpected result.
func (r *Record) Done() {
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
func (r *Record) Msg(msg string) *Record {
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
func (r *Record) Msgf(format string, v ...interface{}) *Record {
	return r.Msg(fmt.Sprintf(format, v...))
}

// Err sets the given err to the error field when err is not nil。
// NOTICE: This method should only be called once。
// Calling multiple times may cause unpredictable results。
func (r *Record) Err(err error) *Record {
	if r.strikeOrNot() || err == nil {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Err(err)
	}
	return r
}

func (r *Record) Str(key, val string) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Str(key, val)
	}
	return r
}

func (r *Record) Strs(key string, vals ...string) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Strs(key, vals...)
	}
	return r
}

func (r *Record) Stringer(key string, val fmt.Stringer) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Stringer(key, val)
	}
	return r
}

func (r *Record) Stringers(key string, vals ...fmt.Stringer) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Stringers(key, vals...)
	}
	return r
}

func (r *Record) Bytes(key string, val []byte) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Bytes(key, val)
	}
	return r
}

func (r *Record) Hex(key string, val []byte) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Hex(key, val)
	}
	return r
}

func (r *Record) Int(key string, val int) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int(key, val)
	}
	return r
}

func (r *Record) Ints(key string, vals ...int) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Ints(key, vals...)
	}
	return r
}

func (r *Record) Int8(key string, val int8) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int8(key, val)
	}
	return r
}

func (r *Record) Int8s(key string, vals ...int8) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int8s(key, vals...)
	}
	return r
}

func (r *Record) Int16(key string, val int16) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int16(key, val)
	}
	return r
}

func (r *Record) Int16s(key string, vals ...int16) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int16s(key, vals...)
	}
	return r
}

func (r *Record) Int32(key string, val int32) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int32(key, val)
	}
	return r
}

func (r *Record) Int32s(key string, vals ...int32) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int32s(key, vals...)
	}
	return r
}

func (r *Record) Int64(key string, val int64) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int64(key, val)
	}
	return r
}

func (r *Record) Int64s(key string, vals ...int64) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Int64s(key, vals...)
	}
	return r
}

func (r *Record) Uint(key string, val uint) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint(key, val)
	}
	return r
}

func (r *Record) Uints(key string, vals ...uint) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uints(key, vals...)
	}
	return r
}

func (r *Record) Uint8(key string, val uint8) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint8(key, val)
	}
	return r
}

func (r *Record) Uint8s(key string, vals ...uint8) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint8s(key, vals...)
	}
	return r
}

func (r *Record) Uint16(key string, val uint16) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint16(key, val)
	}
	return r
}

func (r *Record) Uint16s(key string, vals ...uint16) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint16s(key, vals...)
	}
	return r
}

func (r *Record) Uint32(key string, val uint32) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint32(key, val)
	}
	return r
}

func (r *Record) Uint32s(key string, vals ...uint32) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint32s(key, vals...)
	}
	return r
}

func (r *Record) Uint64(key string, val uint64) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint64(key, val)
	}
	return r
}

func (r *Record) Uint64s(key string, vals ...uint64) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Uint64s(key, vals...)
	}
	return r
}

func (r *Record) Float32(key string, val float32) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Float32(key, val)
	}
	return r
}

func (r *Record) Float32s(key string, vals ...float32) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Float32s(key, vals...)
	}
	return r
}

func (r *Record) Float64(key string, val float64) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Float64(key, val)
	}
	return r
}

func (r *Record) Float64s(key string, vals ...float64) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Float64s(key, vals...)
	}
	return r
}

func (r *Record) Time(key, format string, val time.Time) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Time(key, format, val)
	}
	return r
}

func (r *Record) Times(key, format string, vals ...time.Time) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Times(key, format, vals...)
	}
	return r
}

func (r *Record) Dur(key string, unit, val time.Duration) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Dur(key, unit, val)
	}
	return r
}

func (r *Record) Durs(key string, unit time.Duration, vals ...time.Duration) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Durs(key, unit, vals...)
	}
	return r
}

func (r *Record) Any(key string, i any) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.Any(key, i)
	}
	return r
}

func (r *Record) IPAddr(key string, ip net.IP) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.IPAddr(key, ip)
	}
	return r
}

func (r *Record) IPPrefix(key string, pfx net.IPNet) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.IPPrefix(key, pfx)
	}
	return r
}

func (r *Record) MACAddr(key string, ha net.HardwareAddr) *Record {
	if r.strikeOrNot() {
		return r
	}
	for _, rp := range r.recordPackers {
		rp.MACAddr(key, ha)
	}
	return r
}
