package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	PgConnStr string `env:"POSTGRES_CONN_STR,required"`
	R4uabURL  string `env:"R4UAB_URL,required"`
	HTTPPort  int    `env:"HTTP_PORT,required"`
}

func New() (*Config, error) {
	if godotenv.Load() == nil {
		log.Info().Msg("env variables loaded from .env file")
	}

	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
