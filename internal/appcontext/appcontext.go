package appcontext

import (
	"github.com/PAM-IDAM-Org/asset-discovery/internal/config"
	"github.com/PAM-IDAM-Org/asset-discovery/internal/storage"
	"github.com/PAM-IDAM-Org/asset-discovery/pkg/logger"
)

type AppContext struct {
	Log     *logger.Logger
	Postgre *config.Postgre
	Env     *config.AppEnv
	Redis   *storage.RedisClient
}
