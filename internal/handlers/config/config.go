package config

import (
	"github.com/PAM-IDAM-Org/asset-discovery/pkg/logger"
)

type Config struct {
	Postgre *Postgre
}

func LoadConfig(env *AppEnv, log *logger.Logger) (*Config, error) {
	// Initialize PostgreSQL connections - skip for health check only
	// postgre, err := InitPostgre(env)
	// if err != nil {
	// 	return nil, err
	// }

	return &Config{
		Postgre: nil, // Skip DB for health check
	}, nil
}
