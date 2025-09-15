package handler

import (
	"budget-book/internal/domain/entity"
	"budget-book/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	usecase *usecase.TransactionUseCase
}

type CreateTransactionRequest struct {
	Type            string  `json:"type" validate:"required,oneof=income expense"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	CategoryID      uint64  `json:"category_id" validate:"required"`
	TransactionDate string  `json:"transaction_date" validate:"required"`
	Memo            string  `json:"memo"`
}

type UpdateTransactionRequest struct {
	Type            string  `json:"type" validate:"required,oneof=income expense"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	CategoryID      uint64  `json:"category_id" validate:"required"`
	TransactionDate string  `json:"transaction_date" validate:"required"`
	Memo            string  `json:"memo"`
}

func NewTransactionHandler(usecase *usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{usecase: usecase}
}

func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	var req CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid transaction_date format. Use YYYY-MM-DD"})
	}

	transactionType := entity.TransactionType(req.Type)
	transaction, err := h.usecase.CreateTransaction(transactionType, req.Amount, req.CategoryID, transactionDate, req.Memo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) GetTransaction(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid transaction ID"})
	}

	transaction, err := h.usecase.GetTransactionByID(id)
	if err != nil {
		if _, ok := err.(*entity.NotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) GetTransactions(c echo.Context) error {
	transactions, err := h.usecase.GetAllTransactions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) UpdateTransaction(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid transaction ID"})
	}

	var req UpdateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid transaction_date format. Use YYYY-MM-DD"})
	}

	transactionType := entity.TransactionType(req.Type)
	transaction, err := h.usecase.UpdateTransaction(id, transactionType, req.Amount, req.CategoryID, transactionDate, req.Memo)
	if err != nil {
		if _, ok := err.(*entity.NotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) DeleteTransaction(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid transaction ID"})
	}

	if err := h.usecase.DeleteTransaction(id); err != nil {
		if _, ok := err.(*entity.NotFoundError); ok {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
