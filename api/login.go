// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"time"

	"github.com/clivern/mut/db"
	"github.com/clivern/mut/module"
	"github.com/clivern/mut/service"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Email      string `json:"email" validate:"required,email" label:"Email"`
	Password   string `json:"password" validate:"required" label:"Password"`
	RememberMe bool   `json:"rememberMe" validate:"omitempty,boolean" label:"Remember Me"`
}

// LoginAction handles login requests
func LoginAction(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Login endpoint called")

	var req LoginRequest
	if err := service.DecodeAndValidate(r, &req); err != nil {
		service.WriteValidationError(w, err)
		return
	}

	userRepo := db.NewUserRepository(db.GetDB())
	sessionRepo := db.NewSessionRepository(db.GetDB())
	authModule := module.NewAuth(userRepo)

	user, err := authModule.Login(req.Email, req.Password)
	if err != nil {
		service.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"errorMessage": "Invalid credentials",
		})
		return
	}

	if !user.IsActive {
		service.WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"errorMessage": "User is not active",
		})
		return
	}

	sessionManager := module.NewSessionManager(sessionRepo, userRepo)
	session, err := sessionManager.CreateSession(
		user.ID,
		time.Hour*24*7,
		r.RemoteAddr,
		r.UserAgent(),
	)
	if err != nil {
		service.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"errorMessage": "Failed to create session",
		})
		return
	}

	var cookieOptions *service.CookieOptions
	if viper.GetBool("app.tls.status") {
		cookieOptions = service.SecureCookieOptions()
	} else {
		cookieOptions = service.DefaultCookieOptions()
	}
	if req.RememberMe {
		cookieOptions.MaxAge = int((time.Hour * 24 * 30) / time.Second)
	} else {
		cookieOptions.MaxAge = 0
	}

	service.SetCookie(w, "_mut_session", session.Token, cookieOptions)
	service.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"successMessage": "Login successful",
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
