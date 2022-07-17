package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type (
	ZenithConfig struct {
		App       configApp
		Functions []configFunction
	}

	configFunction struct {
		Name  string
		Path  string
		Route string
	}

	configApp struct {
		Name    string
		Runtime string
	}
)

func ReadConfig(configPath string) (ZenithConfig, error) {
	var cnf ZenithConfig

	if _, err := os.Stat(configPath); err != nil {
		return ZenithConfig{}, err
	}

	if _, err := toml.DecodeFile(configPath, &cnf); err != nil {
		return ZenithConfig{}, err
	}

	return cnf, nil
}
