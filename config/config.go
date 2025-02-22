// Package config handles the loading of environment variables.
package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file.
// If the .env file is not found, it logs a warning and continues using system environment variables.
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using default environment variables.")
	}
	return err
}
