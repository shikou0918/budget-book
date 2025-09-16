package repository

import (
	"budget-book/internal/domain/entity"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// BudgetRepository handles budget data operations
type BudgetRepository struct {
	db *gorm.DB
}

// NewBudgetRepository creates a new budget repository instance
func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

// Create saves a new budget to the database
func (r *BudgetRepository) Create(budget *entity.Budget) error {
	if err := budget.IsValid(); err != nil {
		return err
	}

	exists, err := r.ExistsByCategoryAndMonth(budget.CategoryID, budget.TargetYear, budget.TargetMonth)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("budget for category %d in %d-%02d already exists", budget.CategoryID, budget.TargetYear, budget.TargetMonth)
	}

	result := r.db.Create(budget)
	if result.Error != nil {
		return fmt.Errorf("failed to create budget: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a budget by its ID
func (r *BudgetRepository) GetByID(id uint64) (*entity.Budget, error) {
	var budget entity.Budget
	result := r.db.Preload("Category").First(&budget, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, entity.NewNotFoundError("budget", id)
		}
		return nil, fmt.Errorf("failed to get budget: %w", result.Error)
	}

	return &budget, nil
}

// GetAll retrieves all budgets ordered by target year and month
func (r *BudgetRepository) GetAll() ([]*entity.Budget, error) {
	var budgets []*entity.Budget
	result := r.db.Preload("Category").Order("target_year DESC, target_month DESC").Find(&budgets)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get budgets: %w", result.Error)
	}

	return budgets, nil
}

// GetByMonth retrieves all budgets for a specific year and month
func (r *BudgetRepository) GetByMonth(year, month int) ([]*entity.Budget, error) {
	var budgets []*entity.Budget
	result := r.db.Preload("Category").
		Where("target_year = ? AND target_month = ?", year, month).
		Find(&budgets)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get budgets by month: %w", result.Error)
	}

	return budgets, nil
}

// GetByCategoryAndMonth retrieves a budget by category ID and target month
func (r *BudgetRepository) GetByCategoryAndMonth(categoryID uint64, year, month int) (*entity.Budget, error) {
	var budget entity.Budget
	result := r.db.Preload("Category").
		Where("category_id = ? AND target_year = ? AND target_month = ?", categoryID, year, month).
		First(&budget)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, entity.NewNotFoundError("budget", fmt.Sprintf("category:%d year:%d month:%d", categoryID, year, month))
		}
		return nil, fmt.Errorf("failed to get budget: %w", result.Error)
	}

	return &budget, nil
}

// Update modifies an existing budget in the database
func (r *BudgetRepository) Update(budget *entity.Budget) error {
	if err := budget.IsValid(); err != nil {
		return err
	}

	budget.UpdatedAt = time.Now()
	result := r.db.Save(budget)
	if result.Error != nil {
		return fmt.Errorf("failed to update budget: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return entity.NewNotFoundError("budget", budget.ID)
	}

	return nil
}

// Delete removes a budget from the database by ID
func (r *BudgetRepository) Delete(id uint64) error {
	result := r.db.Delete(&entity.Budget{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete budget: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return entity.NewNotFoundError("budget", id)
	}

	return nil
}

// ExistsByCategoryAndMonth checks if a budget exists for a category in a specific month
func (r *BudgetRepository) ExistsByCategoryAndMonth(categoryID uint64, year, month int) (bool, error) {
	var count int64
	result := r.db.Model(&entity.Budget{}).
		Where("category_id = ? AND target_year = ? AND target_month = ?", categoryID, year, month).
		Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check budget existence: %w", result.Error)
	}

	return count > 0, nil
}
