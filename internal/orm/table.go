package orm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// getTableName extracts the table name from a struct type using reflection.
func GetTableName(model any) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem() // Dereference the pointer if it is a pointer.
	}

	// Use the struct name as the table name by default.  You might want to
	// add a tag like `liteforge:"table_name"` to override this.
	return strings.ToLower(t.Name()) // Convert to lowercase as a convention
}

// getFieldInfo extracts field information from a struct using reflection.
// It returns slices of column names and placeholders for use in SQL queries.
func GetFieldInfo(model any) ([]string, []interface{}) {
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	t := val.Type()

	numFields := val.NumField()
	columns := make([]string, 0, numFields)
	values := make([]interface{}, 0, numFields)

	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		columnName := strings.ToLower(field.Name) // Default column name

		// // Check for a `db` tag to customize the column name.
		// dbTag := field.Tag.Get("db")
		// if dbTag != "" {
		// 	columnName = dbTag
		// }
		columns = append(columns, columnName)
		values = append(values, val.Field(i).Interface())
	}

	return columns, values
}

// GetPrimaryKeyColumn extracts the primary key column name from a struct type.
func GetPrimaryKeyColumn(model any) (string, error) {
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	t := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		pkTag := field.Tag.Get("pk")
		if pkTag == "true" {
			// Default column name is lowercase field name
			return strings.ToLower(field.Name), nil
		}
	}

	return "", errors.New("model has no primary key field (tag: `pk:\"true\"`)")
}

// CreateTable creates a database table based on the provided model using the Datastore's adapter.
func CreateTable(ds *Datastore, model any) error {

	//Check for invalid inputs
	if ds == nil || ds.DB == nil || ds.Adapter == nil {
		return errors.New("datastore, database connection, or adapter was nil")
	}

	createQuery, err := ds.Adapter.CreateTableSQL(model)
	if err != nil {
		return fmt.Errorf("failed to generate create table SQL: %w", err)
	}

	_, err = ds.DB.Exec(createQuery)
	if err != nil {
		return fmt.Errorf("failed to execute create table query: %w", err)
	}

	return nil
}
