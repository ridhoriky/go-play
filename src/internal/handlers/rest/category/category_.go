package category

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *categoryHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.GetCategoriesQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidReqBody,
		})
		return
	}

	categories, err := h.categoryService.GetAllCategories(ctx, &req)
	if err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}
	helpers.ResponseSuccess(c.Writer, http.StatusOK, "Success", categories)
}

func (h *categoryHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var cat dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&cat); err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidReqBody,
		})
		return
	}

	if _, err := h.categoryService.CreateCategory(ctx, &cat); err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}

	helpers.ResponseSuccess(c.Writer, http.StatusCreated, "Category created successfully", cat)
}

func (h *categoryHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidCategoryID,
		})
		return
	}
	category, err := h.categoryService.GetCategoryByID(ctx, id.String())
	if err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}
	helpers.ResponseSuccess(c.Writer, http.StatusOK, "Success", category)
}

func (h *categoryHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidCategoryID,
		})
		return
	}

	var cat dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&cat); err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidReqBody,
		})
		return
	}

	if err := h.categoryService.UpdateCategory(ctx, id.String(), &cat); err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}

	helpers.ResponseSuccess(c.Writer, http.StatusOK, "Category updated successfully", cat)
}

func (h *categoryHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidCategoryID,
		})
		return
	}
	if err := h.categoryService.DeleteCategory(ctx, id.String()); err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}
	helpers.ResponseSuccess(c.Writer, http.StatusOK, "Category deleted successfully", nil)
}
