package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// Create performs an INSERT operation.
func Create(db *sql.DB, model interface{}) error {
	tableName := getTableName(model)
	columns, placeholders := getFieldInfo(model)

	insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	stmt, err := db.Prepare(insertQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	args := make([]interface{}, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		args[i] = val.Field(i).Interface()
	}

	_, err = stmt.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to execute insert statement: %w", err)
	}

	return nil
}

// Get retrieves a record by ID.
func Get(db *sql.DB, table string, id int, model interface{}) error {
	tableName := table // Use the provided table name.

	// Get column names from the model using reflection
	columns, _ := getFieldInfo(model) // only need columns

	selectQuery := fmt.Sprintf("SELECT %s FROM %s WHERE id = ?", strings.Join(columns, ", "), tableName)

	stmt, err := db.Prepare(selectQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare select statement: %w", err)
	}
	defer stmt.Close()

	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	dest := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		dest[i] = val.Field(i).Addr().Interface() // Pass pointers to the fields for Scan
	}

	row := stmt.QueryRow(id)
	err = row.Scan(dest...)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("record not found with id %d", id)
		}
		return fmt.Errorf("failed to scan row: %w", err)
	}

	return nil
}

// Update updates a record.
func Update(db *sql.DB, table string, id int, data interface{}) error {
	tableName := table
	columns, _ := getFieldInfo(data)

	var setClauses []string
	for _, column := range columns {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", column))
	}

	updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", tableName, strings.Join(setClauses, ", "))

	stmt, err := db.Prepare(updateQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	args := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		args[i] = val.Field(i).Interface()
	}

	args = append(args, id) // Add the ID to the arguments for the WHERE clause

	_, err = stmt.Exec(args...)
	if err != nil {
		return fmt.Errorf("failed to execute update statement: %w", err)
	}

	return nil
}

// Delete deletes a record by ID.
func Delete(db *sql.DB, table string, id int) error {
	tableName := table

	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)

	stmt, err := db.Prepare(deleteQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}

	return nil
}
