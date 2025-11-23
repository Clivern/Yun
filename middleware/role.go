// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package middleware provides HTTP middleware functions for the Tut application.
package middleware

import (
	"fmt"
	"net/http"

	"github.com/clivern/mut/service"

	"github.com/rs/zerolog/log"
)

// RequireRole creates a middleware that checks if the authenticated user has one of the required roles
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, ok := GetUserFromContext(r.Context())
			if !ok || user == nil {
				log.Info().Str("path", r.URL.Path).Msg("User not found in context for role check")
				service.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"errorMessage": "Not authenticated",
				})
				return
			}

			// Check if user is active
			if !user.IsActive {
				log.Info().
					Str("path", r.URL.Path).
					Int64("userID", user.ID).
					Msg("Inactive user attempted to access protected route")
				service.WriteJSON(w, http.StatusForbidden, map[string]interface{}{
					"errorMessage": "Account is inactive",
				})
				return
			}

			// Check if user has one of the allowed roles
			hasRole := false
			fmt.Println("Allowed Roles", allowedRoles)
			for _, role := range allowedRoles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				log.Info().
					Str("path", r.URL.Path).
					Int64("userID", user.ID).
					Str("userRole", user.Role).
					Strs("allowedRoles", allowedRoles).
					Msg("User does not have required role")
				service.WriteJSON(w, http.StatusForbidden, map[string]interface{}{
					"errorMessage": "Insufficient permissions",
				})
				return
			}

			// User has required role, proceed
			next.ServeHTTP(w, r)
		})
	}
}
