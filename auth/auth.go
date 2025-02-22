// Package auth provides core authentication functionality using graphical passwords.
// Instead of a traditional text-based password, users provide a sequence of image indices.
// This sequence is converted to a string, hashed with bcrypt, and stored in MongoDB.
package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/YashSaini99/graphical-password-authentication/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserExists is returned when a user with the same username or email already exists.
	ErrUserExists = errors.New("username or email already exists")
	// ErrUserNotFound is returned when no user is found for a given username.
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidGraphicalPassword is returned when the provided graphical password does not match the stored hash.
	ErrInvalidGraphicalPassword = errors.New("invalid graphical password")
)

// User represents a user in the system.
type User struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty"`
	Username              string             `bson:"username"`
	Email                 string             `bson:"email"`
	GraphicalPasswordHash string             `bson:"graphical_password_hash"`
}

// convertGraphicalPasswordToString converts a slice of integers (representing image indices)
// into a dash-separated string. For example, []int{1, 3, 5} becomes "1-3-5".
func convertGraphicalPasswordToString(gp []int) string {
	strs := make([]string, len(gp))
	for i, num := range gp {
		strs[i] = strconv.Itoa(num)
	}
	return strings.Join(strs, "-")
}

// HashGraphicalPassword converts the provided graphical password (slice of ints)
// into a string and hashes it using bcrypt.
func HashGraphicalPassword(gp []int) (string, error) {
	passStr := convertGraphicalPasswordToString(gp)
	hashed, err := bcrypt.GenerateFromPassword([]byte(passStr), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash graphical password: %w", err)
	}
	return string(hashed), nil
}

// CheckGraphicalPasswordHash compares the provided graphical password (as a slice of ints)
// with the stored hash. It returns nil if the password matches, or an error otherwise.
func CheckGraphicalPasswordHash(gp []int, hash string) error {
	passStr := convertGraphicalPasswordToString(gp)
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passStr)); err != nil {
		return fmt.Errorf("graphical password comparison failed: %w", err)
	}
	return nil
}

// RegisterUser registers a new user using their username, email, and a graphical password.
// The graphicalPassword parameter is the sequence of image indices chosen by the user.
// It returns an error if the user already exists or if registration fails.
func RegisterUser(username, email string, graphicalPassword []int) error {
	collection := db.GetCollection("auth", "users")
	if collection == nil {
		return errors.New("failed to get users collection")
	}

	// Check if the username or email already exists.
	filter := bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": email},
		},
	}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to count existing users: %w", err)
	}
	if count > 0 {
		return ErrUserExists
	}

	// Hash the graphical password.
	hashedGP, err := HashGraphicalPassword(graphicalPassword)
	if err != nil {
		return err
	}

	// Create and insert the new user.
	user := User{
		Username:              username,
		Email:                 email,
		GraphicalPasswordHash: hashedGP,
	}
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return fmt.Errorf("failed to insert new user: %w", err)
	}
	return nil
}

// AuthenticateUser authenticates a user using their username and graphical password.
// It returns true if authentication is successful, false otherwise, and an error if any issues occur.
func AuthenticateUser(username string, graphicalPassword []int) (bool, error) {
	collection := db.GetCollection("auth", "users")
	if collection == nil {
		return false, errors.New("failed to get users collection")
	}

	// Retrieve the user by username.
	var user User
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, ErrUserNotFound
		}
		return false, fmt.Errorf("failed to retrieve user: %w", err)
	}

	// Compare the provided graphical password with the stored hash.
	if err = CheckGraphicalPasswordHash(graphicalPassword, user.GraphicalPasswordHash); err != nil {
		return false, ErrInvalidGraphicalPassword
	}
	return true, nil
}
