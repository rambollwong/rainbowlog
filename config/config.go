package config

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/rambollwong/rainbowcat/writer/filewriter"
	"github.com/spf13/viper"
)

type LoggerConfig struct {
	Enable                bool                        `mapstructure:"enable" json:"enable" yaml:"enable"`
	Level                 string                      `mapstructure:"level" json:"level" yaml:"level"`
	Label                 string                      `mapstructure:"label" json:"label" yaml:"label"`
	Stack                 bool                        `mapstructure:"stack" json:"stack" yaml:"stack"`
	EnableConsolePrinting bool                        `mapstructure:"enableConsolePrinting" json:"enableConsolePrinting" yaml:"enableConsolePrinting"`
	EnableRainbowConsole  bool                        `mapstructure:"enableRainbowConsole" json:"enableRainbowConsole" yaml:"enableRainbowConsole"`
	TimeFormat            string                      `mapstructure:"timeFormat" json:"timeFormat" yaml:"timeFormat"`
	SizeRollingFileConfig LoggerSizeRollingFileConfig `mapstructure:"sizeRollingFileConfig" json:"sizeRollingFileConfig" yaml:"sizeRollingFileConfig"`
	TimeRollingFileConfig LoggerTimeRollingFileConfig `mapstructure:"timeRollingFileConfig" json:"timeRollingFileConfig" yaml:"timeRollingFileConfig"`
}

type LoggerTimeRollingFileConfig struct {
	Enable            bool                     `mapstructure:"enable" json:"enable" yaml:"enable"`
	LogFilePath       string                   `mapstructure:"logFilePath" json:"logFilePath" yaml:"logFilePath"`
	LogFileBaseName   string                   `mapstructure:"logFileBaseName" json:"logFileBaseName" yaml:"logFileBaseName"`
	MaxBackups        int                      `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"`
	RollingPeriod     filewriter.RollingPeriod `mapstructure:"rollingPeriod" json:"rollingPeriod" yaml:"rollingPeriod"`
	Encoder           string                   `mapstructure:"encoder" json:"encoder" yaml:"encoder"`
	UseBufferedWriter bool                     `mapstructure:"useBufferedWriter" json:"useBufferedWriter" yaml:"useBufferedWriter"`
	WriterBufferSize  string                   `mapstructure:"writerBufferSize" json:"writerBufferSize" yaml:"writerBufferSize"`
}

type LoggerSizeRollingFileConfig struct {
	Enable            bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	LogFilePath       string `mapstructure:"logFilePath" json:"logFilePath" yaml:"logFilePath"`
	LogFileBaseName   string `mapstructure:"logFileBaseName" json:"logFileBaseName" yaml:"logFileBaseName"`
	MaxBackups        int    `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"`
	FileSizeLimit     string `mapstructure:"fileSizeLimit" json:"fileSizeLimit" yaml:"fileSizeLimit"`
	Encoder           string `mapstructure:"encoder" json:"encoder" yaml:"encoder"`
	UseBufferedWriter bool   `mapstructure:"useBufferedWriter" json:"useBufferedWriter" yaml:"useBufferedWriter"`
	WriterBufferSize  string `mapstructure:"writerBufferSize" json:"writerBufferSize" yaml:"writerBufferSize"`
}

func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Enable:                true,
		Level:                 "DEBUG",
		Label:                 "",
		Stack:                 false,
		EnableConsolePrinting: true,
		EnableRainbowConsole:  true,
		TimeFormat:            "2006-01-02 15:04:05.000",
		SizeRollingFileConfig: LoggerSizeRollingFileConfig{
			Enable:            false,
			LogFilePath:       "./log",
			LogFileBaseName:   "rainbow.log",
			MaxBackups:        10,
			FileSizeLimit:     "100M",
			Encoder:           "json",
			UseBufferedWriter: true,
			WriterBufferSize:  "4K",
		},
		TimeRollingFileConfig: LoggerTimeRollingFileConfig{
			Enable:            false,
			LogFilePath:       "./log",
			LogFileBaseName:   "rainbow.log",
			MaxBackups:        7,
			RollingPeriod:     "DAY",
			Encoder:           "json",
			UseBufferedWriter: true,
			WriterBufferSize:  "4K",
		},
	}
}

