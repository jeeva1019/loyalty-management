package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidateEmail checks if the email is valid and returns an error if not.
func ValidateEmail(email string) error {
	const emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidatePassword checks for strong password requirements.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`)
		hasLower   = regexp.MustCompile(`[a-z]`)
		hasNumber  = regexp.MustCompile(`[0-9]`)
		hasSpecial = regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?]`)
	)

	switch {
	case !hasUpper.MatchString(password):
		return errors.New("password must contain at least one uppercase letter")
	case !hasLower.MatchString(password):
		return errors.New("password must contain at least one lowercase letter")
	case !hasNumber.MatchString(password):
		return errors.New("password must contain at least one digit")
	case !hasSpecial.MatchString(password):
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func Validator(fields map[string]string) error {
	for name, val := range fields {
		if strings.TrimSpace(val) == "" {
			return fmt.Errorf("field '%s' is required", name)
		}
	}
	return nil
}
