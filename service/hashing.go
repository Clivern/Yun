// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash from a plain text password.
// It uses the default cost factor (bcrypt.DefaultCost = 10).
//
// Example:
//
//	hashedPassword, err := HashPassword("mySecurePassword123")
//	if err != nil {
//		return err
//	}
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword compares a bcrypt hashed password with a plain text password.
// Returns true if the password matches the hash, false otherwise.
//
// Example:
//
//	isValid := ComparePassword(hashedPassword, "mySecurePassword123")
//	if !isValid {
//		return errors.New("invalid password")
//	}
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
