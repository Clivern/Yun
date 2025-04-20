// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Email    string `json:"email" validate:"required,email"`
	URL      string `json:"url" validate:"required,url"`
	Password string `json:"password" validate:"required,strong_password"`
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Age      int    `json:"age" validate:"omitempty,gte=0,lte=150"`
}

func TestUnitValidateStruct(t *testing.T) {
	t.Run("Valid struct", func(t *testing.T) {
		data := TestStruct{
			Email:    "test@example.com",
			URL:      "https://example.com",
			Password: "SecurePass123!",
			Name:     "John Doe",
			Age:      30,
		}

		err := ValidateStruct(data)
		assert.NoError(t, err)
	})

	t.Run("Invalid email", func(t *testing.T) {
		data := TestStruct{
			Email:    "invalid-email",
			URL:      "https://example.com",
			Password: "SecurePass123!",
			Name:     "John",
		}

		err := ValidateStruct(data)
		assert.Error(t, err)

		errorMsg := FormatValidationErrors(err)
		assert.NotEmpty(t, errorMsg)
		assert.Contains(t, errorMsg, "email")
	})

	t.Run("Invalid URL", func(t *testing.T) {
		data := TestStruct{
			Email:    "test@example.com",
			URL:      "not-a-url",
			Password: "SecurePass123!",
			Name:     "John",
		}

		err := ValidateStruct(data)
		assert.Error(t, err)

		errorMsg := FormatValidationErrors(err)
		assert.NotEmpty(t, errorMsg)
		assert.Contains(t, errorMsg, "URL")
	})

	t.Run("Weak password", func(t *testing.T) {
		data := TestStruct{
			Email:    "test@example.com",
			URL:      "https://example.com",
			Password: "weak",
			Name:     "John",
		}

		err := ValidateStruct(data)
		assert.Error(t, err)

		errorMsg := FormatValidationErrors(err)
		assert.NotEmpty(t, errorMsg)
		assert.Contains(t, errorMsg, "password")
	})

	t.Run("Missing required fields", func(t *testing.T) {
		data := TestStruct{}

		err := ValidateStruct(data)
		assert.Error(t, err)

		errorMsg := FormatValidationErrors(err)
		assert.NotEmpty(t, errorMsg)
		assert.True(t,
			assert.ObjectsAreEqual(errorMsg, errorMsg) &&
				(len(errorMsg) > 0),
			"Error message should not be empty")
	})

	t.Run("Name too short", func(t *testing.T) {
		data := TestStruct{
			Email:    "test@example.com",
			URL:      "https://example.com",
			Password: "SecurePass123!",
			Name:     "J",
		}

		err := ValidateStruct(data)
		assert.Error(t, err)

		errorMsg := FormatValidationErrors(err)
		assert.NotEmpty(t, errorMsg)
		assert.Contains(t, errorMsg, "name")
		assert.Contains(t, errorMsg, "at least 2")
	})

	t.Run("Age out of range", func(t *testing.T) {
		data := TestStruct{
			Email:    "test@example.com",
			URL:      "https://example.com",
			Password: "SecurePass123!",
			Name:     "John",
			Age:      200,
		}

		err := ValidateStruct(data)
		assert.Error(t, err)

		errorMsg := FormatValidationErrors(err)
		assert.NotEmpty(t, errorMsg)
		assert.Contains(t, errorMsg, "age")
	})
}

func TestUnitFormatValidationErrors(t *testing.T) {
	data := TestStruct{
		Email: "invalid",
		URL:   "invalid",
	}

	err := ValidateStruct(data)
	assert.Error(t, err)

	errorMsg := FormatValidationErrors(err)
	assert.NotEmpty(t, errorMsg)
	assert.Contains(t, errorMsg, "email")
}
