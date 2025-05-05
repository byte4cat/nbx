# Usage Example

## Application Initialization
During your application's setup (e.g., in `main` or an `init` function), you
initialize your base GORM DB connection, the `TransactionManager`, and inject dependencies.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
	"gorm.io/driver/sqlite" // Example driver

	"yourpackage/transaction" // Replace with your package path
	"yourpackage/model"      // Replace with your model/repository package path
)

// Assume these are defined elsewhere
// type UserRepository interface { Create(ctx context.Context, user *model.User) error }
// type AccountRepository interface { Create(ctx context.Context, account *model.Account) error }
// type UserService struct { ... }
// func NewUserRepository(db *gorm.DB) model.UserRepository { ... }
// func NewAccountRepository(db *gorm.DB) model.AccountRepository { ... }
// func NewUserService(txManager transaction.TransactionManager, userRepo model.UserRepository, accountRepo model.AccountRepository) *UserService { ... }

var baseDB *gorm.DB
var appTxManager transaction.TransactionManager
var appUserService *UserService // Your service instance

func init() {
	var err error
	// 1. Initialize the base GORM DB connection
	baseDB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto-migrate schema (example)
	baseDB.AutoMigrate(&model.User{}, &model.Account{})

	// 2. Create the Transaction Manager instance
	appTxManager = transaction.NewGormTransactionManager(baseDB)

	// 3. Create Repository instances (they hold the base DB)
	userRepo := NewUserRepository(baseDB)
	accountRepo := NewAccountRepository(baseDB)

	// 4. Create Service instance, injecting dependencies
	appUserService = NewUserService(appTxManager, userRepo, accountRepo)

	log.Println("Application initialized successfully.")
}

func main() {
    // Example usage of the service layer method that uses transactions
    ctx := context.Background()
    err := appUserService.RegisterUserWithInitialBalance(ctx, "John Doe", "john.doe@example.com", 100.0)
    if err != nil {
        log.Printf("Error registering user: %v", err)
    } else {
        log.Println("User registered successfully.")
    }

	// Example of a service method that might NOT need a transaction
	// user, err := appUserService.GetUserByID(ctx, 1) // This method might not use ExecuteInTx
}

// Example Service and Repository definitions (implementations would go elsewhere)
// Assume model.User and model.Account structs are defined.

// Example UserRepository
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	// Use transaction.GetDB to get the appropriate DB instance
	db := transaction.GetDB(ctx, r.db)
	// Execute GORM operation using the obtained DB instance
	if result := db.WithContext(ctx).Create(user); result.Error != nil {
		return result.Error
	}
	return nil
}

// Example AccountRepository
type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, account *model.Account) error {
	db := transaction.GetDB(ctx, r.db)
	if result := db.WithContext(ctx).Create(account); result.Error != nil {
		return result.Error
	}
	return nil
}

// Example UserService
type UserService struct {
	txManager transaction.TransactionManager // Dependency on the interface
	userRepo  *UserRepository            // Dependency on Repository implementation (or interface)
	accountRepo *AccountRepository
}

func NewUserService(txManager transaction.TransactionManager, userRepo *UserRepository, accountRepo *AccountRepository) *UserService {
	return &UserService{txManager: txManager, userRepo: userRepo, accountRepo: accountRepo}
}

// Service method that requires a transaction
func (s *UserService) RegisterUserWithInitialBalance(ctx context.Context, name, email string, initialBalance float64) error {
	log.Println("UserService: Starting transaction...")
	// Use the injected TransactionManager to execute the logic in a transaction
	err := s.txManager.ExecuteInTx(ctx, func(txCtx context.Context) error {
		// --- Business Logic inside the transaction ---

		// 1. Create User using the repository (receives txCtx)
		newUser := &model.User{Name: name, Email: email}
		log.Println("UserService: Creating user in transaction...")
		if err := s.userRepo.Create(txCtx, newUser); err != nil { // userRepo.Create uses GetDB(txCtx, ...)
			return fmt.Errorf("user creation failed: %w", err) // Return error to trigger rollback
		}
        // Assume newUser.ID is populated here by GORM

		// 2. Create Account using the repository (receives txCtx)
		newAccount := &model.Account{UserID: newUser.ID, Balance: initialBalance}
		log.Println("UserService: Creating account in transaction...")
		if err := s.accountRepo.Create(txCtx, newAccount); err != nil { // accountRepo.Create uses GetDB(txCtx, ...)
			return fmt.Errorf("account creation failed: %w", err) // Return error to trigger rollback
		}

		log.Println("UserService: Transaction logic successful.")
		return nil // Return nil to trigger commit
		// --- End Business Logic ---
	})

	if err != nil {
		// Handle transaction failure (either Begin, fn error/rollback, or Commit error)
		return fmt.Errorf("RegisterUserWithInitialBalance transaction failed: %w", err)
	}

	log.Println("UserService: Transaction committed.")
	return nil // Success
}

