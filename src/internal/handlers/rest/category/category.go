package category

import (
	"ne-project/src/internal/services/category"

	"github.com/gin-gonic/gin"
)

type CategoryHandlerItf interface {
	RegisterRoutes(r *gin.RouterGroup)
}

type categoryHandler struct {
	gin             *gin.RouterGroup
	categoryService category.CategoryServiceItf
}

func NewCategoryHandler(categoryService category.CategoryServiceItf) CategoryHandlerItf {
	return &categoryHandler{
		categoryService: categoryService,
	}
}

func (h *categoryHandler) RegisterRoutes(r *gin.RouterGroup) {
	categoryRoutes := r.Group("/categories")
	{
		categoryRoutes.GET("", h.GetAll)
		categoryRoutes.POST("", h.Create)
		categoryRoutes.GET("/:id", h.GetByID)
		categoryRoutes.PUT("/:id", h.Update)
		categoryRoutes.DELETE("/:id", h.Delete)
	}
}
