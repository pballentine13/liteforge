package orm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// SanitizeInput sanitizes user input to prevent SQL injection.
// DEPRECATED: This function is provided for legacy purposes and should not be used in new code.
// Always prefer prepared statements to prevent SQL injection.
func SanitizeInput(input string) string {
	//*NOT PRODUCTION READY*
	input = strings.ReplaceAll(input, "'", "''")    // Escape single quotes in SQLite
	input = strings.ReplaceAll(input, "\"", "\"\"") // Escape double quotes
	return input
}

// SanitizeHTML removes potentially malicious HTML, allowing only a safe subset.
func SanitizeHTML(htmlInput string) string {
	p := bluemonday.UGCPolicy() // User Generated Content policy is a good starting point
	return p.Sanitize(htmlInput)
}

var (
	// Example: only allow alphanumeric characters and underscores
	alphaNumericRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

// ValidateAndSanitizeAlphaNumeric validates that a string is alphanumeric.
// It returns the original string if valid, and an error otherwise.
func ValidateAndSanitizeAlphaNumeric(input string) (string, error) {
	if !alphaNumericRegex.MatchString(input) {
		return "", fmt.Errorf("input contains invalid characters")
	}
	return input, nil
}

// ValidateAndSanitizeInt validates that a string is a valid integer.
// It returns the integer value if valid, and an error otherwise.
func ValidateAndSanitizeInt(input string) (int, error) {
	return strconv.Atoi(input)
}
