package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	PostgresURL string `envconfig:"POSTGRES_URL" required:"true"`
	JWTSecret   string `envconfig:"JWT_SECRET" required:"true"`
	JWTExpiry   int    `envconfig:"JWT_EXPIRY_MINUTES" default:"60"`
	Port        string `envconfig:"PORT" default:"8080"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("MARKETPLACE", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
