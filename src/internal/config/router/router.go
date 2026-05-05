package router

import (
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
)

func SetupRouter(log *zerolog.Logger, cfg *appconfig.Config, res *resource.Resources, tokenSvc *token.Token) *gin.Engine {
	mw := middleware.InitMiddleware(log, tokenSvc, &cfg.RateLimit)
	repo := repositories.NewRepository(res.DB)
	service := services.NewServices(repo, tokenSvc, res.DB, res.Redis)
	handlers := routes.NewHandlers(res.DB, service)

	r := gin.New()
	r.Use(mw.Logger(), mw.ErrorHandler(), mw.CORS(), mw.Recovery())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	handlers.RegisterRoutes(v1, tokenSvc, *mw)

	return r
}
