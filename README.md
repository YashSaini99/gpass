
```markdown
# Graphical Password Authentication

[![Go Report Card](https://goreportcard.com/badge/github.com/YashSaini99/graphical-password-authentication)](https://goreportcard.com/report/github.com/YashSaini99/graphical-password-authentication)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Graphical Password Authentication is a Go package that provides a secure and advanced authentication system using graphical passwords. Instead of relying on traditional text-based passwords, users provide a sequence of image selections (represented by indices), which are then converted into a string, hashed using bcrypt, and stored securely in MongoDB.

This package includes:
- **User Registration:** Register a new user with a username, email, and graphical password.
- **User Authentication:** Authenticate a user by verifying the hashed graphical password.
- **Advanced Security Features:** Brute-force protection (with account blocking and alert emails) and password reset functionality with secure token generation.
- **Email Validation:** Ensure that provided email addresses are in the correct format.
- **Facade Pattern:** A single import (via the `gp` package) re-exports all functionality for ease of use.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Basic Authentication](#basic-authentication)
  - [Advanced Security Features](#advanced-security-features)
- [API Reference](#api-reference)
- [Example](#example)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Graphical Password Registration & Authentication:**  
  Users register by selecting a pattern (a sequence of image indices) which is converted to a string and hashed.
  
- **Brute-Force Protection:**  
  Protects against repeated failed login attempts by blocking accounts after a threshold is reached and sending alert emails.

- **Password Reset:**  
  Generate secure reset tokens and send password reset emails to allow users to change their graphical password if forgotten.

- **Email Validation:**  
  Validate the format of user-provided email addresses.

- **Modular & Facade Design:**  
  The core authentication logic is decoupled from the UI. A facade package (`gp`) is provided so that developers can import everything with a single import.

## Installation

Make sure you have Go installed. Then, install the package using `go get`:

```bash
go get github.com/YashSaini99/graphical-password-authentication
```

Alternatively, if you are testing locally, use a replace directive in your project's `go.mod`:

```go
replace github.com/YashSaini99/graphical-password-authentication => /home/yash/Documents/graphical_password_authentication
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

## API Reference

### Core Functions (from `gp` facade)

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

This project is licensed under the MIT License.
```
