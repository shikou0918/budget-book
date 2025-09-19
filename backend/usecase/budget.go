package usecase

import (
	"budget-book/entity"
	"budget-book/infrastructure/repository"
)

// BudgetUseCase handles budget business logic
type BudgetUseCase struct {
	budgetRepo   *repository.BudgetRepository
	categoryRepo *repository.CategoryRepository
}

// NewBudgetUseCase creates a new budget use case instance
func NewBudgetUseCase(budgetRepo *repository.BudgetRepository, categoryRepo *repository.CategoryRepository) *BudgetUseCase {
	return &BudgetUseCase{
		budgetRepo:   budgetRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateBudget creates a new budget with validation
func (uc *BudgetUseCase) CreateBudget(categoryID uint64, amount float64, targetYear, targetMonth int) (*entity.Budget, error) {
	category, err := uc.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, err
	}

	if category.Type != entity.CategoryTypeExpense {
		return nil, entity.NewValidationError("budget can only be set for expense categories")
	}

	budget := entity.NewBudget(categoryID, amount, targetYear, targetMonth)
	if err := uc.budgetRepo.Create(budget); err != nil {
		return nil, err
	}

	return budget, nil
}

// GetBudgetByID retrieves a budget by its ID
func (uc *BudgetUseCase) GetBudgetByID(id uint64) (*entity.Budget, error) {
	return uc.budgetRepo.GetByID(id)
}

// GetAllBudgets retrieves all budgets
func (uc *BudgetUseCase) GetAllBudgets() ([]*entity.Budget, error) {
	return uc.budgetRepo.GetAll()
}

// GetBudgetsByMonth retrieves budgets for a specific month
func (uc *BudgetUseCase) GetBudgetsByMonth(year, month int) ([]*entity.Budget, error) {
	return uc.budgetRepo.GetByMonth(year, month)
}

// GetBudgetByCategoryAndMonth retrieves a budget for a specific category and month
func (uc *BudgetUseCase) GetBudgetByCategoryAndMonth(categoryID uint64, year, month int) (*entity.Budget, error) {
	return uc.budgetRepo.GetByCategoryAndMonth(categoryID, year, month)
}

// UpdateBudget updates an existing budget with validation
func (uc *BudgetUseCase) UpdateBudget(id uint64, categoryID uint64, amount float64, targetYear, targetMonth int) (*entity.Budget, error) {
	budget, err := uc.budgetRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	category, err := uc.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, err
	}

	if category.Type != entity.CategoryTypeExpense {
		return nil, entity.NewValidationError("budget can only be set for expense categories")
	}

	budget.CategoryID = categoryID
	budget.Amount = amount
	budget.TargetYear = targetYear
	budget.TargetMonth = targetMonth

	if err := uc.budgetRepo.Update(budget); err != nil {
		return nil, err
	}

	return budget, nil
}

// DeleteBudget deletes a budget by its ID
func (uc *BudgetUseCase) DeleteBudget(id uint64) error {
	_, err := uc.budgetRepo.GetByID(id)
	if err != nil {
		return err
	}

	return uc.budgetRepo.Delete(id)
}
