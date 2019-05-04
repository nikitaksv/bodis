package config

import (
	"bytes"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

// GetConfig returns application config
func GetConfig() *TomlConfig {
	return configInstance
}

// NewConfig creates new application config with given .toml file
func NewConfig(file string) (*TomlConfig, error) {
	configInstance = &TomlConfig{}

	if _, err := toml.DecodeFile(file, configInstance); err != nil {
		return nil, err
	}
	dump(configInstance)

	return configInstance, nil
}

func dump(cfg *TomlConfig) {
	var buffer bytes.Buffer
	e := toml.NewEncoder(&buffer)
	err := e.Encode(cfg)
	if err != nil {
	}

	fmt.Println(
		time.Now().UTC(),
		"\n---------------------Bodis started with config:\n",
		buffer.String(),
		"\n---------------------")
}
