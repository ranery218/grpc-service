package config

import (
	"os"
	"time"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Server struct {
		HTTPAddr     string `yaml:"http_addr"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
	} `yaml:"server"`
	Backends struct {
		AuthServiceAddr string `yaml:"auth_service_addr"`
		UserServiceAddr string `yaml:"user_service_addr"`
	} `yaml:"backends"`
	JWT struct {
		Secret string `yaml:"secret"`
		Iss    string `yaml:"iss"`
		Aud    string `yaml:"aud"`
	} `yaml:"jwt"`
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
