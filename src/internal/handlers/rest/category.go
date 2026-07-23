package rest

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	_ "ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/category"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	categoryService category.CategoryServiceItf
}

func NewCategoryHandler(categoryService category.CategoryServiceItf) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetCategories godoc
// @Summary      List Category
// @Description  Retrieve list of categories as a tree structure
// @Tags         categories
// @Produce      json
// @Success      200  {object}  	dto.APIResponse{data=[]dto.CategoryTreeNode}
// @Failure      400  {object} 		dto.APIResponse
// @Failure      404  {object}  	dto.APIResponse
// @Router       /categories [get]
func (h *CategoryHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	categories, err := h.categoryService.GetCategoryTree(ctx)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", categories)
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
// @Router       /admin/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var cat dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&cat); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	createdCategory, err := h.categoryService.CreateCategory(ctx, &cat)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Category created successfully", createdCategory)
}

// GetCategoryByID godoc
// @Summary      Get Single category
// @Description  Get data category by id
// @Tags         categories
// @Produce      json
// @Param 		 id   path 		string true "Category ID (UUID)" format(uuid)
// @Success      200  {object}  dto.APIResponse{data=dto.CategoryDetailResponse}
// @Failure      400  {object} 	dto.APIResponse
// @Failure      404  {object}  	dto.APIResponse
// @Router       /categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidCategoryID))
		return
	}
	category, err := h.categoryService.GetCategoryByID(ctx, id.String())
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", category)
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
// @Router       /admin/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidCategoryID))
		return
	}

	var cat dto.UpdateCategoryRequest
	if err = c.ShouldBindJSON(&cat); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	updatedCategory, err := h.categoryService.UpdateCategory(ctx, id.String(), &cat)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Category updated successfully", updatedCategory)
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
// @Router       /admin/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidCategoryID))
		return
	}
	if err := h.categoryService.DeleteCategory(ctx, id.String()); err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Category deleted successfully", nil)
}

// GetCategoryTree godoc
// @Summary      Get Category Tree
// @Description  Get full category tree hierarchy
// @Tags         categories
// @Produce      json
// @Success      200  {object}  dto.APIResponse{data=[]dto.CategoryTreeNode}
// @Failure      400  {object} 	dto.APIResponse
// @Failure      500  {object}  dto.APIResponse
// @Router       /categories/tree [get]
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	ctx := c.Request.Context()
	tree, err := h.categoryService.GetCategoryTree(ctx)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", tree)
}
