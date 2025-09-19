package entity

import (
	"time"
)

// CategoryType represents the type of category (income or expense)
type CategoryType string

const (
	// CategoryTypeIncome represents income categories
	CategoryTypeIncome CategoryType = "income"
	// CategoryTypeExpense represents expense categories
	CategoryTypeExpense CategoryType = "expense"
)

// Category represents a transaction category
type Category struct {
	ID        uint64       `json:"id"`
	Name      string       `json:"name"`
	Type      CategoryType `json:"type"`
	Color     string       `json:"color"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// NewCategory creates a new category instance
func NewCategory(name string, categoryType CategoryType, color string) *Category {
	if color == "" {
		color = "#007BFF"
	}
	return &Category{
		Name:      name,
		Type:      categoryType,
		Color:     color,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// IsValid validates the category data
func (c *Category) IsValid() error {
	if c.Name == "" {
		return NewValidationError("name is required")
	}
	if len(c.Name) > 50 {
		return NewValidationError("name must be 50 characters or less")
	}
	if c.Type != CategoryTypeIncome && c.Type != CategoryTypeExpense {
		return NewValidationError("type must be 'income' or 'expense'")
	}
	if c.Color != "" && len(c.Color) != 7 {
		return NewValidationError("color must be a valid hex color code")
	}
	return nil
}
