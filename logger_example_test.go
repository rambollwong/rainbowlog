package rainbowlog_test

import (
	"github.com/rambollwong/rainbowlog"
	"os"
)

func ExampleNew() {
	log := rainbowlog.New(rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout))
	log.Info().Msg("Hello world!").Done()

	log2 := rainbowlog.New(
		rainbowlog.AppendsEncoderWriters(rainbowlog.TextEnc, os.Stdout),
		rainbowlog.WithMetaKeys(rainbowlog.MsgFieldName))
	log2.Info().Msg("Hello world!").Done()

	// Output: {"message":"Hello world!"}
	//Hello world!
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

//func ExampleLogger_Record() {
//	log := rainbowlog.New(rainbowlog.AppendsEncoderWriters(rainbowlog.JsonEnc, os.Stdout))
//}
