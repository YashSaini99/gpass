// Package gp is a facade that re-exports functionality from the graphical password authentication system.
// It provides a single import point for users.
package gp

import (
	"github.com/YashSaini99/graphical-password-authentication/auth"
	"github.com/YashSaini99/graphical-password-authentication/config"
	"github.com/YashSaini99/graphical-password-authentication/db"
	"github.com/YashSaini99/graphical-password-authentication/email"
)

// Re-export auth functions and types.
var (
	RegisterUser               = auth.RegisterUser
	AuthenticateUser           = auth.AuthenticateUser
	HashGraphicalPassword      = auth.HashGraphicalPassword
	CheckGraphicalPasswordHash = auth.CheckGraphicalPasswordHash
	IsValidEmail               = auth.IsValidEmail
	NewSecureAuthManager       = auth.NewSecureAuthManager
)

type SecureAuthManager = auth.SecureAuthManager

// Re-export configuration functions.
var (
	LoadEnv = config.LoadEnv
)

// Re-export database functions.
var (
	Connect       = db.Connect
	Disconnect    = db.Disconnect
	GetCollection = db.GetCollection
)

// Re-export email functions.
var (
	SendEmail = email.SendEmail
)
