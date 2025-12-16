package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Server struct {
		GRPCAddr string `yaml:"grpc_addr"`
	} `yaml:"server"`
	DB struct {
		DSN string `yaml:"dsn"`
	} `yaml:"db"`
	Clients struct {
		AuthServiceAddr string `yaml:"auth_service_addr"`
	} `yaml:"clients"`
}

func Load(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
