package lightforge

import (
	"strings"
	"testing"

	"github.com/pballentine13/liteforge/internal/orm"
)

func TestAdapterGetPlaceholder(t *testing.T) {
	tests := []struct {
		name     string
		adapter  orm.DBAdapter
		index    int
		expected string
	}{
		{"SQLite", &orm.SQLiteAdapter{}, 1, "?"},
		{"SQLite", &orm.SQLiteAdapter{}, 5, "?"},
		{"Postgres", &orm.PostgresAdapter{}, 1, "$1"},
		{"Postgres", &orm.PostgresAdapter{}, 2, "$2"},
		{"Postgres", &orm.PostgresAdapter{}, 10, "$10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.adapter.GetPlaceholder(tt.index)
			if actual != tt.expected {
				t.Errorf("GetPlaceholder(%d) got = %s, want %s", tt.index, actual, tt.expected)
			}
		})
	}
}

func TestAdapterCreateTableSQL(t *testing.T) {
	userModel := TestUser{}

	tests := []struct {
		name     string
		adapter  orm.DBAdapter
		expected string
	}{
		{
			name:    "SQLite Create Table",
			adapter: &orm.SQLiteAdapter{},
			// Expected: CREATE TABLE IF NOT EXISTS testuser (id INTEGER PRIMARY KEY NOT NULL, username TEXT UNIQUE NOT NULL, email TEXT NOT NULL UNIQUE, age INTEGER , isactive BOOLEAN )
			expected: "CREATE TABLE IF NOT EXISTS testuser (id INTEGER PRIMARY KEY NOT NULL, username TEXT UNIQUE NOT NULL, email TEXT NOT NULL UNIQUE, age INTEGER , isactive BOOLEAN )",
		},
		{
			name:    "Postgres Create Table",
			adapter: &orm.PostgresAdapter{},
			// Expected: CREATE TABLE IF NOT EXISTS testuser (id SERIAL PRIMARY KEY NOT NULL, username TEXT UNIQUE NOT NULL, email TEXT NOT NULL UNIQUE, age INTEGER , isactive BOOLEAN )
			expected: "CREATE TABLE IF NOT EXISTS testuser (id SERIAL PRIMARY KEY NOT NULL, username TEXT UNIQUE NOT NULL, email TEXT NOT NULL UNIQUE, age INTEGER , isactive BOOLEAN )",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.adapter.CreateTableSQL(userModel)
			if err != nil {
				t.Fatalf("CreateTableSQL returned an error: %v", err)
			}
			// Normalize whitespace for comparison
			actual = strings.Join(strings.Fields(actual), " ")
			expected := strings.Join(strings.Fields(tt.expected), " ")

			if actual != expected {
				t.Errorf("CreateTableSQL mismatch.\nGot:      %s\nExpected: %s", actual, expected)
			}
		})
	}
}
