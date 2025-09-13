package entity

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type Transaction struct {
	ID              uint64          `json:"id"`
	Type            TransactionType `json:"type"`
	Amount          float64         `json:"amount"`
	CategoryID      uint64          `json:"category_id"`
	Category        *Category       `json:"category,omitempty"`
	TransactionDate time.Time       `json:"transaction_date"`
	Memo            string          `json:"memo"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func NewTransaction(transactionType TransactionType, amount float64, categoryID uint64, transactionDate time.Time, memo string) *Transaction {
	return &Transaction{
		Type:            transactionType,
		Amount:          amount,
		CategoryID:      categoryID,
		TransactionDate: transactionDate,
		Memo:            memo,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func (t *Transaction) IsValid() error {
	if t.Amount <= 0 {
		return NewValidationError("amount must be greater than 0")
	}
	if t.CategoryID == 0 {
		return NewValidationError("category_id is required")
	}
	if t.TransactionDate.IsZero() {
		return NewValidationError("transaction_date is required")
	}
	if t.Type != TransactionTypeIncome && t.Type != TransactionTypeExpense {
		return NewValidationError("type must be 'income' or 'expense'")
	}
	return nil
}