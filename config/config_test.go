package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestWriteConfigToFile(t *testing.T) {
	tmpPath, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpPath)

	cfg := DefaultLoggerConfig()
	fileNames := []string{"rainbowlog.toml", "rainbowlog.yaml", "rainbowlog.json"}

	for _, name := range fileNames {
		cfgFile := filepath.Join(tmpPath, name)
		err = WriteConfigToFile(cfgFile, cfg)
		require.NoError(t, err)
		_, err = os.Stat(cfgFile)
		require.NoError(t, err)
	}
}

func TestLoadLoggerConfigFromFile(t *testing.T) {
	tmpPath, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpPath)

	cfg := DefaultLoggerConfig()
	fileNames := []string{"rainbowlog.toml", "rainbowlog.yaml", "rainbowlog.json"}
	for _, fileName := range fileNames {
		cfgFile := filepath.Join(tmpPath, fileName)
		err = WriteConfigToFile(cfgFile, cfg)
		require.NoError(t, err)

		cfg2, err := LoadLoggerConfigFromFile(cfgFile)
		require.NoError(t, err)

		require.Equal(t, cfg.Enable, cfg2.Enable)
		require.Equal(t, cfg.Level, cfg2.Level)
		require.Equal(t, cfg.Label, cfg2.Label)
		require.Equal(t, cfg.Stack, cfg2.Stack)
		require.Equal(t, cfg.EnableConsolePrinting, cfg2.EnableConsolePrinting)
		require.Equal(t, cfg.EnableRainbowConsole, cfg2.EnableRainbowConsole)
		require.Equal(t, cfg.TimeFormat, cfg2.TimeFormat)
		require.Equal(t, cfg.SizeRollingFileConfig.Enable, cfg2.SizeRollingFileConfig.Enable)
		require.Equal(t, cfg.SizeRollingFileConfig.LogFilePath, cfg2.SizeRollingFileConfig.LogFilePath)
		require.Equal(t, cfg.SizeRollingFileConfig.LogFileBaseName, cfg2.SizeRollingFileConfig.LogFileBaseName)
		require.Equal(t, cfg.SizeRollingFileConfig.MaxBackups, cfg2.SizeRollingFileConfig.MaxBackups)
		require.Equal(t, cfg.SizeRollingFileConfig.FileSizeLimit, cfg2.SizeRollingFileConfig.FileSizeLimit)
		require.Equal(t, cfg.SizeRollingFileConfig.UseBufferedWriter, cfg2.SizeRollingFileConfig.UseBufferedWriter)
		require.Equal(t, cfg.SizeRollingFileConfig.WriterBufferSize, cfg2.SizeRollingFileConfig.WriterBufferSize)
		require.Equal(t, cfg.TimeRollingFileConfig.Enable, cfg2.TimeRollingFileConfig.Enable)
		require.Equal(t, cfg.TimeRollingFileConfig.LogFilePath, cfg2.TimeRollingFileConfig.LogFilePath)
		require.Equal(t, cfg.TimeRollingFileConfig.LogFileBaseName, cfg2.TimeRollingFileConfig.LogFileBaseName)
		require.Equal(t, cfg.TimeRollingFileConfig.MaxBackups, cfg2.TimeRollingFileConfig.MaxBackups)
		require.Equal(t, cfg.TimeRollingFileConfig.RollingPeriod, cfg2.TimeRollingFileConfig.RollingPeriod)
		require.Equal(t, cfg.TimeRollingFileConfig.UseBufferedWriter, cfg2.TimeRollingFileConfig.UseBufferedWriter)
		require.Equal(t, cfg.TimeRollingFileConfig.WriterBufferSize, cfg2.TimeRollingFileConfig.WriterBufferSize)
	}
}

func TestLoadLoggerConfigFromEnv(t *testing.T) {
	t.Setenv("RAINBOWLOG_ENABLE", "false")
	t.Setenv("RAINBOWLOG_LEVEL", "ERROR")
	t.Setenv("RAINBOWLOG_LABEL", "env-label")
	t.Setenv("RAINBOWLOG_STACK", "true")
	t.Setenv("RAINBOWLOG_ENABLECONSOLEPRINTING", "false")
	t.Setenv("RAINBOWLOG_ENABLERAINBOWCONSOLE", "false")
	t.Setenv("RAINBOWLOG_TIMEFORMAT", "UNIX")
	t.Setenv("RAINBOWLOG_SIZEROLLINGFILECONFIG_ENABLE", "true")
	t.Setenv("RAINBOWLOG_SIZEROLLINGFILECONFIG_LOGFILEPATH", "./env-log")

	cfg, err := LoadLoggerConfigFromEnv("")
	require.NoError(t, err)
	require.False(t, cfg.Enable)
	require.Equal(t, "ERROR", cfg.Level)
	require.Equal(t, "env-label", cfg.Label)
	require.True(t, cfg.Stack)
	require.False(t, cfg.EnableConsolePrinting)
	require.False(t, cfg.EnableRainbowConsole)
	require.Equal(t, "UNIX", cfg.TimeFormat)
	require.True(t, cfg.SizeRollingFileConfig.Enable)
	require.Equal(t, "./env-log", cfg.SizeRollingFileConfig.LogFilePath)
}

func TestLoadLoggerConfigFromEnv_CustomPrefix(t *testing.T) {
	t.Setenv("MYAPP_LEVEL", "WARN")
	t.Setenv("MYAPP_LABEL", "custom-label")

	cfg, err := LoadLoggerConfigFromEnv("MYAPP")
	require.NoError(t, err)
	require.Equal(t, "WARN", cfg.Level)
	require.Equal(t, "custom-label", cfg.Label)
}

func TestLoadLoggerConfigFromViper_EnvOverridesFile(t *testing.T) {
	v := viper.New()
	v.SetConfigType("yaml")
	require.NoError(t, v.ReadConfig(strings.NewReader(`
rainbowlog:
  level: INFO
  label: file-label
`)))

	t.Setenv("RAINBOWLOG_LEVEL", "ERROR")

	cfg, err := LoadLoggerConfigFromViper(v)
	require.NoError(t, err)
	require.Equal(t, "ERROR", cfg.Level, "environment variable should override the file value")
	require.Equal(t, "file-label", cfg.Label, "field absent in env should keep the file value")
}
