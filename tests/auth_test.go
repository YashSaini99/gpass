// Package tests contains tests for the authentication functionality.
package tests

import (
	"testing"

	"github.com/YashSaini99/gpass/auth"
)

func TestGraphicalPasswordHashing(t *testing.T) {
	// Define an example graphical password sequence.
	gpSeq := []int{1, 3, 5, 7}
	hashed, err := auth.HashGraphicalPassword(gpSeq)
	if err != nil {
		t.Fatalf("Error hashing graphical password: %v", err)
	}

	// Verify that the generated hash is valid.
	if err := auth.CheckGraphicalPasswordHash(gpSeq, hashed); err != nil {
		t.Errorf("Graphical password verification failed: %v", err)
	}
}

func TestEmailValidation(t *testing.T) {
	validEmails := []string{"test@example.com", "user.name+tag@example.co.uk"}
	invalidEmails := []string{"plainaddress", "missingatsign.com"}

	for _, email := range validEmails {
		if !auth.IsValidEmail(email) {
			t.Errorf("Expected valid email but got invalid: %s", email)
		}
	}
	for _, email := range invalidEmails {
		if auth.IsValidEmail(email) {
			t.Errorf("Expected invalid email but got valid: %s", email)
		}
	}
}
