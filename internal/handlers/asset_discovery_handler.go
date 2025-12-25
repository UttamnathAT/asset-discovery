package handlers

import (
	"context"

	"github.com/PAM-IDAM-Org/asset-discovery/api/gen"
	healthpb "github.com/PAM-IDAM-Org/asset-discovery/api/gen/health"
	"github.com/PAM-IDAM-Org/asset-discovery/internal/appcontext"
)

type AssetDiscoveryHandler struct {
	gen.UnimplementedAssetDiscoveryServiceServer // grpc interface
	appContext                                   *appcontext.AppContext
	healthHandler                                *HealthHandler // health interface
}

func NewAssetDiscoveryHandler(appContext *appcontext.AppContext) *AssetDiscoveryHandler {
	return &AssetDiscoveryHandler{
		appContext:    appContext,
		healthHandler: NewHealthHandler(appContext),
	}
}

func (h *AssetDiscoveryHandler) Health(ctx context.Context, req *healthpb.HealthRequest) (*healthpb.HealthResponse, error) {
	return h.healthHandler.Health(ctx, req)
}
