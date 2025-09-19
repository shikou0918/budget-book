package entity

// MonthlySummary represents a financial summary for a specific month
type MonthlySummary struct {
	Year            int                         `json:"year"`
	Month           int                         `json:"month"`
	TotalIncome     float64                     `json:"total_income"`
	TotalExpense    float64                     `json:"total_expense"`
	Balance         float64                     `json:"balance"`
	CategorySummary map[uint64]*CategorySummary `json:"category_summary"`
}

// CategorySummary represents a financial summary for a specific category
type CategorySummary struct {
	CategoryID   uint64  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	CategoryType string  `json:"category_type"`
	Total        float64 `json:"total"`
	Budget       float64 `json:"budget"`
	Percentage   float64 `json:"percentage"`
}

// NewMonthlySummary creates a new MonthlySummary instance for the given year and month
func NewMonthlySummary(year, month int) *MonthlySummary {
	return &MonthlySummary{
		Year:            year,
		Month:           month,
		CategorySummary: make(map[uint64]*CategorySummary),
	}
}

// AddTransaction adds a transaction to the monthly summary and updates totals
func (ms *MonthlySummary) AddTransaction(transaction *Transaction) {
	if transaction.Type == TransactionTypeIncome {
		ms.TotalIncome += transaction.Amount
	} else {
		ms.TotalExpense += transaction.Amount
	}
	ms.Balance = ms.TotalIncome - ms.TotalExpense

	if ms.CategorySummary[transaction.CategoryID] == nil {
		ms.CategorySummary[transaction.CategoryID] = &CategorySummary{
			CategoryID: transaction.CategoryID,
		}
	}
	ms.CategorySummary[transaction.CategoryID].Total += transaction.Amount
}

// SetCategoryInfo sets the category name and type for a given category ID
func (ms *MonthlySummary) SetCategoryInfo(categoryID uint64, name, categoryType string) {
	if ms.CategorySummary[categoryID] == nil {
		ms.CategorySummary[categoryID] = &CategorySummary{
			CategoryID: categoryID,
		}
	}
	ms.CategorySummary[categoryID].CategoryName = name
	ms.CategorySummary[categoryID].CategoryType = categoryType
}

// SetBudget sets the budget for a category and calculates the percentage used
func (ms *MonthlySummary) SetBudget(categoryID uint64, budget float64) {
	if ms.CategorySummary[categoryID] == nil {
		ms.CategorySummary[categoryID] = &CategorySummary{
			CategoryID: categoryID,
		}
	}
	ms.CategorySummary[categoryID].Budget = budget
	if budget > 0 {
		ms.CategorySummary[categoryID].Percentage = (ms.CategorySummary[categoryID].Total / budget) * 100
	}
}
