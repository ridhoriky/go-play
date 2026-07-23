package router

import (
	_ "ne-project/docs"
	"ne-project/src/internal/config/appconfig"
	"ne-project/src/internal/config/resource"
	"ne-project/src/internal/config/token"
	"ne-project/src/internal/handlers/rest/middleware"
	"ne-project/src/internal/handlers/rest/routes"
	"ne-project/src/internal/repositories"
	"ne-project/src/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func SetupRouter(log *zerolog.Logger, cfg *appconfig.Config, res *resource.Resources, tokenSvc token.TokenServiceItf) *gin.Engine {

	repo := repositories.NewRepository(res.DB, res.Redis)
	service := services.NewServices(repo, tokenSvc, res.Redis, cfg)
	mw := middleware.InitMiddleware(log, tokenSvc, &cfg.RateLimit, service.Store)
	handlers := routes.NewHandlers(res.DB, service, cfg)

	r := gin.New()
	r.Use(mw.Logger(), mw.CORS(), mw.Recovery(), otelgin.Middleware("kasir-app"))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	handlers.RegisterRoutes(v1, tokenSvc, mw)

	return r
}
