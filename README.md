# Liteforge ORM

A lightweight and flexible ORM for Go, designed for simplicity and ease of use with SQLite and PostgreSQL databases. Liteforge leverages Go's reflection capabilities to provide a clean and efficient way to interact with your database, minimizing boilerplate code and maximizing developer productivity.

## Features

*   **Simple Configuration:** Easy-to-configure database connections for SQLite and PostgreSQL
*   **Schema Generation:** Automatic table creation based on Go struct definitions
*   **Reflection-Based Mapping:** Automatically maps Go struct fields to database columns using reflection and optional `db` tags.
*   **Lightweight:** Minimal dependencies and a focus on performance.
*   **Transactions:** Support for database transactions with `BeginTx`, `Commit`, and `Rollback`.
*   **Prepared Statements:** Built-in protection against SQL injection vulnerabilities.

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
Use Go structs to represent your database tables. Use the db tag to customize the mapping between struct fields and database column names. Use the pk:"true" tag to mark the field as primary key with auto increment.
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

	err := liteforge.CreateTable(db, model.User{})
	if err != nil {
		log.Fatal(err)
	}
}
```
Advanced Usage
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

Input Sanitization: The SanitizeInput function provides minimal sanitization. It's highly recommended to rely on prepared statements for SQL injection prevention and use a dedicated HTML sanitization library (e.g., github.com/microcosm-cc/bluemonday) if you need to sanitize HTML content.

Test Thoroughly: Write comprehensive unit and integration tests to ensure the correctness and reliability of your code.

Contributing
Contributions are welcome! Please submit pull requests with clear descriptions of the changes. Follow the existing coding style and conventions. Be sure to include tests for any new features or bug fixes.

License
MIT License
