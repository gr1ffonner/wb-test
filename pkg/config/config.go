package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	NATS     NATSConfig
	Logger   Logger
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT" env-default:"8080"`
}

type Logger struct {
	Level string `env:"LOG_LEVEL" env-default:"info"`
}

type DatabaseConfig struct {
	DSN string `env:"DB_DSN" env-required:"true"`
}

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR" env-default:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `env:"REDIS_DB" env-default:"0"`
}

type NATSConfig struct {
	URL string `env:"NATS_URL" env-default:"nats://localhost:4222"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	log.Println("Config loaded:", cfg)

	return cfg, nil
}
