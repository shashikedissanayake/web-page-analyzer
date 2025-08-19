package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
)

type Configuration struct {
	Port string `env:"PORT" envDefault:":8080"`
}

func NewConfig(files ...string) (*Configuration, error) {
	err := godotenv.Load(files...)

	if err != nil {
		logger.Error("Failed to load .env file")
	}

	cfg := Configuration{}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
