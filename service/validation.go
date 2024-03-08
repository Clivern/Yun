// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// Validator is a singleton instance of the validator
var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom tag name function to use 'label' tag
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		label := fld.Tag.Get("label")
		if label != "" {
			return label
		}
		// Fallback to json tag or field name
		name := fld.Tag.Get("json")
		if name != "" {
			return name
		}
		return fld.Name
	})

	// Register custom validators
	validate.RegisterValidation("strong_password", validateStrongPassword)
}

// GetValidator returns the global validator instance
func GetValidator() *validator.Validate {
	return validate
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag,omitempty"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// ValidateStruct validates a struct and returns formatted errors
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// FormatValidationErrors returns the first validation error message
func FormatValidationErrors(err error) string {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		if len(validationErrs) > 0 {
			// Return the first error message
			return getErrorMessage(validationErrs[0])
		}
	}
	// Return empty string if no errors
	return ""
}

// getErrorMessage returns a user-friendly error message based on the validation tag
func getErrorMessage(e validator.FieldError) string {
	// Use the field name as set by RegisterTagNameFunc (which uses 'label' tag if available)
	field := e.Field()

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, e.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, e.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, e.Param())
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "alpha":
		return fmt.Sprintf("%s must contain only letters", field)
	case "numeric":
		return fmt.Sprintf("%s must contain only numbers", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	case "containsany":
		return fmt.Sprintf("%s must contain at least one of: %s", field, e.Param())
	case "startswith":
		return fmt.Sprintf("%s must start with %s", field, e.Param())
	case "endswith":
		return fmt.Sprintf("%s must end with %s", field, e.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "strong_password":
		return fmt.Sprintf("%s must contain at least 8 characters, one uppercase, one lowercase, one digit, and one special character", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// DecodeJSON decodes JSON from request body
func DecodeJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("Invalid JSON format: %w", err)
	}
	return nil
}

// DecodeAndValidate decodes JSON and validates the struct in one step
func DecodeAndValidate(r *http.Request, v interface{}) error {
	// Decode JSON
	if err := DecodeJSON(r, v); err != nil {
		return err
	}

	// Validate struct
	return ValidateStruct(v)
}

// WriteValidationError writes validation errors as JSON response
func WriteValidationError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		w.WriteHeader(http.StatusBadRequest)
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": FormatValidationErrors(validationErrs),
		})
	} else {
		WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": err.Error(),
		})
	}
}
