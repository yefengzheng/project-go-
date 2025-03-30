package config

import "time"
import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Rest struct {
		Port         int           `env:"REST_PORT" env-default:"8080"`
		ReadTimeout  time.Duration `env:"REST_READ_TIMEOUT" env-default:"20s"`
		WriteTimeout time.Duration `env:"REST_WRITE_TIMEOUT" env-default:"20s"`
	}
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
