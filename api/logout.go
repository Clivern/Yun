// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/mut/db"
	"github.com/clivern/mut/middleware"
	"github.com/clivern/mut/module"
	"github.com/clivern/mut/service"

	"github.com/rs/zerolog/log"
)

// LogoutAction handles logout requests
func LogoutAction(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Logout endpoint called")

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		service.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"errorMessage": "Not authenticated",
		})
		return
	}

	sessionManager := module.NewSessionManager(
		db.NewSessionRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)

	sessionManager.CleanupExpiredSessions()

	if err := sessionManager.RevokeUserSessions(user.ID); err != nil {
		log.Error().Err(err).Msg("Failed to revoke session")
	}
	service.DeleteCookie(w, "_mut_session")

	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"successMessage": "Logout successful",
	})
}
