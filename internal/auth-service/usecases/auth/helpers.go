package auth

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	usernameRegexp      = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	passwordLetterRegex = regexp.MustCompile(`[A-Za-z]`)
	passwordDigitRegex  = regexp.MustCompile(`[0-9]`)
)

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("is required")
	}

	length := utf8.RuneCountInString(username)
	if length < 3 || length > 32 {
		return errors.New("must be between 3 and 32 characters")
	}

	if !usernameRegexp.MatchString(username) {
		return errors.New("can only contain letters, numbers, underscores, dots, and dashes")
	}

	return nil
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("is required")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("has invalid format")
	}

	return nil
}

func ValidatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return errors.New("is required")
	}

	length := utf8.RuneCountInString(password)
	if length < 8 || length > 128 {
		return errors.New("must be between 8 and 128 characters")
	}

	if !passwordLetterRegex.MatchString(password) {
		return errors.New("must contain at least one letter")
	}

	if !passwordDigitRegex.MatchString(password) {
		return errors.New("must contain at least one digit")
	}

	return nil
}
