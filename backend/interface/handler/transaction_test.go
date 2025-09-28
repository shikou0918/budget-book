package handler

import (
	"budget-book/entity"
	mock_usecase "budget-book/mocks/usecase"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	return e
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTransactionUseCaseInterface(ctrl)
	handler := NewTransactionHandler(mockUseCase)

	e := setupEcho()

	t.Run("正常な取引作成", func(t *testing.T) {
		req := CreateTransactionRequest{
			Type:            "income",
			Amount:          50000.0,
			CategoryID:      1,
			TransactionDate: "2024-01-15",
			Memo:            "給与",
		}

		expectedTransaction := &entity.Transaction{
			ID:              1,
			Type:            entity.TransactionTypeIncome,
			Amount:          50000.0,
			CategoryID:      1,
			TransactionDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Memo:            "給与",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		mockUseCase.EXPECT().
			CreateTransaction(
				entity.TransactionTypeIncome,
				50000.0,
				uint64(1),
				time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				"給与",
			).
			Return(expectedTransaction, nil)

		reqBody, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}
		httpReq := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(reqBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		err = handler.CreateTransaction(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var response entity.Transaction
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedTransaction.ID, response.ID)
		assert.Equal(t, expectedTransaction.Amount, response.Amount)
	})

	t.Run("不正なリクエストボディ", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodPost, "/transactions", strings.NewReader("invalid json"))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		err := handler.CreateTransaction(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		req := CreateTransactionRequest{
			Type:            "invalid",
			Amount:          -100.0,
			CategoryID:      0,
			TransactionDate: "2024-01-15",
			Memo:            "テスト",
		}

		reqBody, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}
		httpReq := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(reqBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		err = handler.CreateTransaction(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("不正な日付フォーマット", func(t *testing.T) {
		req := CreateTransactionRequest{
			Type:            "income",
			Amount:          50000.0,
			CategoryID:      1,
			TransactionDate: "invalid-date",
			Memo:            "給与",
		}

		reqBody, err := json.Marshal(req)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}
		httpReq := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(reqBody))
		httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		err = handler.CreateTransaction(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestTransactionHandler_GetTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTransactionUseCaseInterface(ctrl)
	handler := NewTransactionHandler(mockUseCase)

	e := setupEcho()

	t.Run("正常な取引取得", func(t *testing.T) {
		transactionID := uint64(1)
		expectedTransaction := &entity.Transaction{
			ID:     transactionID,
			Type:   entity.TransactionTypeIncome,
			Amount: 50000.0,
		}

		mockUseCase.EXPECT().
			GetTransactionByID(transactionID).
			Return(expectedTransaction, nil)

		httpReq := httptest.NewRequest(http.MethodGet, "/transactions/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetPath("/transactions/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.GetTransaction(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response entity.Transaction
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedTransaction.ID, response.ID)
	})

	t.Run("無効なID", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodGet, "/transactions/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetPath("/transactions/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.GetTransaction(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("取引が見つからない", func(t *testing.T) {
		transactionID := uint64(999)

		mockUseCase.EXPECT().
			GetTransactionByID(transactionID).
			Return(nil, entity.NewNotFoundError("transaction", transactionID))

		httpReq := httptest.NewRequest(http.MethodGet, "/transactions/999", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetPath("/transactions/:id")
		c.SetParamNames("id")
		c.SetParamValues("999")

		err := handler.GetTransaction(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestTransactionHandler_GetTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTransactionUseCaseInterface(ctrl)
	handler := NewTransactionHandler(mockUseCase)

	e := setupEcho()

	t.Run("正常な取引一覧取得", func(t *testing.T) {
		expectedTransactions := []*entity.Transaction{
			{
				ID:     1,
				Type:   entity.TransactionTypeIncome,
				Amount: 50000.0,
			},
			{
				ID:     2,
				Type:   entity.TransactionTypeExpense,
				Amount: 1200.0,
			},
		}

		mockUseCase.EXPECT().
			GetAllTransactions().
			Return(expectedTransactions, nil)

		httpReq := httptest.NewRequest(http.MethodGet, "/transactions", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)

		err := handler.GetTransactions(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []*entity.Transaction
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, expectedTransactions[0].ID, response[0].ID)
	})
}

func TestTransactionHandler_DeleteTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mock_usecase.NewMockTransactionUseCaseInterface(ctrl)
	handler := NewTransactionHandler(mockUseCase)

	e := setupEcho()

	t.Run("正常な取引削除", func(t *testing.T) {
		transactionID := uint64(1)

		mockUseCase.EXPECT().
			DeleteTransaction(transactionID).
			Return(nil)

		httpReq := httptest.NewRequest(http.MethodDelete, "/transactions/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetPath("/transactions/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.DeleteTransaction(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("存在しない取引の削除", func(t *testing.T) {
		transactionID := uint64(999)

		mockUseCase.EXPECT().
			DeleteTransaction(transactionID).
			Return(entity.NewNotFoundError("transaction", transactionID))

		httpReq := httptest.NewRequest(http.MethodDelete, "/transactions/999", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(httpReq, rec)
		c.SetPath("/transactions/:id")
		c.SetParamNames("id")
		c.SetParamValues("999")

		err := handler.DeleteTransaction(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
