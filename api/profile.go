// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"time"

	"github.com/clivern/mut/middleware"
	"github.com/clivern/mut/service"

	"github.com/rs/zerolog/log"
)

// GetProfileAction handles user profile requests
func GetProfileAction(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Get profile endpoint called")

	// Get user from context (set by auth middleware)
	user, ok := middleware.GetUserFromContext(r.Context())

	if !ok {
		service.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"errorMessage": "Not authenticated",
		})
		return
	}

	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"id":          user.ID,
			"email":       user.Email,
			"role":        user.Role,
			"isActive":    user.IsActive,
			"lastLoginAt": user.LastLoginAt.UTC().Format(time.RFC3339),
			"createdAt":   user.CreatedAt.UTC().Format(time.RFC3339),
			"updatedAt":   user.UpdatedAt.UTC().Format(time.RFC3339),
		},
	})
}

// UpdateProfileAction handles user profile update requests
func UpdateProfileAction(w http.ResponseWriter, _ *http.Request) {
	log.Debug().Msg("Update profile endpoint called")

	service.WriteJSON(w, http.StatusOK, map[string]interface{}{})
}
