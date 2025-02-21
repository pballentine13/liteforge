package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// getTableName extracts the table name from a struct type using reflection.
func getTableName(model interface{}) string {
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
func getFieldInfo(model interface{}) ([]string, []string) {
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	t := val.Type()

	numFields := val.NumField()
	columns := make([]string, 0, numFields)
	placeholders := make([]string, 0, numFields)

	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		columnName := strings.ToLower(field.Name) // Default column name

		// Check for a `db` tag to customize the column name.
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			columnName = dbTag
		}
		columns = append(columns, columnName)
		placeholders = append(placeholders, "?") // SQLite placeholder
	}

	return columns, placeholders
}

// CreateTable creates a database table based on the provided model.
func CreateTable(db *sql.DB, model interface{}) error {
	fmt.Println(model)

	if model == nil {
		return errors.New("no model passed in. model was nil")
	}

	modelKind := reflect.TypeOf(model).Kind()
	if modelKind != reflect.Struct {
		return errors.New("no model passed in")
	}
	tableName := getTableName(model)
	// columns _ := getFieldInfo(model) // We only need the columns for CREATE TABLE.

	var columnDefinitions []string
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		columnName := strings.ToLower(field.Name) // Default column name

		// Check for a `db` tag to customize the column name.
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			columnName = dbTag
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

		columnDefinitions = append(columnDefinitions, fmt.Sprintf("%s %s", columnName, sqlType))

	}

	createQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columnDefinitions, ", "))
	_, err := db.Exec(createQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}
