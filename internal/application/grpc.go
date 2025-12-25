package application

import (
	"fmt"
	"net"

	"github.com/PAM-IDAM-Org/asset-discovery/api/gen"
	"github.com/PAM-IDAM-Org/asset-discovery/internal/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (a *Application) StartGRPCServer() error {
	grpcServer := grpc.NewServer()

	assetDiscoveryHandler := handlers.NewAssetDiscoveryHandler(a.appContext)
	gen.RegisterAssetDiscoveryServiceServer(grpcServer, assetDiscoveryHandler)

	// Enable reflection for testing
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.appContext.Env.Server.GRPCPort))
	if err != nil {
		a.log.Error("Failed to listen on gRPC port", "error", err.Error())
		return err
	}

	go func() {
		a.log.Info("gRPC server starting", "port", a.appContext.Env.Server.GRPCPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			a.log.Error("gRPC server failed", "error", err.Error())
		}
	}()

	return nil
}
