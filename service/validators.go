// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// validateStrongPassword validates that a password meets security requirements
// Requires: min 8 chars, 1 uppercase, 1 lowercase, 1 digit, 1 special character
//
// Usage: Password string `validate:"strong_password"`
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}
	// Has uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// Has lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// Has digit
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	// Has special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`).MatchString(password)

	return hasUpper && hasLower && hasDigit && hasSpecial
}
