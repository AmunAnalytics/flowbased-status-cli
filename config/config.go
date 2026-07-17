package config

import (
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

var NewVersionAvailable = false

type Config struct {
	General  GeneralConfig  `toml:"General"`
	FBStatus FBStatusConfig `toml:"fbstatus"`
}

type GeneralConfig struct {
	DisableTelemetry bool `toml:"disable_telemetry"`
}

type FBStatusConfig struct {
	Debug                bool `toml:"debug_backend"`
	SuppressVersionCheck bool `toml:"surpress_version_check"`
}

func GetConfigFolder() string {
	home, _ := os.UserHomeDir()
	pathFolder := filepath.Join(home, ".amun-analytics")
	os.MkdirAll(pathFolder, os.ModePerm)

	return pathFolder
}

func GetConfig() Config {
	pathFolder := GetConfigFolder()
	path := filepath.Join(pathFolder, "config.toml")
	config := Config{}

	data, err := os.ReadFile(path)
	if err == nil {
		toml.Unmarshal(data, &config)
	}

	if val := os.Getenv("AMUN_DISABLE_TELEMETRY"); val == "1" {
		config.General.DisableTelemetry = true
	}

	if val := os.Getenv("AMUN_DEBUG"); val == "1" {
		config.FBStatus.Debug = true
	}

	if val := os.Getenv("AMUN_SURPRESS_VERSION_CHECK"); val == "1" {
		config.FBStatus.SuppressVersionCheck = true
	}

	return config
}
