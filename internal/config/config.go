package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App AppConfig
}

type AppConfig struct {
	Env         string
	Host        string
	Port        string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func MustLoad(cfgPath string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
