package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	App  AppConfig  `yaml:"app"`
	Db   DbConfig   `yaml:"db"`
	Auth AuthConfig `yaml:"auth"`
}

type AppConfig struct {
	Env         string        `yaml:"env"`
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type DbConfig struct {
	Name            string        `yaml:"name"`
	Host            string        `yaml:"host"`
	Port            string        `yaml:"port"`
	Username        string        `yaml:"user"`
	Password        string        `env:"DB_PASSWORD"`
	MaxIdleConn     int           `yaml:"max_idle_conn"`
	MaxOpenConn     int           `yaml:"max_open_conn"`
	MaxLifetimeConn time.Duration `yaml:"max_conn_time_sec"`
	SSLMode         string        `yaml:"sslmode"`
}

type AuthConfig struct {
	Secret   []byte        `env:"JWT_SECRET"`
	TokenTTL time.Duration `yaml:"token-ttl"`
}

func MustLoad(cfgPath string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, err
	}

	cfg.Db.Password = os.Getenv("DB_PASSWORD")
	cfg.Auth.Secret = []byte(os.Getenv("JWT_SECRET"))

	return &cfg, nil
}

func (d *DbConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		d.Host,
		d.Port,
		d.Name,
		d.Username,
		d.Password,
		d.SSLMode,
	)
}