// Example model structs (usually in model/model.go)
/*
package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"uniqueIndex"`
}

type Account struct {
	gorm.Model
	UserID  uint `gorm:"uniqueIndex"`
	Balance float64
}
*/
```

## Service Layer Usage
The service layer method depends on the `TransactionManager` interface and calls
`ExecuteInTx`, passing the original context and a `TxFn` containing the business logic.

```go
// Assume s is an instance of UserService with dependencies injected

// Service method that requires a transaction
func (s *UserService) RegisterUserWithInitialBalance(ctx context.Context, name, email string, initialBalance float64) error {
	// Use the injected TransactionManager to execute the logic in a transaction
	// ExecuteInTx handles Begin/Commit/Rollback
	// ExecuteInTx passes a transaction-aware context (txCtx) to your function
	err := s.txManager.ExecuteInTx(ctx, func(txCtx context.Context) error {
		// --- Business Logic inside the transaction ---

		// Call repository methods, passing the transaction-aware context (txCtx)
		newUser := &model.User{Name: name, Email: email}
		if err := s.userRepo.Create(txCtx, newUser); err != nil {
			// Return the specific error. ExecuteInTx will catch this and rollback.
			return fmt.Errorf("failed during user creation: %w", err)
		}

		// newUser.ID will be populated after successful creation

		newAccount := &model.Account{UserID: newUser.ID, Balance: initialBalance}
		if err := s.accountRepo.Create(txCtx, newAccount); err != nil {
			// Return the specific error. ExecuteInTx will catch this and rollback.
			return fmt.Errorf("failed during account creation: %w", err)
		}

		// If all operations within the TxFn succeed, return nil
		return nil // This signals ExecuteInTx to commit the transaction
		// --- End Business Logic ---
	})

	// Handle the outcome of the transaction execution (commit success or failure/rollback)
	if err != nil {
		// ExecuteInTx has already handled rollback if needed.
		// Return a higher-level error or the wrapped error depending on needs.
		return fmt.Errorf("RegisterUserWithInitialBalance transaction failed: %w", err)
	}

	// Transaction was committed successfully by ExecuteInTx
	return nil
}
```

## Repository Layer Usage
Repository methods receive the context and use `transaction.GetDB` at the
beginning to get the appropriate `*gorm.DB` instance (either the transaction or
the base connection).

```go
// Assume r is an instance of UserRepository holding the base *gorm.DB instance

// Repository method that performs a DB operation
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	// --- Use transaction.GetDB at the start of the method ---
	// GetDB checks if a transaction is in the context.
	// If yes, it returns the transaction (*gorm.Tx).
	// If no, it returns the base *gorm.DB (r.db).
	db := transaction.GetDB(ctx, r.db)

	// --- Use the 'db' instance returned by GetDB for GORM operations ---
	// This code works seamlessly whether called inside or outside a transaction.
	if result := db.WithContext(ctx).Create(user); result.Error != nil {
		return result.Error
	}

	// Assume GORM populates user.ID on success
	return nil
}

// Other repository methods like FindByID, Update, Delete etc. would follow the same pattern:
// func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
//     db := transaction.GetDB(ctx, r.db) // Get DB instance
//     var user model.User
//     if result := db.WithContext(ctx).First(&user, id); result.Error != nil {
//         if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//             return nil, nil // Or a specific ErrNotFound
//         }
//         return nil, result.Error // Return other errors
//     }
//     return &user, nil
// }
```

