package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP HTTPConfig `env-prefix:"HTTP_"`
	App  AppConfig  `env-prefix:"APP_"`
}

type HTTPConfig struct {
	Port              string        `env:"PORT" env-default:"8081"`
	Address           string        `env:"ADDRESS" env-default:"0.0.0.0"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" env-default:"10s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT" env-default:"10s"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" env-default:"5s"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}

type AppConfig struct {
	Env string `env:"ENV" env-default:"prod"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
