package usecase

import (
	"budget-book/internal/domain/entity"
	"budget-book/internal/infrastructure/repository"
)

type CategoryUseCase struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryUseCase(categoryRepo *repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (uc *CategoryUseCase) CreateCategory(name string, categoryType entity.CategoryType, color string) (*entity.Category, error) {
	category := entity.NewCategory(name, categoryType, color)
	if err := uc.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (uc *CategoryUseCase) GetCategoryByID(id uint64) (*entity.Category, error) {
	return uc.categoryRepo.GetByID(id)
}

func (uc *CategoryUseCase) GetAllCategories() ([]*entity.Category, error) {
	return uc.categoryRepo.GetAll()
}

func (uc *CategoryUseCase) GetCategoriesByType(categoryType entity.CategoryType) ([]*entity.Category, error) {
	return uc.categoryRepo.GetByType(categoryType)
}

func (uc *CategoryUseCase) UpdateCategory(id uint64, name string, categoryType entity.CategoryType, color string) (*entity.Category, error) {
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

func (uc *CategoryUseCase) DeleteCategory(id uint64) error {
	_, err := uc.categoryRepo.GetByID(id)
	if err != nil {
		return err
	}

	return uc.categoryRepo.Delete(id)
}