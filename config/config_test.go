package config

import (
	"os"
	"path/filepath"
	"testing"

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
