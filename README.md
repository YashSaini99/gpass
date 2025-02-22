# Graphical Password Authentication ![Go](https://img.shields.io/badge/Go-100%25-blue) ![MongoDB](https://img.shields.io/badge/MongoDB-Database-green) ![MIT License](https://img.shields.io/badge/License-MIT-yellow.svg)

Graphical Password Authentication is a Go package that secures user login with image-based password patterns. It converts selected image indices into a string, hashes it with bcrypt, and stores it in MongoDB. It also features brute-force protection, email alerts, and secure password resets.

## Features

- 🔒 Secure user login with image-based password patterns
- 🛡️ Brute-force protection
- 📧 Email alerts for suspicious activities
- 🔄 Secure password resets
- 💾 Stores hashed passwords in MongoDB

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Basic Authentication](#basic-authentication)
  - [Advanced Security Features](#advanced-security-features)
  - [Email Validation](#email-validation)
  - [Sending Emails](#sending-emails)
- [API Reference](#api-reference)
  - [Core Functions](#core-functions)
  - [Advanced Security Functions](#advanced-security-functions)
- [Example](#example)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install the package, use:

```bash
go get github.com/YashSaini99/graphical-password-authentication
```

## Configuration

Create a `.env` file in the root of your project with the following keys:

```ini
# Database Configuration
DB_URI=mongodb://localhost:27017/graphicalpasswordauth

# SMTP Configuration (example using Mailtrap for testing)
SMTP_USER=your_mailtrap_username@mailtrap.io
SMTP_PASS=your_mailtrap_password
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
```

- **DB_URI:** Connection string for your MongoDB instance.
- **SMTP_USER, SMTP_PASS, SMTP_HOST, SMTP_PORT:** Credentials and server details for sending emails. You can use a service like Mailtrap for testing purposes.

## Usage

### Basic Authentication

The package uses a facade package `gp` to simplify imports. In your code, you can do:

```go
import (
    "github.com/YashSaini99/graphical-password-authentication/gp"
    "time"
)

func main() {
    // Load environment variables
    gp.LoadEnv()

    // Connect to the database
    err := gp.Connect("your_mongodb_connection_string")
    if err != nil {
        // Handle error
    }
    defer gp.Disconnect()

    // Validate an email
    if !gp.IsValidEmail("user@example.com") {
        // Handle invalid email
    }

    // Register a new user
    err = gp.RegisterUser("username", "user@example.com", []int{1, 3, 5, 7})
    if err != nil {
        // Handle error (e.g., duplicate username/email)
    }

    // Authenticate the user
    ok, err := gp.AuthenticateUser("username", []int{1, 3, 5, 7})
    if err != nil {
        // Handle error
    }
    if ok {
        // Successful login
    }
}
```

### Advanced Security Features

For added security, use the advanced functions that protect against brute-force attacks and support password resets.

```go
// Create a SecureAuthManager instance
secManager := gp.NewSecureAuthManager(3, 10*time.Minute, 15*time.Minute)

// Authenticate with protection (this will block the account on repeated failed attempts and send alert emails)
ok, err := secManager.AuthenticateWithProtection("username", []int{1, 3, 5, 7}, "user@example.com")
if err != nil {
    // Handle authentication error (e.g., account blocked)
}
if ok {
    // Successful login
}

// Initiate a password reset (generates a secure token and sends a reset email)
token, err := secManager.InitiatePasswordReset("username", "user@example.com")
if err != nil {
    // Handle password reset error
}
// Use the token for resetting the password, typically via a dedicated reset endpoint.
```

### Email Validation

```go
// Validate an email
if gp.IsValidEmail("user@example.com") {
    fmt.Println("Email is valid")
} else {
    fmt.Println("Email is invalid")
}
```

### Sending Emails

```go
// Send an email
err := gp.SendEmail("user@example.com", "Subject", "Email body")
if err != nil {
    // Handle email sending error
}
```

## API Reference

### Core Functions

- **`LoadEnv() error`**  
  Loads environment variables from a `.env` file.

- **`Connect(uri string) error`**  
  Connects to MongoDB using the provided URI.

- **`Disconnect() error`**  
  Disconnects from MongoDB.

- **`RegisterUser(username, email string, graphicalPassword []int) error`**  
  Registers a new user.

- **`AuthenticateUser(username string, graphicalPassword []int) (bool, error)`**  
  Authenticates a user with their graphical password.

- **`IsValidEmail(email string) bool`**  
  Validates an email address.

- **`SendEmail(to, subject, body string) error`**  
  Sends an email using the SMTP settings in your `.env` file.

### Advanced Security Functions

- **`NewSecureAuthManager(threshold int, blockDuration, tokenDuration time.Duration) *SecureAuthManager`**  
  Creates a new instance of SecureAuthManager.

- **`(m *SecureAuthManager) AuthenticateWithProtection(username string, graphicalPassword []int, userEmail string) (bool, error)`**  
  Authenticates a user with brute-force protection.

- **`(m *SecureAuthManager) InitiatePasswordReset(username, userEmail string) (string, error)`**  
  Initiates a password reset, sending a reset email with a secure token.

- **`(m *SecureAuthManager) ValidateResetToken(username, token string) bool`**  
  Validates a password reset token.

## Example

See the `cmd/example/main.go` file for a complete example demonstrating registration, authentication, and password reset flows using the package.

## Testing

To run the tests for this package:

```bash
go test ./tests
```

This will execute unit tests for core functionalities such as hashing, email validation, and more.

## Contributing

Contributions are welcome! If you have ideas for enhancements, bug fixes, or additional features, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/YashSaini99/graphical-password-authentication/blob/main/LICENCE) file for details.

[![GitHub Stars](https://img.shields.io/github/stars/YashSaini99/graphical-password-authentication?style=social)](https://github.com/YashSaini99/graphical-password-authentication)
[![GitHub Issues](https://img.shields.io/github/issues/YashSaini99/graphical-password-authentication?style=plastic)](https://github.com/YashSaini99/graphical-password-authentication/issues)
[![GitHub Forks](https://img.shields.io/github/forks/YashSaini99/graphical-password-authentication?style=social)](https://github.com/YashSaini99/graphical-password-authentication)
