package handler

import (
	"budget-book/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// SummaryHandler handles summary HTTP requests
type SummaryHandler struct {
	usecase *usecase.SummaryUseCase
}

// NewSummaryHandler creates a new summary handler instance
func NewSummaryHandler(usecase *usecase.SummaryUseCase) *SummaryHandler {
	return &SummaryHandler{usecase: usecase}
}

// GetMonthlySummary handles GET /summary/:year/:month endpoint
func (h *SummaryHandler) GetMonthlySummary(c echo.Context) error {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid year parameter"})
	}

	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid month parameter"})
	}

	if month < 1 || month > 12 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Month must be between 1 and 12"})
	}

	summary, err := h.usecase.GetMonthlySummary(year, month)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, summary)
}
