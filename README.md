# Liteforge ORM

A lightweight and flexible ORM for Go, designed for simplicity and the **Repository** and **Data Store** patterns. Liteforge provides a clean and efficient way to interact with your database, minimizing boilerplate code and maximizing developer productivity.

## Features

*   **Multi-Database Support:** Seamlessly switch between **SQLite** and **PostgreSQL**.
*   **Model-Centric Repository Pattern:** High-level CRUD operations (`Save`, `FindByID`, `Delete`) using Go structs and the `Repository` interface.
*   **Interface-Based Data Stores:** The recommended pattern for application logic, allowing you to define and implement application-specific data access methods (e.g., `GetUserByID`).
*   **Simple Configuration:** Easy-to-configure database connections.
*   **Schema Generation:** Automatic table creation based on Go struct definitions.
*   **Reflection-Based Mapping:** Automatically maps Go struct fields to database columns.
*   **Transactions:** Support for database transactions.
*   **Prepared Statements:** Built-in protection against SQL injection vulnerabilities.

## Planned Features

*   **Relationships:** Support for defining and querying one-to-many and many-to-many relationships.

## Getting Started

### 1. Installation

```bash
go get github.com/pballentine13/liteforge
```

### 2. Import the Library

```go
import "github.com/pballentine13/liteforge" 
```

### 3. Define Your Data Models

Use Go structs to represent your database tables. Use the `pk:"true"` tag to mark the field as the primary key with auto-increment.

```go
package model

type User struct {
    ID    int    `db:"id" pk:"true"` // Auto-incrementing primary key
    Name  string                     // Maps to the "name" column
    Email string                     // Maps to the "email" column
}
```

### 4. Configure and Open the Database Connection

```go
import (
	"log"
	"os"

	"github.com/joho/godotenv" // Load .env files
	"github.com/pballentine13/liteforge"  
)

func main() {
	// Load environment variables from .env file (if present)
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file: ", err) // Non-fatal error
	}

	dbDriver := os.Getenv("DB_DRIVER")      // e.g., "sqlite3", "postgres"
	dbDataSource := os.Getenv("DB_DATA_SOURCE") // e.g., "mydb.db", "postgres://..."

	cfg := liteforge.Config{
		DriverName:     dbDriver,
		DataSourceName: dbDataSource,
		EncryptAtRest:  false,       // Set to true if using SQLCipher
		EncryptionKey:  "",          // Provide the key if EncryptAtRest is true (DO NOT HARDCODE)
	}

	// Open the database connection. We use 'ds' (DataStore) for the variable name.
	ds, err := liteforge.OpenDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer ds.Close()

	// ... rest of your code
}
```

### 5. Using the Repository

The `Repository` provides a generic, model-centric interface for performing common CRUD operations.

```go
import (
	"log"
	"github.com/pballentine13/liteforge"
	"github.com/pballentine13/liteforge/pkg/model"
)

func main() {
	// ... (Database connection code from above, using 'ds')

	// 1. Create the table
	err := liteforge.CreateTable(ds, model.User{})
	if err != nil {
		log.Fatal(err)
	}

	repo := liteforge.NewRepository(ds)

	// 2. INSERT (Save a new user)
	newUser := &model.User{
		Name:  "Alice",
		Email: "alice@example.com",
	}
	result, err := repo.Save(newUser) // Save handles INSERT when ID is 0
	if err != nil {
		log.Fatal(err)
	}
	newID, _ := result.LastInsertId()
	newUser.ID = int(newID)
	log.Printf("Inserted user with ID: %d", newUser.ID)

	// 3. SELECT (Find by ID)
	foundUser := &model.User{}
	err = repo.FindByID(foundUser, newUser.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found user: %+v", foundUser)

	// 4. UPDATE (Save an existing user)
	foundUser.Email = "alice.updated@example.com"
	_, err = repo.Save(foundUser) // Save handles UPDATE when ID is non-zero
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Updated user: %s", foundUser.Email)

	// 5. DELETE
	_, err = repo.Delete(foundUser)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleted user with ID: %d", foundUser.ID)
}
```

### 6. Using Data Stores (The Recommended Pattern)

The Data Store pattern is recommended for application development. It allows you to define an interface with application-specific methods, keeping your business logic clean and decoupled from generic CRUD.

First, define your Data Store interface (e.g., in `pkg/datastore/user.go`):

```go
package datastore

import "github.com/pballentine13/liteforge/pkg/model"

type UserDataStore interface {
    GetUserByID(id int) (*model.User, error)
    FindUsersByName(name string) ([]*model.User, error)
    SaveUser(user *model.User) error
}
```

Next, implement the Data Store using `liteforge.NewDataStore`. This function uses reflection to implement the interface methods by mapping them to repository calls or custom queries.

```go
// In your main function or initialization logic:
var userDS datastore.UserDataStore

// Create the generic repository (using 'ds' from the connection step)
repo := liteforge.NewRepository(ds)

// Implement the Data Store interface
err = liteforge.NewDataStore(&userDS, repo)
if err != nil {
    log.Fatal(err)
}

// Now use the application-specific Data Store
user, err := userDS.GetUserByID(1)
if err != nil {
    log.Fatal(err)
}
log.Printf("User from Data Store: %+v", user)
```

## Advanced Usage

### Custom Queries

For operations not covered by the Repository, you can execute raw SQL queries using the underlying `Datastore` (`ds`).

```go
// Note: 'ds' is the variable from the connection step (liteforge.DB)
rows, err := liteforge.Query(ds, "SELECT * FROM users WHERE name LIKE ?", "%John%")
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

// Process the rows
```

### Transactions

Transactions are managed using the `liteforge.BeginTx` function, which takes the `Datastore` (`ds`) as an argument.

```go
// Note: 'ds' is the variable from the connection step (liteforge.DB)
tx, err := liteforge.BeginTx(ds)
if err != nil {
    log.Fatal(err)
}
defer tx.Rollback() // Rollback if we don't commit

// Perform database operations within the transaction
// Use 'tx' instead of 'ds' for operations inside the transaction scope

err = tx.Commit()
if err != nil {
    log.Fatal(err)
}
```

## Important Considerations

*   **Security:** Liteforge uses prepared statements to prevent SQL injection. Always use prepared statements when handling user input.
*   **Encryption at Rest (SQLCipher):** If you enable `EncryptAtRest`, ensure you store the `EncryptionKey` securely (e.g., using environment variables or a secrets management solution). Never hardcode the encryption key!
*   **Error Handling:** Handle errors gracefully and provide informative error messages.
*   **Database Migrations:** For production applications, use a database migration tool (e.g., `golang-migrate/migrate`) to manage schema changes.
*   **Input Sanitization:** The `SanitizeInput` function provides minimal sanitization. It's highly recommended to rely on prepared statements for SQL injection prevention and use a dedicated HTML sanitization library (e.g., `github.com/microcosm-cc/bluemonday`) if you need to sanitize HTML content.
*   **Test Thoroughly:** Write comprehensive unit and integration tests to ensure the correctness and reliability of your code.

## Contributing

Contributions are welcome! Please submit pull requests with clear descriptions of the changes. Follow the existing coding style and conventions. Be sure to include tests for any new features or bug fixes.

## License

MIT License