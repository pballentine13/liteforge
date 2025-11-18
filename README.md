# Liteforge ORM

A lightweight and flexible ORM for Go, designed for simplicity and ease of use with SQLite and PostgreSQL databases. Liteforge leverages Go's reflection capabilities to provide a clean and efficient way to interact with your database, minimizing boilerplate code and maximizing developer productivity.

## Features

*   **Simple Configuration:** Easy-to-configure database connections for SQLite and PostgreSQL
*   **Schema Generation:** Automatic table creation based on Go struct definitions
*   **Reflection-Based Mapping:** Automatically maps Go struct fields to database columns using reflection and optional `db` tags.
*   **Lightweight:** Minimal dependencies and a focus on performance.
*   **Transactions:** Support for database transactions with `BeginTx`, `Commit`, and `Rollback`.
*   **Prepared Statements:** Built-in protection against SQL injection vulnerabilities.
*   **Input Sanitization:** A utility function (`Sanitize`) is provided to protect against Cross-Site Scripting (XSS) in HTML output, powered by Bluemonday.

## Planned Features
*   **Data Stores:** Interface-based data stores for flexible data access patterns (e.g., SQLite, API).

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
Use Go structs to represent your database tables. Use the `db` tag to add column constraints (e.g., "not null", "unique"). Use the `pk:"true"` tag to mark the field as primary key with auto increment.
```go
package model

type User struct {
    ID    int    `pk:"true"`         // Auto-incrementing primary key
    Name  string `db:"not null"`     // Required field
    Email string `db:"unique"`       // Unique constraint
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

	db, err := liteforge.OpenDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ... rest of your code
}
```

### 5. Create the Table
```go
import (
	"log"

	"github.com/pballentine13/liteforge"
    "github.com/pballentine13/pkg/model"
)

func main() {
	// ... (Database connection code from above)

	err := liteforge.CreateTable(db, &model.User{})
	if err != nil {
		log.Fatal(err)
	}
}
```
### 6. Using the Repository for CRUD Operations

The `Repository` provides a high-level, model-centric interface for performing common database operations.

```go
import (
	"log"
	"github.com/pballentine13/liteforge"
	"github.com/pballentine13/liteforge/pkg/model"
)

func main() {
	// ... (Database connection and table creation code from above)

	repo := liteforge.NewRepository(db)

	// 1. INSERT (Save a new user)
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

	// 2. SELECT (Find by ID)
	foundUser := &model.User{}
	err = repo.FindByID(foundUser, newUser.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found user: %+v", foundUser)

	// 3. UPDATE (Save an existing user)
	foundUser.Email = "alice.updated@example.com"
	_, err = repo.Save(foundUser) // Save handles UPDATE when ID is non-zero
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Updated user: %s", foundUser.Email)

	// 4. DELETE
	_, err = repo.Delete(foundUser)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleted user with ID: %d", foundUser.ID)
}
```

### Advanced Usage
Custom Queries
```go
rows, err := liteforge.Query(db, "SELECT * FROM users WHERE name LIKE ?", "%John%")
if err != nil {
    log.Fatal(err)
}
defer rows.Close()

// Process the rows
```

Transactions
```go
tx, err := liteforge.BeginTx(db)
if err != nil {
    log.Fatal(err)
}
defer tx.Rollback() // Rollback if we don't commit

// Perform database operations within the transaction

err = tx.Commit()
if err != nil {
    log.Fatal(err)
}
```

Important Considerations
Security:

Prepared Statements: Liteforge uses prepared statements to prevent SQL injection. Always use prepared statements when handling user input.

Encryption at Rest (SQLCipher): If you enable EncryptAtRest, ensure you store the EncryptionKey securely (e.g., using environment variables or a secrets management solution). Never hardcode the encryption key!

Error Handling: Handle errors gracefully and provide informative error messages.

Database Migrations: For production applications, use a database migration tool (e.g., golang-migrate/migrate) to manage schema changes.

Data Types: The CreateTable function supports basic data types (int, string, float64, bool). Extend it to support other data types as needed.

Relationships: Liteforge does not automatically handle database relationships (one-to-many, many-to-many). You'll need to implement relationship management logic yourself.

Input Sanitization: The `Sanitize` function is designed to prevent Cross-Site Scripting (XSS) attacks by sanitizing user-generated content intended for HTML output. It uses the Bluemonday library to strip potentially dangerous HTML tags and attributes while allowing safe formatting.

**Important:** The `Sanitize` function does **NOT** protect against SQL injection. Always rely on Liteforge's prepared statements for database security. Sanitization is only for output rendering, not input validation.

Example usage:
```go
sanitizedContent := liteforge.Sanitize("<script>alert('XSS')</script><p>Safe text</p>")
// Result: "<p>Safe text</p>"
```

Test Thoroughly: Write comprehensive unit and integration tests to ensure the correctness and reliability of your code.

Contributing
Contributions are welcome! Please submit pull requests with clear descriptions of the changes. Follow the existing coding style and conventions. Be sure to include tests for any new features or bug fixes.

License
MIT License
