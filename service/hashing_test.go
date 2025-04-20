// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitHashPassword(t *testing.T) {
	t.Run("HashPassword with valid password", func(t *testing.T) {
		password := "mySecurePassword123"

		hashedPassword, err := HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
		assert.NotEqual(t, password, hashedPassword, "Hashed password should not equal plain text password")
		assert.True(t, len(hashedPassword) > 0, "Hashed password should have length")
	})

	t.Run("HashPassword generates different hashes for same password", func(t *testing.T) {
		password := "samePassword"

		hash1, err1 := HashPassword(password)
		hash2, err2 := HashPassword(password)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2, "Each hash should be unique due to salt")
	})

	t.Run("HashPassword with empty string", func(t *testing.T) {
		password := ""
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("HashPassword with long password", func(t *testing.T) {
		// Password is 70 bytes (under bcrypt's 72 byte limit)
		password := "thisIsAVeryLongPasswordWithCharacters1234567890!@#$%^&*()_+-=[]"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})

	t.Run("HashPassword with special characters", func(t *testing.T) {
		password := "p@ssw0rd!#$%^&*()"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashedPassword)
	})
}

func TestUnitComparePassword(t *testing.T) {
	t.Run("ComparePassword with matching password", func(t *testing.T) {
		password := "mySecurePassword123"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		result := ComparePassword(hashedPassword, password)
		assert.True(t, result, "Password should match the hash")
	})

	t.Run("ComparePassword with non-matching password", func(t *testing.T) {
		password := "correctPassword"
		wrongPassword := "wrongPassword"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		result := ComparePassword(hashedPassword, wrongPassword)
		assert.False(t, result, "Wrong password should not match the hash")
	})

	t.Run("ComparePassword with empty password", func(t *testing.T) {
		password := "somePassword"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		result := ComparePassword(hashedPassword, "")
		assert.False(t, result, "Empty password should not match non-empty hash")
	})

	t.Run("ComparePassword with invalid hash", func(t *testing.T) {
		invalidHash := "notAValidBcryptHash"
		password := "somePassword"
		result := ComparePassword(invalidHash, password)
		assert.False(t, result, "Invalid hash should return false")
	})

	t.Run("ComparePassword case sensitivity", func(t *testing.T) {
		password := "Password123"
		differentCase := "password123"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		result := ComparePassword(hashedPassword, differentCase)
		assert.False(t, result, "Password comparison should be case-sensitive")
	})

	t.Run("ComparePassword with special characters", func(t *testing.T) {
		password := "p@ssw0rd!#$%"
		hashedPassword, err := HashPassword(password)
		assert.NoError(t, err)
		result := ComparePassword(hashedPassword, password)
		assert.True(t, result, "Password with special characters should match")
	})
}

func TestUnitHashAndComparePasswordIntegration(t *testing.T) {
	t.Run("Complete password workflow", func(t *testing.T) {
		passwords := []string{
			"simplePassword",
			"Complex123!@#",
			"",
			"with spaces in password",
			"unicode密码🔐",
		}

		for _, password := range passwords {
			t.Run("Password: "+password, func(t *testing.T) {
				hashed, err := HashPassword(password)
				assert.NoError(t, err)
				resultCorrect := ComparePassword(hashed, password)
				assert.True(t, resultCorrect, "Original password should match")
				resultWrong := ComparePassword(hashed, password+"wrong")
				assert.False(t, resultWrong, "Wrong password should not match")
			})
		}
	})
}
