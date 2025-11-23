// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/clivern/mut/db"
	"github.com/clivern/mut/middleware"
	"github.com/clivern/mut/module"
	"github.com/clivern/mut/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// CreateUserRequest represents the create user request payload
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email,min=4,max=60" label:"Email"`
	Password string `json:"password" validate:"required,strong_password,min=8,max=60" label:"Password"`
	Role     string `json:"role" validate:"required,oneof=admin user readonly" label:"Role"`
	IsActive bool   `json:"isActive" label:"Is Active"`
}

// UpdateUserRequest represents the update user request payload
type UpdateUserRequest struct {
	Email    string `json:"email" validate:"required,email,min=4,max=60" label:"Email"`
	Password string `json:"password" validate:"omitempty,strong_password,min=8,max=60" label:"Password"`
	Role     string `json:"role" validate:"required,oneof=admin user readonly" label:"Role"`
	IsActive bool   `json:"isActive" label:"Is Active"`
}

// CreateUserAction handles user creation requests
func CreateUserAction(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Create user endpoint called")

	var req CreateUserRequest
	if err := service.DecodeAndValidate(r, &req); err != nil {
		service.WriteValidationError(w, err)
		return
	}

	userModule := module.NewUser(db.NewUserRepository(db.GetDB()))
	user, err := userModule.CreateUser(&module.CreateUserOptions{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		IsActive: req.IsActive,
	})

	if err != nil {
		if errors.Is(err, module.ErrUserEmailAlreadyExists) {
			service.WriteJSON(w, http.StatusConflict, map[string]interface{}{
				"errorMessage": "User with this email already exists",
			})
			return
		}
		log.Error().Err(err).Msg("Failed to create user")
		service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"errorMessage": "Failed to create user",
		})
		return
	}

	log.Info().Int64("userID", user.ID).Msg("User created successfully")
	service.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"id":          user.ID,
		"email":       user.Email,
		"role":        user.Role,
		"isActive":    user.IsActive,
		"apiKey":      user.APIKey,
		"lastLoginAt": user.LastLoginAt.UTC().Format(time.RFC3339),
		"createdAt":   user.CreatedAt.UTC().Format(time.RFC3339),
		"updatedAt":   user.UpdatedAt.UTC().Format(time.RFC3339),
	})
}

// GetUserAction handles get user by ID requests
func GetUserAction(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Get user endpoint called")

	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		service.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": "Invalid user ID",
		})
		return
	}

	userModule := module.NewUser(db.NewUserRepository(db.GetDB()))
	user, err := userModule.GetUser(userID)
	if err != nil {
		if errors.Is(err, module.ErrUserNotFound) {
			service.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": "User not found",
			})
			return
		}
		log.Error().Err(err).Msg("Failed to get user")
		service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"errorMessage": "Failed to get user",
		})
		return
	}

	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":          user.ID,
		"email":       user.Email,
		"role":        user.Role,
		"isActive":    user.IsActive,
		"apiKey":      user.APIKey,
		"lastLoginAt": user.LastLoginAt.UTC().Format(time.RFC3339),
		"createdAt":   user.CreatedAt.UTC().Format(time.RFC3339),
		"updatedAt":   user.UpdatedAt.UTC().Format(time.RFC3339),
	})
}

// UpdateUserAction handles user update requests
func UpdateUserAction(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Update user endpoint called")

	// Get user ID from URL
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		service.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": "Invalid user ID",
		})
		return
	}

	var req UpdateUserRequest
	if err := service.DecodeAndValidate(r, &req); err != nil {
		service.WriteValidationError(w, err)
		return
	}

	userModule := module.NewUser(db.NewUserRepository(db.GetDB()))
	user, err := userModule.UpdateUser(&module.UpdateUserOptions{
		UserID:   userID,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		IsActive: req.IsActive,
	})

	if err != nil {
		if errors.Is(err, module.ErrUserNotFound) {
			service.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": "User not found",
			})
			return
		}
		if errors.Is(err, module.ErrUserEmailAlreadyExists) {
			service.WriteJSON(w, http.StatusConflict, map[string]interface{}{
				"errorMessage": "User with this email already exists",
			})
			return
		}
		log.Error().Err(err).Msg("Failed to update user")
		service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"errorMessage": "Failed to update user",
		})
		return
	}

	log.Info().Int64("userID", user.ID).Msg("User updated successfully")
	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":          user.ID,
		"email":       user.Email,
		"role":        user.Role,
		"isActive":    user.IsActive,
		"apiKey":      user.APIKey,
		"lastLoginAt": user.LastLoginAt.UTC().Format(time.RFC3339),
		"createdAt":   user.CreatedAt.UTC().Format(time.RFC3339),
		"updatedAt":   user.UpdatedAt.UTC().Format(time.RFC3339),
	})
}

// ListUsersAction handles user listing requests with pagination
func ListUsersAction(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("List users endpoint called")

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	userModule := module.NewUser(db.NewUserRepository(db.GetDB()))
	result, err := userModule.ListUsers(&module.ListUsersOptions{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to list users")
		service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"errorMessage": "Failed to list users",
		})
		return
	}

	userList := make([]map[string]interface{}, 0, len(result.Users))
	for _, user := range result.Users {
		userList = append(userList, map[string]interface{}{
			"id":          user.ID,
			"email":       user.Email,
			"role":        user.Role,
			"isActive":    user.IsActive,
			"apiKey":      user.APIKey,
			"lastLoginAt": user.LastLoginAt.UTC().Format(time.RFC3339),
			"createdAt":   user.CreatedAt.UTC().Format(time.RFC3339),
			"updatedAt":   user.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"users": userList,
		"_meta": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"total":  result.Total,
		},
	})
}

// DeleteUserAction handles user deletion requests
func DeleteUserAction(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Delete user endpoint called")

	currentUser, _ := middleware.GetUserFromContext(r.Context())

	// Get user ID from URL
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		service.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": "Invalid user ID",
		})
		return
	}

	// Prevent self-deletion
	if currentUser.ID == userID {
		service.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": "You cannot delete your own account",
		})
		return
	}

	userModule := module.NewUser(db.NewUserRepository(db.GetDB()))
	err = userModule.DeleteUser(userID)
	if err != nil {
		if errors.Is(err, module.ErrUserNotFound) {
			service.WriteJSON(w, http.StatusNotFound, map[string]interface{}{
				"errorMessage": "User not found",
			})
			return
		}
		log.Error().Err(err).Msg("Failed to delete user")
		service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"errorMessage": "Failed to delete user",
		})
		return
	}

	log.Info().Int64("userID", userID).Msg("User deleted successfully")
	service.WriteJSON(w, http.StatusNoContent, map[string]interface{}{})
}