// LoadLoggerConfigFromFile loads the logger configuration from a file.
func LoadLoggerConfigFromFile(configFile string) (*LoggerConfig, error) {
	fileName := filepath.Base(configFile)
	fileType := strings.TrimPrefix(filepath.Ext(fileName), ".")
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return LoadLoggerConfigFromViper(v)
}

// LoadLoggerConfigFromViper loads the logger configuration from a *viper.Viper instance.
// It unmarshals the value under the "rainbowlog" key into a *LoggerConfig.
// Environment variables prefixed with "RAINBOWLOG_" override the values loaded from the viper instance.
// Returns the loaded *LoggerConfig, or an error if unmarshaling fails.
func LoadLoggerConfigFromViper(v *viper.Viper) (*LoggerConfig, error) {
	sub := v.Sub("rainbowlog")
	if sub == nil {
		sub = viper.New()
		sub.SetEnvPrefix("RAINBOWLOG")
	}
	// Sub's parents already carry the "rainbowlog" prefix; combined with the
	// "." → "_" replacer it maps to environment variables like "RAINBOWLOG_LEVEL".
	// Setting envPrefix again would duplicate the prefix.
	sub.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	sub.AutomaticEnv()

	cfg := new(LoggerConfig)
	if err := sub.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// LoadLoggerConfigFromEnv loads the logger configuration from environment variables.
// Variables use the "<envPrefix>_" prefix with underscores as separators.
// If envPrefix is empty, "RAINBOWLOG" is used as the default prefix.
// e.g. RAINBOWLOG_LEVEL, RAINBOWLOG_ENABLECONSOLEPRINTING, RAINBOWLOG_SIZEROLLINGFILECONFIG_ENABLE.
// Returns the loaded *LoggerConfig, or an error if unmarshaling fails.
func LoadLoggerConfigFromEnv(envPrefix string) (*LoggerConfig, error) {
	if envPrefix == "" {
		envPrefix = "RAINBOWLOG"
	}
	v := viper.New()
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	// AutomaticEnv is lazy: Unmarshal relies on AllKeys to enumerate keys.
	// Register zero-value defaults for every field so environment variables
	// can be read and override them.
	for _, key := range loggerConfigEnvKeys {
		v.SetDefault(key, "")
	}

	cfg := new(LoggerConfig)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// loggerConfigEnvKeys lists every LoggerConfig field as a flattened key,
// used by LoadLoggerConfigFromEnv to register zero-value defaults so that
// AutomaticEnv can enumerate and override them.
var loggerConfigEnvKeys = []string{
	"enable", "level", "label", "stack",
	"enableconsoleprinting", "enablerainbowconsole", "timeformat",
	"sizerollingfileconfig.enable", "sizerollingfileconfig.logfilepath",
	"sizerollingfileconfig.logfilebasename", "sizerollingfileconfig.maxbackups",
	"sizerollingfileconfig.filesizelimit", "sizerollingfileconfig.encoder",
	"sizerollingfileconfig.usebufferedwriter", "sizerollingfileconfig.writerbuffersize",
	"timerollingfileconfig.enable", "timerollingfileconfig.logfilepath",
	"timerollingfileconfig.logfilebasename", "timerollingfileconfig.maxbackups",
	"timerollingfileconfig.rollingperiod", "timerollingfileconfig.encoder",
	"timerollingfileconfig.usebufferedwriter", "timerollingfileconfig.writerbuffersize",
}

// WriteConfigToFile writes the logger configuration to a file.
// configFilePath: Path to the configuration file.
// config: LoggerConfig object containing the configuration.
// Returns an error if the write operation fails.
func WriteConfigToFile(configFilePath string, config LoggerConfig) error {
	v := viper.New()

	configMap := make(map[string]any)
	configMap["rainbowlog"] = config

	v.SetConfigFile(configFilePath)

	for key, value := range configMap {
		v.Set(key, value)
	}

	err := v.WriteConfig()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err = v.SafeWriteConfig()
		}
	}

	return err
}
