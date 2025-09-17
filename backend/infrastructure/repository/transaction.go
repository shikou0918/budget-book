package repository

import (
	"budget-book/domain/entity"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TransactionRepository handles transaction data operations
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository instance
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create saves a new transaction to the database
func (r *TransactionRepository) Create(transaction *entity.Transaction) error {
	if err := transaction.IsValid(); err != nil {
		return err
	}

	result := r.db.Create(transaction)
	if result.Error != nil {
		return fmt.Errorf("failed to create transaction: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a transaction by its ID
func (r *TransactionRepository) GetByID(id uint64) (*entity.Transaction, error) {
	var transaction entity.Transaction
	result := r.db.Preload("Category").First(&transaction, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, entity.NewNotFoundError("transaction", id)
		}
		return nil, fmt.Errorf("failed to get transaction: %w", result.Error)
	}

	return &transaction, nil
}

// GetAll retrieves all transactions ordered by date and creation time
func (r *TransactionRepository) GetAll() ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	result := r.db.Preload("Category").Order("transaction_date DESC, created_at DESC").Find(&transactions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", result.Error)
	}

	return transactions, nil
}

// GetByDateRange retrieves transactions within a specific date range
func (r *TransactionRepository) GetByDateRange(startDate, endDate time.Time) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	result := r.db.Preload("Category").
		Where("transaction_date >= ? AND transaction_date <= ?", startDate, endDate).
		Order("transaction_date DESC, created_at DESC").
		Find(&transactions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get transactions by date range: %w", result.Error)
	}

	return transactions, nil
}

// GetByCategory retrieves all transactions for a specific category
func (r *TransactionRepository) GetByCategory(categoryID uint64) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	result := r.db.Preload("Category").
		Where("category_id = ?", categoryID).
		Order("transaction_date DESC, created_at DESC").
		Find(&transactions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get transactions by category: %w", result.Error)
	}

	return transactions, nil
}

// GetByMonth retrieves all transactions for a specific year and month
func (r *TransactionRepository) GetByMonth(year, month int) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	result := r.db.Preload("Category").
		Where("YEAR(transaction_date) = ? AND MONTH(transaction_date) = ?", year, month).
		Order("transaction_date DESC, created_at DESC").
		Find(&transactions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get transactions by month: %w", result.Error)
	}

	return transactions, nil
}

// Update modifies an existing transaction in the database
func (r *TransactionRepository) Update(transaction *entity.Transaction) error {
	if err := transaction.IsValid(); err != nil {
		return err
	}

	transaction.UpdatedAt = time.Now()
	result := r.db.Save(transaction)
	if result.Error != nil {
		return fmt.Errorf("failed to update transaction: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return entity.NewNotFoundError("transaction", transaction.ID)
	}

	return nil
}

// Delete removes a transaction from the database by ID
func (r *TransactionRepository) Delete(id uint64) error {
	result := r.db.Delete(&entity.Transaction{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete transaction: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return entity.NewNotFoundError("transaction", id)
	}

	return nil
}
