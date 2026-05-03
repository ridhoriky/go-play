package rest

import (
	"ne-project/src/internal/models/dto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type SystemHandler struct {
        db  *sqlx.DB
}

func NewSystemHandler(db *sqlx.DB) *SystemHandler {
        return &SystemHandler{
                db:  db,
        }
}
func (h *SystemHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/health", h.HealthCheck)
	r.GET("/ready", h.ReadyCheck)
}

// HealthCheck godoc
// @Summary      Liveness Probe
// @Description  Check if the application is alive (lightweight check)
// @Tags         System
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.HealthCheckResponse
// @Router       /health [get]
func (h *SystemHandler) HealthCheck(c *gin.Context) {
	// Liveness: hanya cek proses masih jalan
	// Jangan cek DB atau dependency lain di sini!
	// Biarkan tetap return 200 meskipun DB mati (biar K8s tidak restart pod)

	health := dto.HealthCheckResponse{
		Status:  "alive",
		Time:    time.Now().Format(time.RFC3339),
		Service: "go-play",
		Version: "1.0.0",
	}

	c.JSON(http.StatusOK, health)
}

// ReadyCheck godoc
// @Summary      Readiness Probe
// @Description  Check if the application is ready to accept traffic
// @Tags         System
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.ReadyCheckResponse
// @Failure      503  {object}  dto.ErrorResponse  // 503 Service Unavailable
// @Router       /ready [get]   // ← ganti dari /config ke /ready
func (h *SystemHandler) ReadyCheck(c *gin.Context) {
	ctx := c.Request.Context()

	status := "ready"
	httpStatus := http.StatusOK
	dependencies := make(map[string]string)

	if err := h.db.PingContext(ctx); err != nil {
		status = "not_ready"
		httpStatus = http.StatusServiceUnavailable
		dependencies["database"] = "unavailable"
	} else {
		dependencies["database"] = "available"
	}

	// TODO: Check Redis/Cache
	// if err := h.cache.Ping(ctx); err != nil {
	//     status = "not_ready"
	//     dependencies["cache"] = "unavailable"
	// }

	// TODO: Check other dependencies (MQ, external API, dll)

	ready := dto.ReadyCheckResponse{
		Status:       status,
		Time:         time.Now().Format(time.RFC3339),
		Service:      "go-play",
		Version:      "1.0.0",
		Dependencies: dependencies,
		Message:      map[bool]string{true: "Ready to serve requests", false: "Not ready"}[httpStatus == http.StatusOK],
	}

	c.JSON(httpStatus, ready)
}
