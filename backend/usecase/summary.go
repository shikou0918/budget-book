package usecase

import (
	"budget-book/entity"
)

// SummaryUseCase handles summary business logic
type SummaryUseCase struct {
	transactionRepo TransactionRepositoryInterface
	categoryRepo    CategoryRepositoryInterface
	budgetRepo      BudgetRepositoryInterface
}

// NewSummaryUseCase creates a new summary use case instance
func NewSummaryUseCase(transactionRepo TransactionRepositoryInterface, categoryRepo CategoryRepositoryInterface, budgetRepo BudgetRepositoryInterface) *SummaryUseCase {
	return &SummaryUseCase{
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
		budgetRepo:      budgetRepo,
	}
}

// GetMonthlySummary generates a comprehensive monthly summary with transactions and budgets
func (uc *SummaryUseCase) GetMonthlySummary(year, month int) (*entity.MonthlySummary, error) {
	summary := entity.NewMonthlySummary(year, month)

	transactions, err := uc.transactionRepo.GetByMonth(year, month)
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[uint64]*entity.Category)
	categories, err := uc.categoryRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		categoryMap[category.ID] = category
	}

	for _, transaction := range transactions {
		summary.AddTransaction(transaction)
		if category, exists := categoryMap[transaction.CategoryID]; exists {
			summary.SetCategoryInfo(transaction.CategoryID, category.Name, string(category.Type))
		}
	}

	budgets, err := uc.budgetRepo.GetByMonth(year, month)
	if err != nil {
		return nil, err
	}

	for _, budget := range budgets {
		summary.SetBudget(budget.CategoryID, budget.Amount)
	}

	return summary, nil
}

// GetCategoryTotals calculates total amounts per category for a specific month
func (uc *SummaryUseCase) GetCategoryTotals(year, month int) (map[uint64]float64, error) {
	transactions, err := uc.transactionRepo.GetByMonth(year, month)
	if err != nil {
		return nil, err
	}

	totals := make(map[uint64]float64)
	for _, transaction := range transactions {
		totals[transaction.CategoryID] += transaction.Amount
	}

	return totals, nil
}
