package category

import (
	"encoding/json"
	"log"
	"net/http"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *categoryHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	categories, err := h.categoryService.GetAllCategories(ctx)
	if err != nil {
		dto.ResponseError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}
	dto.ResponseSuccess(c.Writer, http.StatusOK, "Success", categories)
}

func (h *categoryHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var cat dto.CategoryDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&cat); err != nil {
		log.Println(err)
		dto.ResponseError(c.Writer, http.StatusBadRequest, preference.ErrInvalidReqBody)
		return
	}
	if err := cat.Validate(); err != nil {
		log.Println(err)
		dto.ResponseError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.categoryService.CreateCategory(ctx, &cat); err != nil {
		log.Println(err)
		dto.ResponseError(c.Writer, http.StatusInternalServerError, "Failed to create category")
		return
	}
	dto.ResponseSuccess(c.Writer, http.StatusCreated, "Category created successfully", cat)
}

func (h *categoryHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, preference.ErrInvalidCategoryID)
		return
	}
	category, err := h.categoryService.GetCategoryByID(ctx, id.String())
	if err != nil {
		dto.ResponseError(c.Writer, http.StatusInternalServerError, "Failed to fetch category")
		return
	}
	dto.ResponseSuccess(c.Writer, http.StatusOK, "Success", category)
}

func (h *categoryHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, preference.ErrInvalidCategoryID)
		return
	}
	var cat dto.CategoryDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&cat); err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, preference.ErrInvalidReqBody)
		return
	}
	if err := cat.Validate(); err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}
	cat.ID = id.String()
	if err := h.categoryService.UpdateCategory(ctx, &cat); err != nil {
		dto.ResponseError(c.Writer, http.StatusInternalServerError, "Failed to update category")
		return
	}
	dto.ResponseSuccess(c.Writer, http.StatusOK, "Category updated successfully", cat)
}

func (h *categoryHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, preference.ErrInvalidCategoryID)
		return
	}
	if err := h.categoryService.DeleteCategory(ctx, id.String()); err != nil {
		dto.ResponseError(c.Writer, http.StatusInternalServerError, "Failed to delete category")
		return
	}
	dto.ResponseSuccess(c.Writer, http.StatusOK, "Category deleted successfully", nil)
}
