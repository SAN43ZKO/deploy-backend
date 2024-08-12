package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP     HTTP
	Postgres Postgres
	JWT      JWT
	Steam    Steam
	Swagger  Swagger
}

type HTTP struct {
	Host         string        `env:"HTTP_HOST"`
	Port         string        `env:"HTTP_PORT"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT"`
}

type Postgres struct {
	DSN string `env:"PG_DSN"`
}

type JWT struct {
	Key string `env:"JWT_KEY"`
}

type Steam struct {
	APIKey string `env:"STEAM_API_KEY"`
}

type Swagger struct {
	URL string `env:"SWAGGER_URL"`
}

func Init() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("Init: %w", err)
	}

	return &cfg, nil
}
