package usecase

import (
	"budget-book/entity"
	"budget-book/infrastructure/repository"
	"time"
)

// TransactionUseCase handles transaction business logic
type TransactionUseCase struct {
	transactionRepo *repository.TransactionRepository
	categoryRepo    *repository.CategoryRepository
}

// NewTransactionUseCase creates a new transaction use case instance
func NewTransactionUseCase(transactionRepo *repository.TransactionRepository, categoryRepo *repository.CategoryRepository) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
	}
}

// CreateTransaction creates a new transaction with validation
func (uc *TransactionUseCase) CreateTransaction(transactionType entity.TransactionType, amount float64, categoryID uint64, transactionDate time.Time, memo string) (*entity.Transaction, error) {
	category, err := uc.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, err
	}

	if string(category.Type) != string(transactionType) {
		return nil, entity.NewValidationError("transaction type does not match category type")
	}

	transaction := entity.NewTransaction(transactionType, amount, categoryID, transactionDate, memo)
	if err := uc.transactionRepo.Create(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactionByID retrieves a transaction by its ID
func (uc *TransactionUseCase) GetTransactionByID(id uint64) (*entity.Transaction, error) {
	return uc.transactionRepo.GetByID(id)
}

// GetAllTransactions retrieves all transactions
func (uc *TransactionUseCase) GetAllTransactions() ([]*entity.Transaction, error) {
	return uc.transactionRepo.GetAll()
}

// GetTransactionsByDateRange retrieves transactions within a date range
func (uc *TransactionUseCase) GetTransactionsByDateRange(startDate, endDate time.Time) ([]*entity.Transaction, error) {
	return uc.transactionRepo.GetByDateRange(startDate, endDate)
}

// GetTransactionsByCategory retrieves transactions for a specific category
func (uc *TransactionUseCase) GetTransactionsByCategory(categoryID uint64) ([]*entity.Transaction, error) {
	_, err := uc.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, err
	}

	return uc.transactionRepo.GetByCategory(categoryID)
}

// GetTransactionsByMonth retrieves transactions for a specific month
func (uc *TransactionUseCase) GetTransactionsByMonth(year, month int) ([]*entity.Transaction, error) {
	return uc.transactionRepo.GetByMonth(year, month)
}

// UpdateTransaction updates an existing transaction with validation
func (uc *TransactionUseCase) UpdateTransaction(id uint64, transactionType entity.TransactionType, amount float64, categoryID uint64, transactionDate time.Time, memo string) (*entity.Transaction, error) {
	transaction, err := uc.transactionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	category, err := uc.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, err
	}

	if string(category.Type) != string(transactionType) {
		return nil, entity.NewValidationError("transaction type does not match category type")
	}

	transaction.Type = transactionType
	transaction.Amount = amount
	transaction.CategoryID = categoryID
	transaction.TransactionDate = transactionDate
	transaction.Memo = memo

	if err := uc.transactionRepo.Update(transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

// DeleteTransaction deletes a transaction by its ID
func (uc *TransactionUseCase) DeleteTransaction(id uint64) error {
	_, err := uc.transactionRepo.GetByID(id)
	if err != nil {
		return err
	}

	return uc.transactionRepo.Delete(id)
}
