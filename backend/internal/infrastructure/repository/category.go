package repository

import (
	"budget-book/internal/domain/entity"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CategoryRepository handles category data operations
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository instance
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create saves a new category to the database
func (r *CategoryRepository) Create(category *entity.Category) error {
	if err := category.IsValid(); err != nil {
		return err
	}

	exists, err := r.ExistsByNameAndType(category.Name, category.Type)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("category with name '%s' and type '%s' already exists", category.Name, category.Type)
	}

	result := r.db.Create(category)
	if result.Error != nil {
		return fmt.Errorf("failed to create category: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a category by its ID
func (r *CategoryRepository) GetByID(id uint64) (*entity.Category, error) {
	var category entity.Category
	result := r.db.First(&category, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, entity.NewNotFoundError("category", id)
		}
		return nil, fmt.Errorf("failed to get category: %w", result.Error)
	}

	return &category, nil
}

// GetAll retrieves all categories ordered by type and name
func (r *CategoryRepository) GetAll() ([]*entity.Category, error) {
	var categories []*entity.Category
	result := r.db.Order("type ASC, name ASC").Find(&categories)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get categories: %w", result.Error)
	}

	return categories, nil
}

// GetByType retrieves all categories of a specific type
func (r *CategoryRepository) GetByType(categoryType entity.CategoryType) ([]*entity.Category, error) {
	var categories []*entity.Category
	result := r.db.Where("type = ?", categoryType).Order("name ASC").Find(&categories)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get categories by type: %w", result.Error)
	}

	return categories, nil
}

// Update modifies an existing category in the database
func (r *CategoryRepository) Update(category *entity.Category) error {
	if err := category.IsValid(); err != nil {
		return err
	}

	category.UpdatedAt = time.Now()
	result := r.db.Save(category)
	if result.Error != nil {
		return fmt.Errorf("failed to update category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return entity.NewNotFoundError("category", category.ID)
	}

	return nil
}

// Delete removes a category from the database by ID
func (r *CategoryRepository) Delete(id uint64) error {
	var transactionCount int64
	r.db.Model(&entity.Transaction{}).Where("category_id = ?", id).Count(&transactionCount)
	if transactionCount > 0 {
		return fmt.Errorf("cannot delete category: it is referenced by %d transactions", transactionCount)
	}

	result := r.db.Delete(&entity.Category{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return entity.NewNotFoundError("category", id)
	}

	return nil
}

// ExistsByNameAndType checks if a category exists with the given name and type
func (r *CategoryRepository) ExistsByNameAndType(name string, categoryType entity.CategoryType) (bool, error) {
	var count int64
	result := r.db.Model(&entity.Category{}).Where("name = ? AND type = ?", name, categoryType).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check category existence: %w", result.Error)
	}

	return count > 0, nil
}
