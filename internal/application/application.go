package application

import (
	"context"

	"github.com/PAM-IDAM-Org/asset-discovery/internal/appcontext"
	"github.com/PAM-IDAM-Org/asset-discovery/internal/config"
	"github.com/PAM-IDAM-Org/asset-discovery/internal/storage"
	"github.com/PAM-IDAM-Org/asset-discovery/pkg/logger"
	"github.com/joho/godotenv"
)

type Application struct {
	appContext *appcontext.AppContext
	log        *logger.Logger
}

func New() *Application {
	return &Application{}
}

func (a *Application) Initialize() bool {
	ctx := context.Background()

	// Load env
	_ = godotenv.Load()
	env, err := config.LoadEnv()
	if err != nil {
		logger.New("none", nil).Error("Failed to load .env", "error", err.Error())
		return false
	}

	// Set logger
	log := logger.New(env.Server.Environment, nil)

	// Load config (PostgreSQL) - optional for health check
	var cfg *config.Config
	cfg, err = config.LoadConfig(&env, log)
	if err != nil {
		log.Warn("PostgreSQL not available, continuing without DB", "error", err.Error())
		cfg = nil
	}

	// Load Redis - optional for health check
	var redis *storage.RedisClient
	redis, err = storage.NewRedisClient(ctx, env.Server.RedisAddr, "", env.Server.RedisDB)
	if err != nil {
		log.Warn("Redis not available, continuing without Redis", "error", err.Error())
		redis = nil
	}

	// DI container
	a.log = log
	var postgre *config.Postgre
	if cfg != nil {
		postgre = cfg.Postgre
	}
	a.appContext = &appcontext.AppContext{
		Log:     log,
		Postgre: postgre,
		Env:     &env,
		Redis:   redis,
	}

	return true
}

func (a *Application) Run() error {

	// Start gRPC server
	if err := a.StartGRPCServer(); err != nil {
		return err
	}

	// Start gRPC-Gateway server
	if err := a.StartGRPCGateway(); err != nil {
		return err
	}

	return nil
}
