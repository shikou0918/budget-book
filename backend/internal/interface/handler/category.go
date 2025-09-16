package handler

import (
	"budget-book/internal/domain/entity"
	"budget-book/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// CategoryHandler handles category HTTP requests
type CategoryHandler struct {
	usecase *usecase.CategoryUseCase
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name  string `json:"name" validate:"required,max=50"`
	Type  string `json:"type" validate:"required,oneof=income expense"`
	Color string `json:"color"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name  string `json:"name" validate:"required,max=50"`
	Type  string `json:"type" validate:"required,oneof=income expense"`
	Color string `json:"color"`
}

// NewCategoryHandler creates a new category handler instance
func NewCategoryHandler(usecase *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{usecase: usecase}
}

// CreateCategory handles POST /categories endpoint
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var req CreateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	categoryType := entity.CategoryType(req.Type)
	category, err := h.usecase.CreateCategory(req.Name, categoryType, req.Color)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, category)
}

// GetCategory handles GET /categories/:id endpoint
func (h *CategoryHandler) GetCategory(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	category, err := h.usecase.GetCategoryByID(id)
	if err != nil {
		if _, ok := err.(*entity.NotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, category)
}

// GetCategories handles GET /categories endpoint
func (h *CategoryHandler) GetCategories(c echo.Context) error {
	categoryType := c.QueryParam("type")

	if categoryType != "" {
		if categoryType != "income" && categoryType != "expense" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category type. Use 'income' or 'expense'"})
		}

		categories, err := h.usecase.GetCategoriesByType(entity.CategoryType(categoryType))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, categories)
	}

	categories, err := h.usecase.GetAllCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, categories)
}

// UpdateCategory handles PUT /categories/:id endpoint
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	var req UpdateCategoryRequest
	if bindErr := c.Bind(&req); bindErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if validErr := c.Validate(&req); validErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": validErr.Error()})
	}

	categoryType := entity.CategoryType(req.Type)
	category, err := h.usecase.UpdateCategory(id, req.Name, categoryType, req.Color)
	if err != nil {
		if _, ok := err.(*entity.NotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, category)
}

// DeleteCategory handles DELETE /categories/:id endpoint
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID"})
	}

	if err := h.usecase.DeleteCategory(id); err != nil {
		if _, ok := err.(*entity.NotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
