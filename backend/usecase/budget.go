package usecase

import (
	"budget-book/entity"
)

// BudgetRepositoryInterface defines the interface for budget repository
type BudgetRepositoryInterface interface {
	Create(budget *entity.Budget) error
	GetByID(id uint64) (*entity.Budget, error)
	GetAll() ([]*entity.Budget, error)
	GetByMonth(year, month int) ([]*entity.Budget, error)
	GetByCategoryAndMonth(categoryID uint64, year, month int) (*entity.Budget, error)
	Update(budget *entity.Budget) error
	Delete(id uint64) error
}

// BudgetUseCase handles budget business logic
type BudgetUseCase struct {
	budgetRepo   BudgetRepositoryInterface
	categoryRepo CategoryRepositoryInterface
}

// NewBudgetUseCase creates a new budget use case instance
func NewBudgetUseCase(budgetRepo BudgetRepositoryInterface, categoryRepo CategoryRepositoryInterface) *BudgetUseCase {
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

	if category.Type != entity.TransactionTypeExpense {
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

	if category.Type != entity.TransactionTypeExpense {
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
