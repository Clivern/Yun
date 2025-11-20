// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"errors"
	"fmt"
	"time"

	"github.com/clivern/mut/db"
	"github.com/clivern/mut/service"
	"github.com/google/uuid"
)

// Setup handles the initial installation and configuration of the gateway.
type Setup struct {
	OptionRepository *db.OptionRepository
	UserRepository   *db.UserRepository
}

// SetupOptions contains the configuration options for gateway setup.
type SetupOptions struct {
	GatewayURL    string
	GatewayEmail  string
	GatewayName   string
	AdminEmail    string
	AdminPassword string
}

// NewSetup creates a new Setup instance with the provided repositories.
func NewSetup(optionRepository *db.OptionRepository, userRepository *db.UserRepository) *Setup {
	return &Setup{OptionRepository: optionRepository, UserRepository: userRepository}
}

// IsInstalled checks whether the gateway has been installed.
func (s *Setup) IsInstalled() bool {
	option, err := s.OptionRepository.Get("is_installed")
	if err != nil {
		return false
	}
	return option != nil
}

// Install performs the initial gateway installation with the provided options.
func (s *Setup) Install(options *SetupOptions) error {
	if s.IsInstalled() {
		return errors.New("gateway is already installed")
	}

	hashedPassword, err := service.HashPassword(options.AdminPassword)
	if err != nil {
		return err
	}
	user := &db.User{
		Email:       options.AdminEmail,
		Password:    hashedPassword,
		Role:        db.UserRoleAdmin,
		APIKey:      uuid.New().String(),
		IsActive:    true,
		LastLoginAt: time.Now().UTC(),
	}
	fmt.Println("user", user)
	err = s.UserRepository.Create(user)
	if err != nil {
		return err
	}

	err = s.OptionRepository.Create("is_installed", "1")
	if err != nil {
		return err
	}

	err = s.OptionRepository.Create("gateway_url", options.GatewayURL)
	if err != nil {
		return err
	}

	err = s.OptionRepository.Create("gateway_email", options.GatewayEmail)
	if err != nil {
		return err
	}

	err = s.OptionRepository.Create("gateway_name", options.GatewayName)
	if err != nil {
		return err
	}

	err = s.OptionRepository.Create("maintenance_mode", "0")
	if err != nil {
		return err
	}

	err = s.OptionRepository.Create("gateway_description", "")
	if err != nil {
		return err
	}

	err = s.OptionRepository.Create("smtp_server", "")
	if err != nil {
		return err
	}

	err = s.OptionRepository.Create("smtp_port", "587")
	if err != nil {
		return err
	}

	err = s.OptionRepository.Create("smtp_from_email", "")
	if err != nil {
		return err
	}

	err = s.OptionRepository.Create("smtp_username", "")
	if err != nil {
		return err
	}

	err = s.OptionRepository.Create("smtp_password", "")
	if err != nil {
		return err
	}

	err = s.OptionRepository.Create("smtp_use_tls", "0")
	if err != nil {
		return err
	}

	return nil
}
