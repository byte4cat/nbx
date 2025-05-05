package transaction

import (
	"context"

	"gorm.io/gorm"
)

// GetDB retrieves the appropriate *gorm.DB instance for database operations.
// It checks if a transaction (*gorm.DB) exists within the provided context.
//   - If a transaction is found (placed there by TransactionManager.ExecuteInTx),
//     it returns the transaction instance (*gorm.Tx, which is *gorm.DB).
//   - If no transaction is found in the context, it returns the original *gorm.DB
//     instance passed as an argument, allowing for non-transactional operations.
//
// This function simplifies repository logic, allowing repository methods to work
// seamlessly both within and outside a transaction without needing explicit
// Context checks everywhere. Repository methods should call GetDB(ctx, repo.db)
// at the start.
func GetDB(ctx context.Context, originalDB *gorm.DB) *gorm.DB {
	if tx, ok := GetTx(ctx); ok {
		// Use the transaction DB if available in context
		// GetTx is from context.go
		return tx
	}
	// Fallback to the original DB connection if no transaction in context
	if originalDB == nil {
		// This indicates a potential issue with dependency injection or
		// calling a repository method without a valid DB instance available.
		panic("transaction.GetDB: originalDB is nil and no transaction found in context")
	}
	return originalDB
}
