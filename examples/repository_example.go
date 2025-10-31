package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/pballentine13/liteforge"
)

// User represents a user model with struct tags for ORM mapping.
// The `pk:"true"` tag indicates the primary key field.
// The `db` tag can specify additional constraints like "not null".
type RepoUser struct {
	ID       int    `db:"not null" pk:"true"`
	Username string `db:"unique not null"`
	Email    string `db:"not null unique"`
	Age      int
	IsActive bool
}

func repoExample() {
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
	err = liteforge.CreateTable(ds, RepoUser{})
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Step 4: Create a Repository instance for model-centric operations.
	// The Repository pattern provides high-level CRUD methods.
	repo := liteforge.NewRepository(ds)

	// Step 5: Demonstrate CRUD operations.

	// CREATE: Insert a new user.
	fmt.Println("=== CREATE: Inserting a new user ===")
	user := &RepoUser{
		Username: "johndoe",
		Email:    "john@example.com",
		Age:      30,
		IsActive: true,
	}
	result, err := repo.Save(user) // Save handles INSERT since ID is 0.
	if err != nil {
		log.Fatalf("Failed to save user: %v", err)
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Failed to get last insert ID: %v", err)
	}
	user.ID = int(lastID) // Set the ID for future operations.
	fmt.Printf("Inserted user with ID: %d\n", user.ID)

	// READ: Find the user by ID.
	fmt.Println("\n=== READ: Finding user by ID ===")
	foundUser := &RepoUser{}
	err = repo.FindByID(foundUser, user.ID)
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	fmt.Printf("Found user: %+v\n", foundUser)

	// UPDATE: Modify the user's email and age.
	fmt.Println("\n=== UPDATE: Updating user ===")
	user.Email = "john.doe@example.com"
	user.Age = 31
	_, err = repo.Save(user) // Save handles UPDATE since ID is non-zero.
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}
	fmt.Println("User updated successfully.")

	// Verify the update by reading again.
	err = repo.FindByID(foundUser, user.ID)
	if err != nil {
		log.Fatalf("Failed to find updated user: %v", err)
	}
	fmt.Printf("Updated user: %+v\n", foundUser)

	// DELETE: Remove the user.
	fmt.Println("\n=== DELETE: Deleting user ===")
	_, err = repo.Delete(user)
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	fmt.Println("User deleted successfully.")

	// Verify deletion by trying to find the user again.
	err = repo.FindByID(foundUser, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User not found after deletion (expected).")
		} else {
			log.Fatalf("Unexpected error finding deleted user: %v", err)
		}
	}

	fmt.Println("\n=== Demo completed successfully! ===")
}
