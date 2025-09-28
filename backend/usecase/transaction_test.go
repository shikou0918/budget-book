package usecase

import (
	"budget-book/entity"
	mock_repository "budget-book/mocks/repository"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTransactionUseCase_CreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock_repository.NewMockTransactionRepositoryInterface(ctrl)
	mockCategoryRepo := mock_repository.NewMockCategoryRepositoryInterface(ctrl)

	usecase := NewTransactionUseCase(mockTransactionRepo, mockCategoryRepo)

	// テストデータ
	categoryID := uint64(1)
	category := &entity.Category{
		ID:   categoryID,
		Name: "給与",
		Type: entity.TransactionTypeIncome,
	}
	transactionType := entity.TransactionTypeIncome
	amount := 50000.0
	transactionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	memo := "給与"

	t.Run("正常な取引作成", func(t *testing.T) {
		// モックの設定
		mockCategoryRepo.EXPECT().
			GetByID(categoryID).
			Return(category, nil)

		mockTransactionRepo.EXPECT().
			Create(gomock.Any()).
			Return(nil)

		// テスト実行
		result, err := usecase.CreateTransaction(transactionType, amount, categoryID, transactionDate, memo)

		// 結果検証
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, transactionType, result.Type)
		assert.Equal(t, amount, result.Amount)
		assert.Equal(t, categoryID, result.CategoryID)
		assert.Equal(t, transactionDate, result.TransactionDate)
		assert.Equal(t, memo, result.Memo)
	})

	t.Run("カテゴリが見つからない場合", func(t *testing.T) {
		mockCategoryRepo.EXPECT().
			GetByID(categoryID).
			Return(nil, entity.NewNotFoundError("category", categoryID))

		result, err := usecase.CreateTransaction(transactionType, amount, categoryID, transactionDate, memo)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("取引タイプとカテゴリタイプが一致しない場合", func(t *testing.T) {
		// 支出カテゴリを設定
		expenseCategory := &entity.Category{
			ID:   categoryID,
			Name: "食費",
			Type: entity.TransactionTypeExpense,
		}

		mockCategoryRepo.EXPECT().
			GetByID(categoryID).
			Return(expenseCategory, nil)

		// 収入タイプで取引を作成しようとする
		result, err := usecase.CreateTransaction(transactionType, amount, categoryID, transactionDate, memo)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "transaction type does not match category type")
	})

	t.Run("リポジトリでエラーが発生", func(t *testing.T) {
		mockCategoryRepo.EXPECT().
			GetByID(categoryID).
			Return(category, nil)

		mockTransactionRepo.EXPECT().
			Create(gomock.Any()).
			Return(errors.New("database error"))

		result, err := usecase.CreateTransaction(transactionType, amount, categoryID, transactionDate, memo)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestTransactionUseCase_GetTransactionByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock_repository.NewMockTransactionRepositoryInterface(ctrl)
	mockCategoryRepo := mock_repository.NewMockCategoryRepositoryInterface(ctrl)

	usecase := NewTransactionUseCase(mockTransactionRepo, mockCategoryRepo)

	transactionID := uint64(1)
	expectedTransaction := &entity.Transaction{
		ID:     transactionID,
		Type:   entity.TransactionTypeIncome,
		Amount: 50000.0,
	}

	t.Run("正常な取引取得", func(t *testing.T) {
		mockTransactionRepo.EXPECT().
			GetByID(transactionID).
			Return(expectedTransaction, nil)

		result, err := usecase.GetTransactionByID(transactionID)

		assert.NoError(t, err)
		assert.Equal(t, expectedTransaction, result)
	})

	t.Run("取引が見つからない場合", func(t *testing.T) {
		mockTransactionRepo.EXPECT().
			GetByID(transactionID).
			Return(nil, entity.NewNotFoundError("transaction", transactionID))

		result, err := usecase.GetTransactionByID(transactionID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTransactionUseCase_UpdateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock_repository.NewMockTransactionRepositoryInterface(ctrl)
	mockCategoryRepo := mock_repository.NewMockCategoryRepositoryInterface(ctrl)

	usecase := NewTransactionUseCase(mockTransactionRepo, mockCategoryRepo)

	transactionID := uint64(1)
	categoryID := uint64(1)
	existingTransaction := &entity.Transaction{
		ID:     transactionID,
		Type:   entity.TransactionTypeIncome,
		Amount: 50000.0,
	}
	category := &entity.Category{
		ID:   categoryID,
		Name: "給与",
		Type: entity.TransactionTypeIncome,
	}
	transactionType := entity.TransactionTypeIncome
	amount := 60000.0
	transactionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	memo := "給与（更新）"

	t.Run("正常な取引更新", func(t *testing.T) {
		mockTransactionRepo.EXPECT().
			GetByID(transactionID).
			Return(existingTransaction, nil)

		mockCategoryRepo.EXPECT().
			GetByID(categoryID).
			Return(category, nil)

		mockTransactionRepo.EXPECT().
			Update(gomock.Any()).
			Return(nil)

		result, err := usecase.UpdateTransaction(transactionID, transactionType, amount, categoryID, transactionDate, memo)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, amount, result.Amount)
		assert.Equal(t, memo, result.Memo)
	})
}

func TestTransactionUseCase_DeleteTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactionRepo := mock_repository.NewMockTransactionRepositoryInterface(ctrl)
	mockCategoryRepo := mock_repository.NewMockCategoryRepositoryInterface(ctrl)

	usecase := NewTransactionUseCase(mockTransactionRepo, mockCategoryRepo)

	transactionID := uint64(1)
	existingTransaction := &entity.Transaction{
		ID:     transactionID,
		Type:   entity.TransactionTypeIncome,
		Amount: 50000.0,
	}

	t.Run("正常な取引削除", func(t *testing.T) {
		mockTransactionRepo.EXPECT().
			GetByID(transactionID).
			Return(existingTransaction, nil)

		mockTransactionRepo.EXPECT().
			Delete(transactionID).
			Return(nil)

		err := usecase.DeleteTransaction(transactionID)

		assert.NoError(t, err)
	})

	t.Run("存在しない取引の削除", func(t *testing.T) {
		mockTransactionRepo.EXPECT().
			GetByID(transactionID).
			Return(nil, entity.NewNotFoundError("transaction", transactionID))

		err := usecase.DeleteTransaction(transactionID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
