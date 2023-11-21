//go:build !binary_log

package rainbowlog_test

import (
	"os"

	"github.com/rambollwong/rainbowlog"
	"github.com/rambollwong/rainbowlog/level"
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

	// Output: {"META_LEVEL":"INFO","message":"Hello world!"}
	//DEBUG > Hello world!
}

func ExampleLogger_SubLogger() {
	log := rainbowlog.New(rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout))

	subLog := log.SubLogger(
		rainbowlog.WithMetaKeys(rainbowlog.MetaLabelFieldName),
		rainbowlog.WithLabel("NewLabel"),
	)
	subLog.Info().Msg("Hello world!").Done()
	// Output: {"META_LABEL":"NewLabel","message":"Hello world!"}
}

func ExampleLogger_Record() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MetaLevelFieldName),
	)
	log.Record().WithLevel(level.Info).Msg("Hello world!").Done()
	log.Record().WithLevel(level.Debug).Msg("Hello world!").Done()

	// Output: {"META_LEVEL":"INFO","message":"Hello world!"}
	//{"META_LEVEL":"DEBUG","message":"Hello world!"}
}

func ExampleLogger_Info() {
	log := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MetaLevelFieldName),
	)
	log.Info().Msg("Hello world!").Done()

	// Output: {"META_LEVEL":"INFO","message":"Hello world!"}
}

// TODO: Debug

// TODO: Error

// TODO: CustomMetaKeyFieldName

// TODO: CustomTimestamp
