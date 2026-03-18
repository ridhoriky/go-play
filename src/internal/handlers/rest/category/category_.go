package category

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	_ "ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetCategories godoc
// @Summary      List Category
// @Description  Retrieve list of categories with optional filters and pagination
// @Tags         categories
// @Produce      json
// @Param		 search				query		string	false	"Filter by name"
// @Param 		 include_deleted 	query 		bool 	false 	"Include soft deleted categories"
// @Param		 page				query		int		false	"Page number"	default(1)
// @Param		 limit				query		int		false	"Page size"		default(10)
// @Param		 sort_by			query		string	false	"Sort by field"
// @Param		 sort_dir			query		string	false	"Sort direction (asc/desc)"	default(asc)
// @Success      200  {object}  	dto.APIResponse{data=[]entity.Category}
// @Failure      400  {object} 		dto.APIResponse
// @Failure      404  {object}  	dto.APIResponse
// @Router       /categories [get]
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

// CreateCategory godoc
// @Summary      Create Single category
// @Description  Create a category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param		 category	body		dto.CreateCategoryRequest		true	"Category data"
// @Success      201		{object}  	dto.APIResponse{data=entity.Category}
// @Failure      400  		{object} 	dto.APIResponse
// @Failure      404  		{object}  	dto.APIResponse
// @Router       /categories [post]
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

	createdCategory, err := h.categoryService.CreateCategory(ctx, &cat)
	if err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}

	helpers.ResponseSuccess(c.Writer, http.StatusCreated, "Category created successfully", createdCategory)
}

// GetCategoryByID godoc
// @Summary      Get Single category
// @Description  Get data category by id
// @Tags         categories
// @Produce      json
// @Param 		 id   path 		string true "Category ID (UUID)" format(uuid)
// @Success      200  {object}  dto.APIResponse{data=entity.Category}
// @Failure      400  {object} 	dto.APIResponse
// @Failure      404  {object}  dto.APIResponse
// @Router       /categories/{id} [get]
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

// UpdateCategory godoc
// @Summary      Update Category
// @Description  Update a category by their id
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param 		 id 	  path 	string true "Category ID (UUID)" format(uuid)
// @Param 		 category body 	dto.UpdateCategoryRequest true "Category data"
// @Success      200  {object}  dto.APIResponse{data=entity.Category}
// @Failure      400  {object} 	dto.APIResponse
// @Failure      404  {object}  dto.APIResponse
// @Router       /categories/{id} [put]
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

	updatedCategory, err := h.categoryService.UpdateCategory(ctx, id.String(), &cat)
	if err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}

	helpers.ResponseSuccess(c.Writer, http.StatusOK, "Category updated successfully", updatedCategory)
}

// DeleteCategory godoc
// @Summary Delete Category
// @Description  Delete a category by their id
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param 		 id   path 	string true "Category ID (UUID)" format(uuid)
// @Success      200  {object}  dto.APIResponse
// @Failure      400  {object} 	dto.APIResponse
// @Failure      404  {object}  dto.APIResponse
// @Router       /categories/{id} [delete]
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
