// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/mut/db"
	"github.com/clivern/mut/service"

	"github.com/rs/zerolog/log"
)

// ReadyAction handles readiness check requests
func ReadyAction(w http.ResponseWriter, _ *http.Request) {
	log.Debug().Msg("Readiness check endpoint called")

	// Try to ping the database
	database := db.GetDB()
	if err := database.Ping(); err != nil {
		log.Error().Err(err).Msg("Database ping failed during readiness check")

		service.WriteJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status": "not_ok",
		})
		return
	}

	log.Debug().Msg("Readiness check passed")
	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
