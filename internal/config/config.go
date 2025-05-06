package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Rest struct {
		Port         int           `env:"REST_PORT" env-default:"8080"`
		ReadTimeout  time.Duration `env:"REST_READ_TIMEOUT" env-default:"20s"`
		WriteTimeout time.Duration `env:"REST_WRITE_TIMEOUT" env-default:"20s"`
	}

	Worker struct {
		DownloadWorkerCount int `env:"DOWNLOAD_WORKER_COUNT" env-default:"3"`
		ScanWorkerCount     int `env:"SCAN_WORKER_COUNT" env-default:"2"`
		QueueSize           int `env:"QUEUE_SIZE" env-default:"100"`
	}

	Redis struct {
		Address  string `env:"REDIS_ADDRESS" env-default:"127.0.0.1"`
		Port     int    `env:"REDIS_PORT" env-default:"6379"`
		Password string `env:"REDIS_PASSWORD" env-default:""`
		Database int    `env:"REDIS_DB" env-default:"0"`
	}

	PGSQL struct {
		Address        string `env:"MYSQL_ADDRESS" env-default:"127.0.0.1"`
		Port           int    `env:"MYSQL_PORT" env-default:"3306"`
		User           string `env:"MYSQL_USER" env-default:"root"`
		Password       string `env:"MYSQL_PASSWORD" env-default:"990824"`
		ConnectTimeout int    `env:"MYSQL_CONNECT_TIMEOUT" env-default:"5"`
		ResultDb       string `env:"MYSQL_RESULT_DB" env-default:"mydb"`
	}
}

type ScanRequest struct {
	ImageName      string   `json:"imageName"`
	ScanStartTime  string   `json:"scanStartTime"`
	ScanFinishTime string   `json:"scanFinishTime"`
	ScanResult     string   `json:"scanResult"`
	MaliciousFiles []string `json:"maliciousFiles"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
