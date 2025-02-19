package orm

import (
	"database/sql"
	"fmt"
)

// BeginTx starts a database transaction.
func BeginTx(db *sql.DB) (*sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}
