package orm

import (
	"database/sql"
	"fmt"
	"strings"
)

// Query performs a custom SQL query using the Datastore's adapter.
func Query(ds *Datastore, query string, args ...any) (*sql.Rows, error) {
	if ds == nil || ds.DB == nil || ds.Adapter == nil {
		return nil, fmt.Errorf("datastore, database connection, or adapter was nil")
	}
	return ds.Adapter.Query(ds.DB, query, args...)
}

// QueryRow performs a custom SQL query for a single row.
func QueryRow(ds *Datastore, query string, args ...any) (*sql.Row, error) {
	if ds == nil || ds.DB == nil {
		return nil, fmt.Errorf("datastore or database connection was nil")
	}

	stmt, err := ds.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(args...)

	return row, nil
}

// Exec performs a custom SQL execution (INSERT, UPDATE, DELETE).
func Exec(ds *Datastore, query string, args ...any) (sql.Result, error) {
	if ds == nil || ds.DB == nil {
		return nil, fmt.Errorf("datastore or database connection was nil")
	}

	stmt, err := ds.DB.Prepare(query)
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

// postgresResult is a custom sql.Result implementation for PostgreSQL
// to correctly return the last inserted ID via the RETURNING clause.
type postgresResult struct {
	lastInsertID int64
	rowsAffected int64
}

func (r postgresResult) LastInsertId() (int64, error) {
	return r.lastInsertID, nil
}

func (r postgresResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

// Insert performs an INSERT operation for a given model.
func Insert(ds *Datastore, model any) (sql.Result, error) {
	if ds == nil || ds.DB == nil || ds.Adapter == nil {
		return nil, fmt.Errorf("datastore, database connection, or adapter was nil")
	}

	tableName := GetTableName(model)
	allColumns, allValues := GetFieldInfo(model)

	pkCol, err := GetPrimaryKeyColumn(model)
	// If no PK is defined, we insert all fields.
	if err != nil {
		pkCol = ""
	}

	columns := make([]string, 0, len(allColumns))
	values := make([]any, 0, len(allValues))

	for i, col := range allColumns {
		if col == pkCol {
			continue // Skip primary key for auto-increment
		}
		columns = append(columns, col)
		values = append(values, allValues[i])
	}

	placeholders := make([]string, len(columns))
	for i := range columns {
		// Index starts at 1 for SQL placeholders
		placeholders[i] = ds.Adapter.GetPlaceholder(i + 1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	// Check if the adapter is PostgresAdapter to handle ID retrieval
	if _, ok := ds.Adapter.(*PostgresAdapter); ok {
		if pkCol == "" {
			// If no PK is defined, just run a standard Exec
			return Exec(ds, query, values...)
		}

		query += fmt.Sprintf(" RETURNING %s", pkCol)

		var lastInsertID int64
		row := ds.DB.QueryRow(query, values...)
		if err := row.Scan(&lastInsertID); err != nil {
			return nil, fmt.Errorf("failed to execute insert query and scan ID: %w", err)
		}

		// Return a custom result with the retrieved ID
		return postgresResult{lastInsertID: lastInsertID, rowsAffected: 1}, nil
	}

	// For SQLite and other adapters, use standard Exec
	return Exec(ds, query, values...)
}
