package config

import (
	"os"
	"time"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Server struct {
		GRPCAddr string `yaml:"grpc_addr"`
	} `yaml:"server"`
	DB struct {
		DSN string `yaml:"dsn"`
	} `yaml:"db"`
	JWT struct {
		Secret     string `yaml:"secret"`
		AccessTTL  time.Duration `yaml:"access_ttl"`
		RefreshTTL time.Duration `yaml:"refresh_ttl"`
		Iss        string `yaml:"iss"`
		Aud        string `yaml:"aud"`
	} `yaml:"jwt"`
	Clients struct {
		UserServiceAddr string `yaml:"user_service_addr"`
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
