package handler

import (
	"budget-book/internal/domain/entity"
	"budget-book/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BudgetHandler struct {
	usecase *usecase.BudgetUseCase
}

type CreateBudgetRequest struct {
	CategoryID  uint64  `json:"category_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	TargetYear  int     `json:"target_year" validate:"required,min=1900,max=2100"`
	TargetMonth int     `json:"target_month" validate:"required,min=1,max=12"`
}

type UpdateBudgetRequest struct {
	CategoryID  uint64  `json:"category_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	TargetYear  int     `json:"target_year" validate:"required,min=1900,max=2100"`
	TargetMonth int     `json:"target_month" validate:"required,min=1,max=12"`
}

func NewBudgetHandler(usecase *usecase.BudgetUseCase) *BudgetHandler {
	return &BudgetHandler{usecase: usecase}
}

func (h *BudgetHandler) CreateBudget(c echo.Context) error {
	var req CreateBudgetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	budget, err := h.usecase.CreateBudget(req.CategoryID, req.Amount, req.TargetYear, req.TargetMonth)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, budget)
}

func (h *BudgetHandler) GetBudget(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid budget ID"})
	}

	budget, err := h.usecase.GetBudgetByID(id)
	if err != nil {
		if _, ok := err.(*entity.NotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, budget)
}

func (h *BudgetHandler) GetBudgets(c echo.Context) error {
	yearParam := c.QueryParam("year")
	monthParam := c.QueryParam("month")

	if yearParam != "" && monthParam != "" {
		year, err := strconv.Atoi(yearParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid year parameter"})
		}

		month, err := strconv.Atoi(monthParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid month parameter"})
		}

		if month < 1 || month > 12 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Month must be between 1 and 12"})
		}

		budgets, err := h.usecase.GetBudgetsByMonth(year, month)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, budgets)
	}

	budgets, err := h.usecase.GetAllBudgets()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, budgets)
}

func (h *BudgetHandler) UpdateBudget(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid budget ID"})
	}

	var req UpdateBudgetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	budget, err := h.usecase.UpdateBudget(id, req.CategoryID, req.Amount, req.TargetYear, req.TargetMonth)
	if err != nil {
		if _, ok := err.(*entity.NotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, budget)
}

func (h *BudgetHandler) DeleteBudget(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid budget ID"})
	}

	if err := h.usecase.DeleteBudget(id); err != nil {
		if _, ok := err.(*entity.NotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}