package user

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint      `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Address      string    `json:"addres,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Email        string    `json:"email,omitempty"`
	Password     string    `json:"password,omitempty"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

func (u *User) HashPassword() error {
	// Check if the password meets the requirements
	if !validatePassword(u.Password) {
		missingRequirements := getMissingRequirements(u.Password)
		errorMessage := "The password does not meet the minimum requirements. The following requirements are missing: " + strings.Join(missingRequirements, ", ")
		return errors.New(errorMessage)
	}

	// Generate the password hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(passwordHash)
	return nil
}

func (u *User) PasswordMatch(password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func validatePassword(password string) bool {
	// Regular expression to validate the password
	regex := regexp.MustCompile(`[A-Z]|\d|[@#$%^&+=!]`)

	// Check if the password matches the regular expression
	return regex.MatchString(password)
}

func getMissingRequirements(password string) []string {
	missingRequirements := []string{}

	// Check if the password has at least 6 characters
	if len(password) < 6 {
		missingRequirements = append(missingRequirements, "minimum of 6 characters")
	}

	// Check if the password has at least one uppercase letter
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		missingRequirements = append(missingRequirements, "at least one uppercase letter")
	}

	return missingRequirements
}
