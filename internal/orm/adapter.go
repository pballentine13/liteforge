package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	_ "github.com/lib/pq"           // For PostgreSQL
	_ "github.com/mattn/go-sqlite3" // For SQLite
)

// DBAdapter defines the interface for database-specific operations.
type DBAdapter interface {
	Connect(cfg Config) (*sql.DB, error)
	CreateTableSQL(model interface{}) (string, error)
	GetPlaceholder(index int) string
	Query(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error)
	BeginTx(db *sql.DB) (*sql.Tx, error)
}

// SQLiteAdapter implements the DBAdapter for SQLite.
type SQLiteAdapter struct{}

// Connect establishes a SQLite database connection.
func (a *SQLiteAdapter) Connect(cfg Config) (*sql.DB, error) {
	dbDirPerm := 0755 //Directory permission
	dbPath := cfg.DataSourceName
	dbDir := filepath.Dir(dbPath)

	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		//creating directory
		if err := os.MkdirAll(dbDir, os.FileMode(dbDirPerm)); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	db, err := sql.Open(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	//Set up wal mode, if needed. (SQLite-specific)
	if cfg.UseWriteAheadLogs {
		_, err = db.Exec("PRAGMA journal_mode=WAL;")
		if err != nil {
			return nil, fmt.Errorf("Failed to switch to WAL mode: %w", err)
		}
	}
	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

// CreateTableSQL generates the SQLite-specific CREATE TABLE SQL statement.
func (a *SQLiteAdapter) CreateTableSQL(model interface{}) (string, error) {
	// Check for invalid inputs
	if model == nil {
		return "", errors.New("no model passed in. model was nil")
	}

	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", errors.New("model must be a struct or pointer to struct")
	}
	tableName := GetTableName(model)

	var columnDefinitions []string
	var columnConstraint string

	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	t = val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		columnName := strings.ToLower(field.Name) // Default column name

		// Check for a `db` tag to customize the column name.
		columnConstraint = ""
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			columnConstraint = strings.ToUpper(dbTag)
		}
		fieldType := field.Type.String()
		sqlType := ""

		switch fieldType {
		case "int", "int64", "int32", "int16", "int8":
			sqlType = "INTEGER"
		case "string":
			sqlType = "TEXT"
		case "float64", "float32":
			sqlType = "REAL"
		case "bool":
			sqlType = "BOOLEAN"
		default:
			sqlType = "TEXT" // Default to TEXT if type is unknown
		}
		// Check for primary key tag
		pkTag := field.Tag.Get("pk")
		if pkTag == "true" {
			sqlType += " PRIMARY KEY"
		}

		columnDefinitions = append(columnDefinitions, fmt.Sprintf("%s %s %s", columnName, sqlType, columnConstraint))

	}

	createQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columnDefinitions, ", "))
	return createQuery, nil
}

// GetPlaceholder returns the SQLite placeholder '?'.
func (a *SQLiteAdapter) GetPlaceholder(index int) string {
	return "?"
}

// Query executes a generic query.
func (a *SQLiteAdapter) Query(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return rows, nil
}

// BeginTx starts a database transaction.
func (a *SQLiteAdapter) BeginTx(db *sql.DB) (*sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

// PostgresAdapter implements the DBAdapter for PostgreSQL.
type PostgresAdapter struct{}

// Connect establishes a PostgreSQL database connection.
func (a *PostgresAdapter) Connect(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

// CreateTableSQL generates the PostgreSQL-specific CREATE TABLE SQL statement.
func (a *PostgresAdapter) CreateTableSQL(model interface{}) (string, error) {
	// Check for invalid inputs
	if model == nil {
		return "", errors.New("no model passed in. model was nil")
	}

	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", errors.New("model must be a struct or pointer to struct")
	}
	tableName := GetTableName(model)

	var columnDefinitions []string
	var columnConstraint string

	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	t = val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		columnName := strings.ToLower(field.Name) // Default column name

		// Check for a `db` tag to customize the column name.
		columnConstraint = ""
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			columnConstraint = strings.ToUpper(dbTag)
		}
		fieldType := field.Type.String()
		sqlType := ""

		switch fieldType {
		case "int", "int64", "int32", "int16", "int8":
			sqlType = "INTEGER"
		case "string":
			sqlType = "TEXT"
		case "float64", "float32":
			sqlType = "DOUBLE PRECISION"
		case "bool":
			sqlType = "BOOLEAN"
		default:
			sqlType = "TEXT" // Default to TEXT if type is unknown
		}
		// Check for primary key tag
		pkTag := field.Tag.Get("pk")
		if pkTag == "true" {
			if strings.Contains(sqlType, "INTEGER") {
				sqlType = "SERIAL PRIMARY KEY" // Use SERIAL for auto-incrementing PK
			} else {
				sqlType += " PRIMARY KEY"
			}
		}

		columnDefinitions = append(columnDefinitions, fmt.Sprintf("%s %s %s", columnName, sqlType, columnConstraint))

	}

	createQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columnDefinitions, ", "))
	return createQuery, nil
}

// GetPlaceholder returns the PostgreSQL numbered placeholder '$N'.
func (a *PostgresAdapter) GetPlaceholder(index int) string {
	return fmt.Sprintf("$%d", index)
}

// Query executes a generic query.
func (a *PostgresAdapter) Query(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return rows, nil
}

// BeginTx starts a database transaction.
func (a *PostgresAdapter) BeginTx(db *sql.DB) (*sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}
