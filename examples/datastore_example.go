package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/pballentine13/liteforge"
	"github.com/pballentine13/liteforge/internal/orm"
)

// User represents a user model with struct tags for ORM mapping.
// The `pk:"true"` tag indicates the primary key field.
// The `db` tag can specify additional constraints like "not null".
type DatastoreUser struct {
	ID       int    `db:"not null" pk:"true"`
	Username string `db:"unique not null"`
	Email    string `db:"not null unique"`
	Age      int
	IsActive bool
}

func datastoreExample() {
	// Step 1: Configure the database connection.
	// Using SQLite with an in-memory database for simplicity.
	// In a real application, you might use a file-based database.
	cfg := liteforge.Config{
		DriverName:     "sqlite3",
		DataSourceName: ":memory:", // Use ":memory:" for in-memory DB
	}

	// Step 2: Open the database connection and get a Datastore.
	ds, err := liteforge.OpenDB(cfg)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer ds.DB.Close() // Ensure the database connection is closed when done.

	// Step 3: Create the table based on the User struct.
	// This uses reflection to generate the CREATE TABLE SQL.
	err = liteforge.CreateTable(ds, DatastoreUser{})
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Step 4: Demonstrate CRUD operations using the Datastore directly.
	// The Data Store pattern involves lower-level SQL operations.

	// CREATE: Insert a new user using orm.Insert.
	fmt.Println("=== CREATE: Inserting a new user ===")
	user := &DatastoreUser{
		Username: "janedoe",
		Email:    "jane@example.com",
		Age:      25,
		IsActive: true,
	}
	result, err := orm.Insert(ds, user)
	if err != nil {
		log.Fatalf("Failed to insert user: %v", err)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Failed to get last insert ID: %v", err)
	}
	user.ID = int(lastID) // Set the ID for future operations.
	fmt.Printf("Inserted user with ID: %d\n", user.ID)

	// READ: Query the user by ID using orm.QueryRow.
	fmt.Println("\n=== READ: Querying user by ID ===")
	query := "SELECT id, username, email, age, isactive FROM user WHERE id = ?"
	row, err := orm.QueryRow(ds, query, user.ID)
	if err != nil {
		log.Fatalf("Failed to query user: %v", err)
	}
	foundUser := &DatastoreUser{}
	err = row.Scan(&foundUser.ID, &foundUser.Username, &foundUser.Email, &foundUser.Age, &foundUser.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("User not found")
		}
		log.Fatalf("Failed to scan user: %v", err)
	}
	fmt.Printf("Found user: %+v\n", foundUser)

	// UPDATE: Update the user's email and age using orm.Exec.
	fmt.Println("\n=== UPDATE: Updating user ===")
	updateQuery := "UPDATE user SET email = ?, age = ? WHERE id = ?"
	_, err = orm.Exec(ds, updateQuery, "jane.doe@example.com", 26, user.ID)
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}
	fmt.Println("User updated successfully.")

	// Verify the update by querying again.
	row, err = orm.QueryRow(ds, query, user.ID)
	if err != nil {
		log.Fatalf("Failed to query updated user: %v", err)
	}
	err = row.Scan(&foundUser.ID, &foundUser.Username, &foundUser.Email, &foundUser.Age, &foundUser.IsActive)
	if err != nil {
		log.Fatalf("Failed to scan updated user: %v", err)
	}
	fmt.Printf("Updated user: %+v\n", foundUser)

	// DELETE: Delete the user using orm.Exec.
	fmt.Println("\n=== DELETE: Deleting user ===")
	deleteQuery := "DELETE FROM user WHERE id = ?"
	_, err = orm.Exec(ds, deleteQuery, user.ID)
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	fmt.Println("User deleted successfully.")

	// Verify deletion by querying again.
	row, err = orm.QueryRow(ds, query, user.ID)
	if err != nil {
		log.Fatalf("Failed to query deleted user: %v", err)
	}
	err = row.Scan(&foundUser.ID, &foundUser.Username, &foundUser.Email, &foundUser.Age, &foundUser.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User not found after deletion (expected).")
		} else {
			log.Fatalf("Unexpected error querying deleted user: %v", err)
		}
	}

	fmt.Println("\n=== Demo completed successfully! ===")
}
