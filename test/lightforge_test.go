package lightforge

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pballentine13/liteforge"
)

func TestOpenDB_Unit(t *testing.T) {
	testCases := []struct {
		name        string
		config      liteforge.Config
		expectedErr bool
	}{
		{
			name: "In-Memory Database",
			config: liteforge.Config{
				DriverName:     "sqlite3",
				DataSourceName: ":memory:", // In-memory database for testing
			},
			expectedErr: false,
		},
		{
			name: "File-Based Database",
			config: liteforge.Config{
				DriverName:     "sqlite3",
				DataSourceName: "test.db", // File-based database
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clean up the file-based database after the test (if applicable)
			if tc.name == "File-Based Database" {
				defer os.Remove("test.db") // Remove test.db after the test
			}

			db, err := liteforge.OpenDB(tc.config)

			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				if db != nil {
					db.Close() // Close the connection after the test
				}
			}
		})
	}
}

// TestUser is a sample struct for testing table creation
type TestUser struct {
	ID        int    `db:"id,primarykey,autoincrement"`
	Username  string `db:"username,unique,notnull"`
	Email     string `db:"email,unique"`
	Age       int    `db:"age"`
	IsActive  bool   `db:"is_active,default:true"`
}

var cfg liteforge.Config = liteforge.Config{
	DriverName:			"sqlite3",
	DataSourceName:	":memory:",
}

func TestCreateTable(t *testing.T) {
	// Test cases
	tests := []struct {
		name    string
		model   interface{}
		wantErr bool
	}{
		{
			name:    "Valid struct with tags",
			model:   TestUser{},
			wantErr: false,
		},
		{
			name:    "Nil model",
			model:   nil,
			wantErr: true,
		},
		{
			name:    "Non-struct type",
			model:   "not a struct",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup in-memory SQLite database for testing
			db, err := liteforge.OpenDB(cfg)
			if err != nil {
				t.Fatalf("Failed to open test database: %v", err)
			}
			defer db.Close()

			// Test the CreateTable function
			err = liteforge.CreateTable(db, tt.model)

			// Check if error matches expected outcome
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// For successful cases, verify table creation
			if !tt.wantErr {
				// Verify table exists and has correct schema
				var tableName string
				if tt.name == "Valid struct with tags" {
					tableName = "test_users" // Assuming table name is pluralized lowercase of struct name
					
					// Query SQLite schema for table info
					var count int
					err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master 
						WHERE type='table' AND name=?`, tableName).Scan(&count)
					
					if err != nil {
						t.Errorf("Failed to query table existence: %v", err)
					}
					
					if count != 1 {
						t.Errorf("Table %s was not created", tableName)
					}

					// Verify columns exist with correct types
					rows, err := db.Query(`PRAGMA table_info(` + tableName + `)`)
					if err != nil {
						t.Errorf("Failed to query table schema: %v", err)
					}
					defer rows.Close()

					// Map to store expected columns and their types
					expectedColumns := map[string]string{
						"id":        "INTEGER",
						"username":  "TEXT",
						"email":     "TEXT",
						"age":       "INTEGER",
						"is_active": "BOOLEAN",
					}

					// Verify each column
					for rows.Next() {
						var (
							cid          int
							name         string
							columnType   string
							notNull      bool
							defaultVal   interface{}
							primaryKey   bool
						)
						
						if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultVal, &primaryKey); err != nil {
							t.Errorf("Failed to scan column info: %v", err)
						}

						expectedType, exists := expectedColumns[name]
						if !exists {
							t.Errorf("Unexpected column %s in table", name)
						} else if columnType != expectedType {
							t.Errorf("Column %s has type %s, expected %s", name, columnType, expectedType)
						}
					}
				}
			}
		})
	}
}
