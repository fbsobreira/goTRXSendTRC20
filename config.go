package main

import (
	"./util"
	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v2"
)

// APIConfig configuration
type APIConfig struct {
	FullNode string `yaml:"fullNode"`
}

// Config defines the config for server
type Config struct {
	API        APIConfig `yaml:"api"`
	PrivateKey string    `yaml:"privateKey"`
}

func loadConfig(configPath string) Config {
	data, err := util.ReadFile(configPath)
	if err != nil {
		zap.L().Fatal("failed to load config file", zap.Error(err))
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		zap.L().Fatal("failed to unmarshal config", zap.Error(err))
	}
	return config
}
