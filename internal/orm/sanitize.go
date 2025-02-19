package orm

import "strings"

// SanitizeInput sanitizes user input to prevent SQL injection.
func SanitizeInput(input string) string {
	//*NOT PRODUCTION READY*
	input = strings.ReplaceAll(input, "'", "''")    // Escape single quotes in SQLite
	input = strings.ReplaceAll(input, "\"", "\"\"") // Escape double quotes
	return input
}
