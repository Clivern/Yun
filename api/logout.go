// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/mut/db"
	"github.com/clivern/mut/module"
	"github.com/clivern/mut/service"

	"github.com/rs/zerolog/log"
)

// LogoutAction handles logout requests
func LogoutAction(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Logout endpoint called")

	// Get session token from cookie
	sessionToken := service.GetCookie(r, "_mut_session")
	if sessionToken == "" {
		service.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"errorMessage": "Not authenticated",
		})
		return
	}

	// Revoke the session
	sessionManager := module.NewSessionManager(
		db.NewSessionRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)

	if err := sessionManager.RevokeSession(sessionToken); err != nil {
		log.Error().Err(err).Msg("Failed to revoke session")
	}

	// Clear the session cookie
	service.DeleteCookie(w, "_mut_session")

	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"successMessage": "Logout successful",
	})
}
