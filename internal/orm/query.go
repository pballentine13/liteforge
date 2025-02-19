package orm

import (
	"database/sql"
	"fmt"
)

// Query performs a custom SQL query.
func Query(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return rows, nil // The caller is responsible for closing the rows.
}

// Exec performs a custom SQL execution (INSERT, UPDATE, DELETE).
func Exec(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare exec statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute exec statement: %w", err)
	}

	return result, nil
}
