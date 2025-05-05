package transaction

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// TxFn defines the function signature for operations to be executed within a transaction.
// The context passed to this function will contain the transaction instance.
type TxFn func(ctx context.Context) error

// TransactionManager defines the interface for executing operations within a transaction.
// The service layer should depend on this interface.
type TransactionManager interface {
	// ExecuteInTx runs the given function within a transaction.
	// The transaction is committed if fn returns nil, otherwise it's rolled back.
	// The context passed to fn will contain the transaction instance, which repositories
	// should retrieve using GetDB.
	ExecuteInTx(ctx context.Context, fn TxFn) error
}

// GormTransactionManager is a TransactionManager implementation for GORM.
// It holds the base *gorm.DB instance used for starting transactions.
type GormTransactionManager struct {
	db *gorm.DB // GormTransactionManager holds the base DB instance
}

// NewGormTransactionManager creates a new GormTransactionManager.
// It takes the base *gorm.DB instance used for starting transactions.
// This should be called during application initialization.
func NewGormTransactionManager(db *gorm.DB) *GormTransactionManager {
	if db == nil {
		// Panic here as a nil DB is a critical configuration error for a transaction manager.
		panic("transaction: NewGormTransactionManager received nil DB instance")
	}
	return &GormTransactionManager{db: db}
}

// ExecuteInTx implements the TransactionManager interface.
// It wraps the standard GORM Begin/Commit/Rollback logic around the execution of fn.
func (m *GormTransactionManager) ExecuteInTx(ctx context.Context, fn TxFn) error {
	// Start the transaction using the DB instance held by the manager
	tx := m.db.Begin() // Uses m.db to begin transaction
	if tx.Error != nil {
		return fmt.Errorf("transaction: failed to begin: %w", tx.Error)
	}

	// Create a new context containing the transaction instance
	txCtx := SetTx(ctx, tx) // Use SetTx from context.go

	// Execute the provided function within the transaction context
	err := fn(txCtx) // Call the business logic function with the transaction context

	if err != nil {
		// If an error occurred in the business logic function, rollback the transaction.
		if rbErr := tx.Rollback().Error; rbErr != nil {
			// If rollback also fails, return a wrapped error containing both the original
			// business logic error and the rollback error.
			return fmt.Errorf("transaction: %w (rollback failed: %v)", err, rbErr)
		}
		// Return the original error from the business logic function.
		return err
	}

	// If the business logic function returned no error, commit the transaction.
	if cErr := tx.Commit().Error; cErr != nil {
		// If commit fails, attempt a rollback (though it might fail or be redundant).
		// We ignore the error from this rollback attempt as the primary failure is the commit.
		_ = tx.Rollback()
		return fmt.Errorf("transaction: failed to commit: %w", cErr)
	}

	// Transaction successful.
	return nil
}
