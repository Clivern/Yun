// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStrongPassword(t *testing.T) {
	// Test struct for password validation
	type PasswordTest struct {
		Password string `validate:"strong_password"`
	}

	testCases := []struct {
		name     string
		password string
		expected bool
	}{
		// Valid passwords
		{
			name:     "Valid password with all requirements",
			password: "Password123!",
			expected: true,
		},
		{
			name:     "Valid password with multiple special chars",
			password: "P@ssw0rd!#$",
			expected: true,
		},
		{
			name:     "Valid password exactly 8 characters",
			password: "Pass123!",
			expected: true,
		},
		{
			name:     "Valid long password",
			password: "MyVerySecureP@ssw0rd2025",
			expected: true,
		},
		{
			name:     "Valid password with brackets",
			password: "Test[Pass]123",
			expected: true,
		},
		{
			name:     "Valid password with various special chars",
			password: "Abc123!@#$%^&*()",
			expected: true,
		},
		// Invalid passwords - missing requirements
		{
			name:     "Missing uppercase letter",
			password: "password123!",
			expected: false,
		},
		{
			name:     "Missing lowercase letter",
			password: "PASSWORD123!",
			expected: false,
		},
		{
			name:     "Missing digit",
			password: "Password!@#",
			expected: false,
		},
		{
			name:     "Missing special character",
			password: "Password123",
			expected: false,
		},
		{
			name:     "Too short (7 characters)",
			password: "Pass12!",
			expected: false,
		},
		{
			name:     "Too short (6 characters) with all types",
			password: "Pas1!a",
			expected: false,
		},
		{
			name:     "Only lowercase letters",
			password: "password",
			expected: false,
		},
		{
			name:     "Only uppercase letters",
			password: "PASSWORD",
			expected: false,
		},
		{
			name:     "Only digits",
			password: "12345678",
			expected: false,
		},
		{
			name:     "Only special characters",
			password: "!@#$%^&*",
			expected: false,
		},
		{
			name:     "Empty string",
			password: "",
			expected: false,
		},
		{
			name:     "Single character",
			password: "A",
			expected: false,
		},
		{
			name:     "Missing two requirements (no upper, no special)",
			password: "password123",
			expected: false,
		},
		{
			name:     "Missing two requirements (no digit, no special)",
			password: "PasswordTest",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test struct with password
			testData := PasswordTest{Password: tc.password}

			// Validate using the validator framework
			err := ValidateStruct(testData)

			if tc.expected {
				// Should be valid (no error)
				assert.NoError(t, err, "Password %q should be valid", tc.password)
			} else {
				// Should be invalid (error expected)
				assert.Error(t, err, "Password %q should be invalid", tc.password)
			}
		})
	}
}
