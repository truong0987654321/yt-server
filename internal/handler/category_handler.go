package handler

import (
	"net/http"
	"yt-clone-server/internal/apperrors"
	"yt-clone-server/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	categoryService domain.CategoryService
}

func NewCategoryHandler(categoryService domain.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

type createCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type updateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// POST /api/categories
func (h *CategoryHandler) Create(c *gin.Context) {
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := apperrors.NewBadRequest("Name is required")
		c.JSON(appErr.Status(), appErr)
		return
	}
	category, err := h.categoryService.CreateCategory(c.Request.Context(), req.Name, req.Description)
	if err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GET /api/categories
func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories(c.Request.Context())
	if err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GET /api/categories/:id
func (h *CategoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		appErr := apperrors.NewBadRequest("Invalid category ID format")
		c.JSON(appErr.Status(), appErr)
		return
	}
	category, err := h.categoryService.GetCategoryByID(c.Request.Context(), id)
	if err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusOK, category)
}

// PUT /api/categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		appErr := apperrors.NewBadRequest("Invalid category ID format")
		c.JSON(appErr.Status(), appErr)
		return
	}

	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := apperrors.NewBadRequest("Invalid request payload")
		c.JSON(appErr.Status(), appErr)
		return
	}

	category, err := h.categoryService.UpdateCategory(c.Request.Context(), id, req.Name, req.Description)
	if err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusOK, category)
}

// DELETE /api/categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		appErr := apperrors.NewBadRequest("Invalid category ID format")
		c.JSON(appErr.Status(), appErr)
		return
	}

	if err := h.categoryService.DeleteCategory(c.Request.Context(), id); err != nil {
		appErr := apperrors.Parse(err)
		c.JSON(appErr.Status(), appErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
