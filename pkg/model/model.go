package model

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/pballentine13/liteforge/internal/orm"
)

// Repository defines the high-level, model-centric interface for CRUD operations.
type Repository interface {
	// Save handles both INSERT (if ID is zero/default) and UPDATE (if ID is set).
	Save(model any) (sql.Result, error)
	// FindByID populates the provided model struct with data for the given ID.
	FindByID(model any, id int) error
	// Update explicitly updates an existing record.
	Update(model any) (sql.Result, error)
	// Delete deletes a record based on the model's primary key.
	Delete(model any) (sql.Result, error)
}

// ORMRepository is a concrete implementation of the Repository interface
// that holds a reference to *orm.Datastore.
type ORMRepository struct {
	DS *orm.Datastore
}

// NewORMRepository creates a new ORMRepository instance.
func NewORMRepository(ds *orm.Datastore) *ORMRepository {
	return &ORMRepository{DS: ds}
}

// Save handles both INSERT and UPDATE.
// It checks the primary key value to determine the operation.
func (r *ORMRepository) Save(model any) (sql.Result, error) {
	if r.DS == nil {
		return nil, fmt.Errorf("datastore is nil")
	}

	pkValue, err := orm.GetPrimaryKeyValue(model)
	if err != nil {
		// If no PK is found, default to Insert.
		return orm.Insert(r.DS, model)
	}

	// Check if the PK value is zero/default (e.g., 0 for int).
	v := reflect.ValueOf(pkValue)
	if v.Kind() == reflect.Int || v.Kind() == reflect.Int64 {
		if v.Int() == 0 {
			return orm.Insert(r.DS, model)
		}
	}

	// If PK is set (non-zero int), perform an update.
	return r.Update(model)
}

// FindByID populates the provided model struct with data for the given ID.
func (r *ORMRepository) FindByID(model any, id int) error {
	if r.DS == nil {
		return fmt.Errorf("datastore is nil")
	}

	// Ensure model is a non-nil pointer to a struct
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("model must be a non-nil pointer to a struct")
	}
	v = v.Elem()
	t := v.Type()

	// 1. Get table name and column names
	tableName := orm.GetTableName(model)
	columns, _ := orm.GetFieldInfo(model) // columns are DB column names (lowercase field names)

	if len(columns) == 0 {
		return fmt.Errorf("model has no fields to query")
	}

	// 2. Get primary key column name
	pkCol, err := orm.GetPrimaryKeyColumn(model)
	if err != nil {
		return fmt.Errorf("model must have a primary key field with 'pk' tag: %w", err)
	}

	// 3. Build the query
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = %s",
		strings.Join(columns, ", "),
		tableName,
		pkCol,
		r.DS.Adapter.GetPlaceholder(1),
	)

	// 4. Execute the query
	row, err := orm.QueryRow(r.DS, query, id)
	if err != nil {
		return fmt.Errorf("failed to query row: %w", err)
	}

	// 5. Prepare destination pointers for Scan
	dest := make([]any, len(columns))
	for i, col := range columns {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Safety check: ensure the column name matches the lowercase field name
		if strings.ToLower(field.Name) != col {
			// This should not happen if GetFieldInfo is correct, but is a good safeguard.
			return fmt.Errorf("internal error: column name mismatch for index %d: expected %s, got %s", i, strings.ToLower(field.Name), col)
		}

		if !fieldValue.CanSet() {
			return fmt.Errorf("field %s is not settable", field.Name)
		}
		dest[i] = fieldValue.Addr().Interface()
	}

	// 6. Scan the row
	if err := row.Scan(dest...); err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return fmt.Errorf("failed to scan row into model: %w", err)
	}

	return nil
}

// Update explicitly updates an existing record.
func (r *ORMRepository) Update(model any) (sql.Result, error) {
	if r.DS == nil {
		return nil, fmt.Errorf("datastore is nil")
	}

	// 1. Get table name, columns, and values
	tableName := orm.GetTableName(model)
	columns, values := orm.GetFieldInfo(model)

	// 2. Get primary key column and value
	pkCol, err := orm.GetPrimaryKeyColumn(model)
	if err != nil {
		return nil, fmt.Errorf("model must have a primary key field with 'pk' tag: %w", err)
	}
	pkValue, err := orm.GetPrimaryKeyValue(model)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary key value: %w", err)
	}

	// 3. Build SET clause and new values slice (excluding PK)
	setClauses := make([]string, 0, len(columns))
	updateValues := make([]any, 0, len(values))
	placeholderIndex := 1

	for i, col := range columns {
		if col == pkCol {
			continue // Skip primary key in SET clause
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = %s", col, r.DS.Adapter.GetPlaceholder(placeholderIndex)))
		updateValues = append(updateValues, values[i])
		placeholderIndex++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// 4. Append PK value to the end of the values slice for the WHERE clause
	updateValues = append(updateValues, pkValue)

	// 5. Build the query
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = %s",
		tableName,
		strings.Join(setClauses, ", "),
		pkCol,
		r.DS.Adapter.GetPlaceholder(placeholderIndex), // The last placeholder for the PK value
	)

	// 6. Execute the query
	return orm.Exec(r.DS, query, updateValues...)
}

// Delete deletes a record based on the model's primary key.
func (r *ORMRepository) Delete(model any) (sql.Result, error) {
	if r.DS == nil {
		return nil, fmt.Errorf("datastore is nil")
	}

	// 1. Get table name
	tableName := orm.GetTableName(model)

	// 2. Get primary key column and value
	pkCol, err := orm.GetPrimaryKeyColumn(model)
	if err != nil {
		return nil, fmt.Errorf("model must have a primary key field with 'pk' tag: %w", err)
	}
	pkValue, err := orm.GetPrimaryKeyValue(model)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary key value: %w", err)
	}

	// 3. Build the query
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = %s",
		tableName,
		pkCol,
		r.DS.Adapter.GetPlaceholder(1),
	)

	// 4. Execute the query
	return orm.Exec(r.DS, query, pkValue)
}

// User is a sample model for demonstration purposes.
type User struct {
	ID   int `liteforge:"pk"`
	Name string
	Age  int
}
