//go:build !binary_log

package log_test

import (
	"errors"
	"net"
	"os"
	"time"

	"github.com/rambollwong/rainbowlog/level"
	rainbowlog "github.com/rambollwong/rainbowlog/log"
)

func ExampleNew() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MetaLevelFieldName, rainbowlog.MsgFieldName),
	)
	log.Info().Msg("Hello world!").Done()

	log2 := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.TextEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MetaLevelFieldName, rainbowlog.MsgFieldName))
	log2.Debug().Msg("Hello world!").Done()

	// Output: {"_LEVEL_":"INFO","message":"Hello world!"}
	//DEBUG > Hello world!
}

func ExampleLogger_SubLogger() {
	log := rainbowlog.New(rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout))

	subLog := log.SubLogger(
		rainbowlog.WithMetaKeys(rainbowlog.MetaLabelFieldName),
		rainbowlog.WithLabels("NewLabel"),
	)
	subLog.Info().Msg("Hello world!").Done()
	// Output: {"_LABEL_":"NewLabel","message":"Hello world!"}
}

func ExampleLogger_Record() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MetaLevelFieldName),
	)
	log.Record().WithLevel(level.Info).Msg("Hello world!").Done()
	log.Record().WithLevel(level.Debug).Msg("Hello world!").Done()

	// Output: {"_LEVEL_":"INFO","message":"Hello world!"}
	//{"_LEVEL_":"DEBUG","message":"Hello world!"}
}

func ExampleLogger_Info() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MetaLevelFieldName),
	)
	log.Info().Msg("Hello world!").Done()

	// Output: {"_LEVEL_":"INFO","message":"Hello world!"}
}

func ExampleLogger_Debug() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MetaLevelFieldName),
		rainbowlog.WithLevel(level.Debug),
	)
	log.Debug().Msg("debug message").Done()
	// Output: {"_LEVEL_":"DEBUG","message":"debug message"}
}

func ExampleLogger_Error() {
	err := errors.New("something went wrong")

	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MetaLevelFieldName, rainbowlog.ErrFieldName),
	)
	log.Error().Err(err).Msg("request failed").Done()
	// Output: {"_LEVEL_":"ERROR","error":"something went wrong","message":"request failed"}
}

func Example_customTimestamp() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MetaTimeFieldName, rainbowlog.MsgFieldName),
		rainbowlog.WithTimeFormat(rainbowlog.TimeFormatUnix),
	)
	log.Info().Msg("unix timestamp").Done()
}

func Example_disabledLogger() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithLevel(level.Disabled),
	)
	// All log calls are silently discarded when the logger level is Disabled.
	log.Info().Msg("this will not be logged").Done()
	// Output:
}

func Example_fieldTypes() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MsgFieldName),
	)

	log.Info().
		Str("str", "value").
		Int64("count", 42).
		Float64("pi", 3.14).
		Msg("field types").Done()
}

func Example_timeAndDuration() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MsgFieldName),
	)

	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	log.Info().
		Time("timestamp", rainbowlog.TimeFormatUnix, now).
		Dur("latency", time.Millisecond, 150*time.Millisecond).
		Msg("time and duration").Done()
}

func Example_bytesAndHex() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MsgFieldName),
	)

	log.Info().
		Bytes("data", []byte("hello")).
		Hex("hash", []byte{0xde, 0xad, 0xbe, 0xef}).
		Msg("bytes and hex").Done()
}

func Example_networkTypes() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MsgFieldName),
	)

	log.Info().
		IPAddr("ip", net.ParseIP("192.168.1.1")).
		MACAddr("mac", net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}).
		Msg("network types").Done()
}

func Example_stringerSlice() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MsgFieldName),
	)

	log.Info().
		Strs("tags", "api", "v1", "production").
		Int64s("ids", 1, 2, 3).
		Msg("slice fields").Done()
}
