package validator

import "regexp"

// EmailRegex is a regular expression for validating email addresses
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail validates that the given email is properly formatted
func ValidateEmail(email string) bool {
	return EmailRegex.MatchString(email)
}

// IsEmpty checks if a string is empty
func IsEmpty(s string) bool {
	return s == ""
}
