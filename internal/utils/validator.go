package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail validates an email address
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword validates a password (minimum 8 characters, at least one letter and one number)
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasLetter && hasNumber
}

// SanitizeString removes leading/trailing whitespace
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// IsEmpty checks if a string is empty after trimming
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
