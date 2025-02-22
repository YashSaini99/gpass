// Package gp provides graphical password authentication functionality.
// This package no longer acts as a facade.
//But it imports all the functions and you can use it directly using this package
package gpass

import (
	"github.com/YashSaini99/gpass/auth"
	"github.com/YashSaini99/gpass/config"
	"github.com/YashSaini99/gpass/db"
	"github.com/YashSaini99/gpass/email"
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
