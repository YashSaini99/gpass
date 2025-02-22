// Package auth provides advanced security features such as brute-force protection and password reset capabilities.
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/YashSaini99/gpass/email"
)

// Predefined error for blocked accounts.
var (
	ErrAccountBlocked = fmt.Errorf("account is temporarily blocked")
)

// SecureAuthManager manages brute-force protection and password reset tokens.
type SecureAuthManager struct {
	// failedAttempts tracks the number of consecutive failed login attempts per username.
	failedAttempts map[string]int
	// blockUntil stores the time until which a username is blocked.
	blockUntil map[string]time.Time
	// resetTokens stores password reset tokens with expiration for each username.
	resetTokens map[string]resetTokenInfo

	mu sync.Mutex

	// attemptThreshold is the maximum number of failed attempts before blocking.
	attemptThreshold int
	// blockDuration is the duration for which a user is blocked after exceeding the threshold.
	blockDuration time.Duration
	// tokenDuration is the validity duration of a password reset token.
	tokenDuration time.Duration
}

// resetTokenInfo holds a reset token and its expiration time.
type resetTokenInfo struct {
	token     string
	expiresAt time.Time
}

// NewSecureAuthManager creates a new SecureAuthManager with the specified settings.
func NewSecureAuthManager(threshold int, blockDuration, tokenDuration time.Duration) *SecureAuthManager {
	return &SecureAuthManager{
		failedAttempts:   make(map[string]int),
		blockUntil:       make(map[string]time.Time),
		resetTokens:      make(map[string]resetTokenInfo),
		attemptThreshold: threshold,
		blockDuration:    blockDuration,
		tokenDuration:    tokenDuration,
	}
}

// AuthenticateWithProtection wraps AuthenticateUser with brute-force protection.
// If authentication fails repeatedly, the account is blocked and an alert email is sent.
func (m *SecureAuthManager) AuthenticateWithProtection(username string, graphicalPassword []int, userEmail string) (bool, error) {
	m.mu.Lock()
	// Check if the account is blocked.
	if unblockTime, exists := m.blockUntil[username]; exists && time.Now().Before(unblockTime) {
		m.mu.Unlock()
		return false, fmt.Errorf("account is temporarily blocked until %s", unblockTime.Format(time.RFC1123))
	}
	m.mu.Unlock()

	// Call core authentication.
	ok, err := AuthenticateUser(username, graphicalPassword)
	if err != nil || !ok {
		m.mu.Lock()
		m.failedAttempts[username]++
		attempts := m.failedAttempts[username]
		// If threshold is reached, block the account.
		if attempts >= m.attemptThreshold {
			blockUntil := time.Now().Add(m.blockDuration)
			m.blockUntil[username] = blockUntil
			alertMsg := fmt.Sprintf("Multiple failed login attempts detected for your account. Your account is temporarily blocked until %s.", blockUntil.Format(time.RFC1123))
			// Send alert email asynchronously.
			go func(emailAddr, message string) {
				if err := email.SendEmail(emailAddr, "Alert: Suspicious Login Attempts", message); err != nil {
					fmt.Printf("Error sending alert email: %v\n", err)
				}
			}(userEmail, alertMsg)
		}
		m.mu.Unlock()
		return false, err
	}

	// On successful login, reset the failed attempt counter.
	m.mu.Lock()
	m.failedAttempts[username] = 0
	delete(m.blockUntil, username)
	m.mu.Unlock()
	return true, nil
}

// InitiatePasswordReset generates a secure reset token for the given username,
// stores it, and sends a reset link to the user's email.
// Returns the generated token or an error.
func (m *SecureAuthManager) InitiatePasswordReset(username, userEmail string) (string, error) {
	token, err := generateSecureToken(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	m.mu.Lock()
	m.resetTokens[username] = resetTokenInfo{
		token:     token,
		expiresAt: time.Now().Add(m.tokenDuration),
	}
	m.mu.Unlock()

	// Construct the reset link (configure your domain as needed).
	resetLink := fmt.Sprintf("https://yourdomain.com/reset?username=%s&token=%s", username, token)
	emailBody := fmt.Sprintf("Click the link below to reset your graphical password. This link is valid for %d minutes:\n%s", int(m.tokenDuration.Minutes()), resetLink)
	if err := email.SendEmail(userEmail, "Password Reset Request", emailBody); err != nil {
		return "", fmt.Errorf("failed to send reset email: %w", err)
	}

	return token, nil
}

// ValidateResetToken checks if the provided token for the given username is valid and not expired.
func (m *SecureAuthManager) ValidateResetToken(username, token string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	info, exists := m.resetTokens[username]
	if !exists || time.Now().After(info.expiresAt) {
		delete(m.resetTokens, username)
		return false
	}
	return info.token == token
}

// generateSecureToken generates a secure random token of the specified byte length,
// returning it as a hexadecimal string.
func generateSecureToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
