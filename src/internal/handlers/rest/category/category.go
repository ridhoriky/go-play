package category

import (
	"ne-project/src/internal/services/category"

	"github.com/gin-gonic/gin"
)

type CategoryHandlerItf interface {
	RegisterRoutes(r *gin.Engine)
}

type categoryHandler struct {
	gin             *gin.Engine
	categoryService category.CategoryServiceItf
}

func NewCategoryHandler(categoryService category.CategoryServiceItf) CategoryHandlerItf {
	return &categoryHandler{
		categoryService: categoryService,
	}
}

func (h *categoryHandler) RegisterRoutes(r *gin.Engine) {
	categoryRoutes := r.Group("/categories")
	{
		categoryRoutes.GET("", h.GetAll)
		categoryRoutes.POST("", h.Create)
		categoryRoutes.GET("/:id", h.GetByID)
		categoryRoutes.PUT("/:id", h.Update)
		categoryRoutes.DELETE("/:id", h.Delete)
	}
}
