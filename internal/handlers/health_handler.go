package handlers

import (
	"context"
	"fmt"
	"time"

	healthpb "github.com/PAM-IDAM-Org/asset-discovery/api/gen/health"
	"github.com/PAM-IDAM-Org/asset-discovery/internal/appcontext"
)

type HealthHandler struct {
	appContext *appcontext.AppContext
}

func NewHealthHandler(appContext *appcontext.AppContext) *HealthHandler {
	return &HealthHandler{
		appContext: appContext,
	}
}

func (h *HealthHandler) Health(ctx context.Context, req *healthpb.HealthRequest) (*healthpb.HealthResponse, error) {
	// Get current timestamp
	timestamp := time.Now().Format(time.RFC3339)

	// Check database status
	dbStatus := "healthy"
	dbMessage := "Database connection is healthy"
	if h.appContext.Postgre != nil && h.appContext.Postgre.ReadWriteDB != nil {
		sqlDB, err := h.appContext.Postgre.ReadWriteDB.DB()
		if err == nil {
			if err := sqlDB.PingContext(ctx); err != nil {
				dbStatus = "unhealthy"
				dbMessage = fmt.Sprintf("Database ping failed: %v", err)
				h.appContext.Log.Warn("Database health check failed", "error", err.Error())
			}
		} else {
			dbStatus = "unhealthy"
			dbMessage = fmt.Sprintf("Database connection error: %v", err)
			h.appContext.Log.Warn("Database connection check failed", "error", err.Error())
		}
	} else {
		dbStatus = "unavailable"
		dbMessage = "Database not configured"
	}

	// Check Redis status
	redisStatus := "healthy"
	redisMessage := "Redis connection is healthy"
	if h.appContext.Redis != nil && h.appContext.Redis.Client != nil {
		_, err := h.appContext.Redis.Client.Ping(ctx).Result()
		if err != nil {
			redisStatus = "unhealthy"
			redisMessage = fmt.Sprintf("Redis ping failed: %v", err)
			h.appContext.Log.Warn("Redis health check failed", "error", err.Error())
		}
	} else {
		redisStatus = "unavailable"
		redisMessage = "Redis not configured"
	}

	// Determine overall status and message
	overallStatus := "healthy"
	overallMessage := "All systems operational"

	unhealthyCount := 0
	unavailableCount := 0

	if dbStatus == "unhealthy" {
		unhealthyCount++
	} else if dbStatus == "unavailable" {
		unavailableCount++
	}

	if redisStatus == "unhealthy" {
		unhealthyCount++
	} else if redisStatus == "unavailable" {
		unavailableCount++
	}

	if unhealthyCount > 0 {
		overallStatus = "degraded"
		overallMessage = fmt.Sprintf("Service is running but %d component(s) are unhealthy", unhealthyCount)
	}

	if unhealthyCount == 0 && unavailableCount == 2 {
		overallStatus = "healthy"
		overallMessage = "Service is healthy (dependencies not configured)"
	}

	// Get service version
	version := h.appContext.Env.Server.Environment
	if version == "" {
		version = "1.0.0"
	}

	// Get service name
	serviceName := "asset-discovery"

	// Get bind address
	bindAddress := fmt.Sprintf(":%d", h.appContext.Env.Server.GRPCPort)

	// Build response
	response := &healthpb.HealthResponse{
		Status:  overallStatus,
		Message: overallMessage,
		Details: &healthpb.HealthDetails{
			DatabaseStatus:  dbStatus,
			DatabaseMessage: dbMessage,
			RedisStatus:     redisStatus,
			RedisMessage:    redisMessage,
		},
		Metadata: &healthpb.HealthMetadata{
			Timestamp:          timestamp,
			Version:            version,
			Service:            serviceName,
			BindAddress:        bindAddress,
			SupportedProtocols: []string{"grpc", "http"},
		},
	}

	return response, nil
}
