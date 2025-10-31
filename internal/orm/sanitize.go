package orm

import "github.com/microcosm-cc/bluemonday"

// policy is a strict HTML sanitizer policy that only allows basic formatting tags.
// This policy is used by the Sanitize function to prevent XSS attacks.
var policy = bluemonday.StrictPolicy()

func init() {
	// Allow paragraphs, bold, italics, strong, and emphasis tags.
	policy.AllowElements("p", "b", "i", "strong", "em")
}

// Sanitize cleans a string to prevent XSS attacks by stripping unwanted HTML.
//
// IMPORTANT: This function is for preventing Cross-Site Scripting (XSS) and
// SHOULD NOT be used for preventing SQL injection. Use parameterized queries or
// prepared statements to protect against SQL injection.
func Sanitize(input string) string {
	return policy.Sanitize(input)
}
