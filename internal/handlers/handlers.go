package handlers

import (
	"net/http"
	"time"

	"github.com/bdrbt/todo/internal/config"
	"github.com/bdrbt/todo/internal/usecases"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

const (
	HTTPPrefix = "/api/v1/"
)

func MountHandlers(cfg *config.Config, uc *usecases.UC, logger *zap.Logger) http.Handler {
	if cfg.Environment != config.EnvDevelopment {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	api := router.Group(HTTPPrefix)

	// Attach swagger

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Auth paths
	authRoutes := api.Group("/auth")
	{
		authRoutes.GET(auth.PrefixRequestOTP, auth.RequestOTP(logger, uc))
		authRoutes.POST(auth.PrefixClaimOTP, auth.ClaimOTP(logger, uc))
	}

	// User routes
	userRoutes := api.Group("/user")
	userRoutes.Use(token.Middleware())
	{
		userRoutes.GET(user.PrefixProfile, user.Profile(logger, uc))
	}

	for _, item := range router.Routes() {
		println("method:", item.Method, "path:", item.Path)
	}

	return router
}
