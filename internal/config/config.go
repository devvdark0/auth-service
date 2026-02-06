package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App AppConfig `yaml:"app"`
	DB  DBConfig  `yaml:"db"`
}

type AppConfig struct {
	Env             string        `yaml:"env"`
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	Timeout         time.Duration `yaml:"timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DBConfig struct {
	URL string `env:"DATABASE_URL" env_required:"true"`
}

func MustLoad(path string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	cfg.DB.URL = os.Getenv("DATABASE_URL")

	return &cfg, nil
}
