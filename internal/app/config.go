package app

import (
	"context"
	"log/slog"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DBHost     string `env:"DB_HOST,default=localhost"`
	DBPort     string `env:"DB_PORT,default=3306"`
	DBUser     string `env:"DB_USER,default=root"`
	DBPassword string `env:"DB_PASSWORD,default="`
	DBName     string `env:"DB_NAME,default=sw_config"`
	ServerAddr string `env:"SERVER_ADDR,default=:8080"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}

	slog.Info("configuration loaded",
		"db_host", config.DBHost,
		"db_port", config.DBPort,
		"db_name", config.DBName,
		"server_addr", config.ServerAddr)

	return &config, nil
}
