package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rambollwong/rainbowcat/filewriter"
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
	Enable          bool                     `mapstructure:"enable" json:"enable" yaml:"enable"`
	LogFilePath     string                   `mapstructure:"logFilePath" json:"logFilePath" yaml:"logFilePath"`
	LogFileBaseName string                   `mapstructure:"logFileBaseName" json:"logFileBaseName" yaml:"logFileBaseName"`
	MaxBackups      int                      `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"`
	RollingPeriod   filewriter.RollingPeriod `mapstructure:"rollingPeriod" json:"rollingPeriod" yaml:"rollingPeriod"`
	Encoder         string                   `mapstructure:"encoder" json:"encoder" yaml:"encoder"`
}

type LoggerSizeRollingFileConfig struct {
	Enable          bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	LogFilePath     string `mapstructure:"logFilePath" json:"logFilePath" yaml:"logFilePath"`
	LogFileBaseName string `mapstructure:"logFileBaseName" json:"logFileBaseName" yaml:"logFileBaseName"`
	MaxBackups      int    `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"`
	FileSizeLimit   string `mapstructure:"fileSizeLimit" json:"fileSizeLimit" yaml:"fileSizeLimit"`
	Encoder         string `mapstructure:"encoder" json:"encoder" yaml:"encoder"`
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
			Enable:          false,
			LogFilePath:     "./log",
			LogFileBaseName: "rainbow.log",
			MaxBackups:      10,
			FileSizeLimit:   "100M",
			Encoder:         "json",
		},
		TimeRollingFileConfig: LoggerTimeRollingFileConfig{
			Enable:          false,
			LogFilePath:     "./log",
			LogFileBaseName: "rainbow.log",
			MaxBackups:      7,
			RollingPeriod:   "DAY",
			Encoder:         "json",
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
	cfg := new(LoggerConfig)
	if err := v.UnmarshalKey("rainbowlog", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// WriteConfigToFile writes the logger configuration to a file.
// configFilePath: Path to the configuration file.
// config: LoggerConfig object containing the configuration.
// Returns an error if the write operation fails.
func WriteConfigToFile(configFilePath string, config LoggerConfig) error {
	v := viper.New()

	configMap := make(map[string]interface{})
	configMap["rainbowlog"] = config

	v.SetConfigFile(configFilePath)

	for key, value := range configMap {
		v.Set(key, value)
	}

	err := v.WriteConfig()
	if err != nil {
		if os.IsNotExist(err) {
			err = v.SafeWriteConfig()
		}
	}

	return err
}
