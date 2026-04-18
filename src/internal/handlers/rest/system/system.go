package system

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type SystemHandlerItf interface {
	HealthCheck(c *gin.Context)
}

type systemHandler struct {
	db *sqlx.DB
}

func NewSystemHandler(db *sqlx.DB) *systemHandler {
	return &systemHandler{db: db}
}

func (h *systemHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/health", h.HealthCheck)
	r.GET("/ready", h.ReadyCheck)
}
