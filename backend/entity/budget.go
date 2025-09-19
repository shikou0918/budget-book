package entity

import (
	"time"
)

// Budget represents a budget for a specific category and month
type Budget struct {
	ID          uint64    `json:"id"`
	CategoryID  uint64    `json:"category_id"`
	Category    *Category `json:"category,omitempty"`
	Amount      float64   `json:"amount"`
	TargetYear  int       `json:"target_year"`
	TargetMonth int       `json:"target_month"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewBudget creates a new budget instance
func NewBudget(categoryID uint64, amount float64, targetYear, targetMonth int) *Budget {
	return &Budget{
		CategoryID:  categoryID,
		Amount:      amount,
		TargetYear:  targetYear,
		TargetMonth: targetMonth,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// IsValid validates the budget data
func (b *Budget) IsValid() error {
	if b.CategoryID == 0 {
		return NewValidationError("category_id is required")
	}
	if b.Amount <= 0 {
		return NewValidationError("amount must be greater than 0")
	}
	if b.TargetYear < 1900 || b.TargetYear > 2100 {
		return NewValidationError("target_year must be between 1900 and 2100")
	}
	if b.TargetMonth < 1 || b.TargetMonth > 12 {
		return NewValidationError("target_month must be between 1 and 12")
	}
	return nil
}
