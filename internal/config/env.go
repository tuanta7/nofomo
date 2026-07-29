package config

import "github.com/joho/godotenv"

type EnvConfig struct {
	APIKey    string
	APISecret string
}

func NewEnvConfig() (EnvConfig, error) {
	var envConfig EnvConfig

	if err := godotenv.Load(); err != nil {
		return envConfig, err
	}

	return envConfig, nil
}
