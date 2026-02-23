package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App  AppConfig      `yaml:"app"`
	DB   DatabaseConfig `yaml:"db"`
	Auth JWTConfig      `yaml:"auth"`
}

type AppConfig struct {
	Env         string        `yaml:"env"`
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `env:"DB_NAME"`
	User     string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	SSLMode  string `yaml:"sslmode"`
	MaxConns int    `yaml:"max_conns"`
	MinConns int    `yaml:"min_conns"`
}

type JWTConfig struct {
	Secret   []byte        `env:"JWT_SECRET"`
	TokenTTL time.Duration `yaml:"tokenttl"`
}

func MustLoad(configPath string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}

	//Sensible database data inside env file
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.User = os.Getenv("DB_USERNAME")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	cfg.Auth.Secret = []byte(os.Getenv("JWT_SECRET"))

	return &cfg, nil
}

func (db *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host,
		db.Port,
		db.User,
		db.Password,
		db.Name,
		db.SSLMode,
	)
}
