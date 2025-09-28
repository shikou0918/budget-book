package usecase

import (
	"budget-book/entity"
)

// CategoryUseCase handles category business logic
type CategoryUseCase struct {
	categoryRepo CategoryRepositoryInterface
}

// NewCategoryUseCase creates a new category use case instance
func NewCategoryUseCase(categoryRepo CategoryRepositoryInterface) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory creates a new category with validation
func (uc *CategoryUseCase) CreateCategory(name string, categoryType entity.TransactionType, color string) (*entity.Category, error) {
	category := entity.NewCategory(name, categoryType, color)
	if err := uc.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID retrieves a category by its ID
func (uc *CategoryUseCase) GetCategoryByID(id uint64) (*entity.Category, error) {
	return uc.categoryRepo.GetByID(id)
}

// GetAllCategories retrieves all categories
func (uc *CategoryUseCase) GetAllCategories() ([]*entity.Category, error) {
	return uc.categoryRepo.GetAll()
}

// GetCategoriesByType retrieves categories by their type
func (uc *CategoryUseCase) GetCategoriesByType(categoryType entity.TransactionType) ([]*entity.Category, error) {
	return uc.categoryRepo.GetByType(categoryType)
}

// UpdateCategory updates an existing category with validation
func (uc *CategoryUseCase) UpdateCategory(id uint64, name string, categoryType entity.TransactionType, color string) (*entity.Category, error) {
	category, err := uc.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	category.Name = name
	category.Type = categoryType
	category.Color = color

	if err := uc.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category by its ID
func (uc *CategoryUseCase) DeleteCategory(id uint64) error {
	_, err := uc.categoryRepo.GetByID(id)
	if err != nil {
		return err
	}

	return uc.categoryRepo.Delete(id)
}
