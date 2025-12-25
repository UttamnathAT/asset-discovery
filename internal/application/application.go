package application

import (
	"context"
	"fmt"

	"github.com/PAM-IDAM-Org/asset-discovery/internal/appcontext"
	"github.com/Uttamnath64/arvo-fin/app/config"
	"github.com/Uttamnath64/arvo-fin/app/requests"
	"github.com/Uttamnath64/arvo-fin/app/storage"
	"github.com/Uttamnath64/arvo-fin/fin-api/internal/routes"
	"github.com/Uttamnath64/arvo-fin/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	requests.NewResponse()

	// Load env
	_ = godotenv.Load()
	env, err := config.LoadEnv()
	if err != nil {
		logger.New("none", nil).Error("Failed to load .env", err.Error())
		return false
	}

	// Set logger
	log := logger.New(env.Server.Environment, nil)

	// Load config
	var con config.Config
	err = config.LoadConfig(env, &con)
	if err != nil {
		log.Error("Failed to load config", err.Error())
		return false
	}

	// Load Redis
	redis, err := storage.NewRedisClient(ctx, env.Redis.Addr, env.Redis.Password, env.Redis.DB)
	if err != nil {
		log.Error("Failed to load redis", err.Error())
		return false
	}

	// DI container
	a.log = log
	a.appContext = &appcontext.AppContext{
		Log:     log,
		Postgre: con.Postgre,
		Env:     &env,
		Redis:   redis,
	}

	return true

}

func (a *Application) Run() {

	server := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{a.appContext.Env.Server.ClientOrigin}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	server.Use(cors.New(corsConfig))

	// Setup routes
	routes.New(a.appContext, server).Handlers()

	// Run server
	if err := server.Run(fmt.Sprintf(":%d", a.appContext.Env.Server.HTTPPort)); err != nil {
		a.log.Error("api-application-server", err.Error())
		return
	}
}
