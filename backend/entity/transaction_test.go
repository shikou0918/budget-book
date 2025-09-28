package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	// テストデータの準備
	transactionType := TransactionTypeIncome
	amount := 10000.0
	categoryID := uint64(1)
	transactionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	memo := "給与"

	// NewTransaction関数のテスト
	transaction := NewTransaction(transactionType, amount, categoryID, transactionDate, memo)

	// 結果の検証
	assert.Equal(t, transactionType, transaction.Type)
	assert.Equal(t, amount, transaction.Amount)
	assert.Equal(t, categoryID, transaction.CategoryID)
	assert.Equal(t, transactionDate, transaction.TransactionDate)
	assert.Equal(t, memo, transaction.Memo)
	assert.NotZero(t, transaction.CreatedAt)
	assert.NotZero(t, transaction.UpdatedAt)
}

func TestTransaction_IsValid(t *testing.T) {
	validTransaction := &Transaction{
		Type:            TransactionTypeIncome,
		Amount:          1000.0,
		CategoryID:      1,
		TransactionDate: time.Now(),
	}

	tests := []struct {
		name        string
		transaction *Transaction
		wantErr     bool
		errMessage  string
	}{
		{
			name:        "有効な取引",
			transaction: validTransaction,
			wantErr:     false,
		},
		{
			name: "金額が0以下",
			transaction: &Transaction{
				Type:            TransactionTypeIncome,
				Amount:          0,
				CategoryID:      1,
				TransactionDate: time.Now(),
			},
			wantErr:    true,
			errMessage: "amount must be greater than 0",
		},
		{
			name: "金額が負数",
			transaction: &Transaction{
				Type:            TransactionTypeIncome,
				Amount:          -100,
				CategoryID:      1,
				TransactionDate: time.Now(),
			},
			wantErr:    true,
			errMessage: "amount must be greater than 0",
		},
		{
			name: "カテゴリIDが未設定",
			transaction: &Transaction{
				Type:            TransactionTypeIncome,
				Amount:          1000,
				CategoryID:      0,
				TransactionDate: time.Now(),
			},
			wantErr:    true,
			errMessage: "category_id is required",
		},
		{
			name: "取引日が未設定",
			transaction: &Transaction{
				Type:       TransactionTypeIncome,
				Amount:     1000,
				CategoryID: 1,
			},
			wantErr:    true,
			errMessage: "transaction_date is required",
		},
		{
			name: "無効な取引タイプ",
			transaction: &Transaction{
				Type:            "invalid",
				Amount:          1000,
				CategoryID:      1,
				TransactionDate: time.Now(),
			},
			wantErr:    true,
			errMessage: "type must be 'income' or 'expense'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.transaction.IsValid()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTransactionType_Constants(t *testing.T) {
	assert.Equal(t, TransactionType("income"), TransactionTypeIncome)
	assert.Equal(t, TransactionType("expense"), TransactionTypeExpense)
}
