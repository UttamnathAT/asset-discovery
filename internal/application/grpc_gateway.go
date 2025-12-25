package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PAM-IDAM-Org/asset-discovery/api/gen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (a *Application) StartGRPCGateway() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := gen.RegisterAssetDiscoveryServiceHandlerFromEndpoint(
		ctx, mux, fmt.Sprintf("localhost:%d", a.appContext.Env.Server.GRPCPort), opts)
	if err != nil {
		a.log.Error("Failed to register asset discovery service gateway", "error", err.Error())
		return err
	}

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{a.appContext.Env.Server.ClientOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := corsConfig.Handler(mux)

	httpAddr := fmt.Sprintf(":%d", a.appContext.Env.Server.HTTPPort)
	a.log.Info("HTTP gateway server starting", "port", a.appContext.Env.Server.HTTPPort)

	if err := http.ListenAndServe(httpAddr, handler); err != nil {
		a.log.Error("HTTP gateway server failed", "error", err.Error())
		return err
	}

	return nil
}
