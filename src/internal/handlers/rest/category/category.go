package category

import (
	"ne-project/src/internal/services/category"

	"github.com/gin-gonic/gin"
)

type CategoryHandlerItf interface {
	GetAll(c *gin.Context)
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type categoryHandler struct {
	categoryService category.CategoryServiceItf
}

func NewCategoryHandler(categoryService category.CategoryServiceItf) CategoryHandlerItf {
	return &categoryHandler{
		categoryService: categoryService,
	}
}
