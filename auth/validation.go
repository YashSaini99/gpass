// Package auth provides helper functions including email validation.
package auth

import (
	"net/mail"
)

// IsValidEmail returns true if the provided email address is valid; false otherwise.
// It uses the net/mail package to parse and verify the email address format.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
