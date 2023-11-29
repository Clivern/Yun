// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/yun/db"
	"github.com/clivern/yun/module"
	"github.com/clivern/yun/service"

	"github.com/rs/zerolog/log"
)

// SetupRequest represents the setup request payload
type SetupRequest struct {
	GatewayURL    string `json:"gatewayURL" validate:"required,url,min=4,max=60" label:"Gateway URL"`
	GatewayEmail  string `json:"gatewayEmail" validate:"required,email,min=4,max=60" label:"Gateway Email"`
	GatewayName   string `json:"gatewayName" validate:"required,min=2,max=50" label:"Gateway Name"`
	AdminEmail    string `json:"adminEmail" validate:"required,email,min=4,max=60" label:"Admin Email"`
	AdminPassword string `json:"adminPassword" validate:"required,strong_password,min=8,max=60" label:"Admin Password"`
}

// SetupAction handles the setup installation
func SetupAction(w http.ResponseWriter, r *http.Request) {
	// Parse and validate request
	var req SetupRequest

	if err := service.DecodeAndValidate(r, &req); err != nil {
		service.WriteValidationError(w, err)
		return
	}

	// Create setup instance
	setupModule := module.NewSetup(
		db.NewOptionRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)

	// Check if already installed
	if setupModule.IsInstalled() {
		service.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errorMessage": "Gateway is already installed",
		})
		return
	}

	// Perform installation
	err := setupModule.Install(&module.SetupOptions{
		GatewayURL:    req.GatewayURL,
		GatewayEmail:  req.GatewayEmail,
		GatewayName:   req.GatewayName,
		AdminEmail:    req.AdminEmail,
		AdminPassword: req.AdminPassword,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to complete setup")
		service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"errorMessage": "Failed to complete setup",
		})
		return
	}

	log.Info().Msg("Gateway setup completed successfully")
	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"successMessage": "Gateway setup completed successfully",
	})
}

// SetupStatusAction checks if the gateway is already installed
func SetupStatusAction(w http.ResponseWriter, r *http.Request) {
	setupModule := module.NewSetup(
		db.NewOptionRepository(db.GetDB()),
		db.NewUserRepository(db.GetDB()),
	)
	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"installed": setupModule.IsInstalled(),
	})
}
