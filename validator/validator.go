package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func validateString(value string, min int, max int) error {

	n := len(value)

	if n < min || n > max {
		return fmt.Errorf("length must be between %d and %d", min, max)
	}

	return nil
}

func ValidateUsername(username string) error {

	if err := validateString(username, 3, 50); err != nil {
		return err
	}

	if !isValidUsername(username) {
		return fmt.Errorf("username must be alphanumeric")
	}

	return nil
}

func ValidatePassword(password string) error {

	if err := validateString(password, 6, 100); err != nil {
		return err
	}

	return nil
}

func ValidateEmail(value string) error {

	if err := validateString(value, 3, 50); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

func ValidateFullName(value string) error {

	if err := validateString(value, 3, 50); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return fmt.Errorf("full name must be alphabetic")
	}

	return nil
}
