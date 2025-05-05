package transaction

import (
	"context"

	"gorm.io/gorm"
)

// contextKey is an unexported type used as key for context values.
// This prevents collisions with keys defined in other packages.
type contextKey string

const (
	// txKey is the key used to store the gorm.DB transaction instance in the context.
	txKey contextKey = "gorm_transaction"
)

// SetTx stores the GORM transaction database instance (*gorm.DB) within the context.
// It returns a new context containing the transaction.
func SetTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

// GetTx retrieves the GORM transaction database instance (*gorm.DB) from the context.
// It returns the transaction instance and a boolean indicating whether a transaction was found.
// If no transaction is found in the context, it returns (nil, false).
func GetTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey).(*gorm.DB)
	if !ok || tx == nil {
		return nil, false // Ensure tx is not nil even if key exists with nil value
	}
	return tx, true
}
