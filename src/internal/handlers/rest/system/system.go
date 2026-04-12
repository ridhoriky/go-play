package system

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type SystemHandlerItf interface {
	RegisterRoutes(r *gin.Engine)
}

type systemHandler struct {
	db *sqlx.DB
}

func NewSystemHandler(db *sqlx.DB) *systemHandler {
	return &systemHandler{db: db}
}

func (h *systemHandler) RegisterRoutes(r *gin.Engine) {
	systemRoutes := r.Group("/api/v1")
	{
		systemRoutes.GET("/health", h.HealthCheck)
		systemRoutes.GET("/ready", h.ReadyCheck)
	}
}
