// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package middleware provides HTTP middleware functions for the Mut application.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/clivern/mut/db"
	"github.com/clivern/mut/module"
	"github.com/clivern/mut/service"

	"github.com/rs/zerolog/log"
)

// Context keys for storing user and session data
const (
	// ContextKeyUser is the key for storing user in context
	ContextKeyUser contextKey = "user"
)

// SessionAuth creates a session-based authentication middleware
// It validates the session cookie and stores user/session in the request context
func SessionAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication for specific routes
			if shouldSkipAuth(r.URL.Path) {
				log.Info().Str("path", r.URL.Path).Msg("Skipping authentication for API route")
				next.ServeHTTP(w, r)
				return
			}

			// Check if API key is present in the request header "X-API-Key"
			apiKey := r.Header.Get("X-API-Key")
			if apiKey != "" {
				user, err := db.NewUserRepository(db.GetDB()).GetByAPIKey(apiKey)
				if err != nil {
					log.Info().Err(err).Str("path", r.URL.Path).Msg("API key validation failed")
					service.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
						"errorMessage": "Invalid API key",
					})
					return
				}
				log.Info().Str("path", r.URL.Path).Msg("API key validation successful")
				// Store user in context
				ctx := context.WithValue(r.Context(), ContextKeyUser, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Get session token from cookie
			sessionToken := service.GetCookie(r, "_mut_session")
			if sessionToken == "" {
				log.Info().Str("path", r.URL.Path).Msg("No session cookie found")
				service.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"errorMessage": "Not authenticated",
				})
				return
			}

			// Validate session
			sessionManager := module.NewSessionManager(
				db.NewSessionRepository(db.GetDB()),
				db.NewUserRepository(db.GetDB()),
			)

			user, _, err := sessionManager.ValidateSession(sessionToken)
			if err != nil {
				log.Info().Err(err).Str("path", r.URL.Path).Msg("Session validation failed")
				service.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"errorMessage": "Invalid or expired session",
				})
				return
			}

			log.Info().Str("path", r.URL.Path).Msg("Session validation successful")
			// Store user and session in context
			ctx := context.WithValue(r.Context(), ContextKeyUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// shouldSkipAuth determines if authentication should be skipped for a given path
func shouldSkipAuth(path string) bool {
	// Skip auth for public API routes
	// Skip auth for non-API routes (static files, frontend, etc.)
	if strings.HasPrefix(path, "/api/v1/public/") || !strings.HasPrefix(path, "/api/v1/") {
		return true
	}

	return false
}

// GetUserFromContext retrieves the user from the request context
func GetUserFromContext(ctx context.Context) (*db.User, bool) {
	user, ok := ctx.Value(ContextKeyUser).(*db.User)
	return user, ok
}
