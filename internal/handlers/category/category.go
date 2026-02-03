package category

import (
	"ne-project/internal/services/category"
	"net/http"
)


type CategoryHandlerItf interface {
	HandleCategories( w http.ResponseWriter, r *http.Request)
	HandleCategoryByID( w http.ResponseWriter, r *http.Request)
}

type categoryHandler struct {
	categoryService category.CategoryServiceItf
}

func NewCategoryHandler(categoryService category.CategoryServiceItf) CategoryHandlerItf {
	return &categoryHandler{
		categoryService: categoryService,
	}
}