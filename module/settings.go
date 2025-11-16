// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/clivern/mut/db"
)

// Settings manages application configuration options.
type Settings struct {
	OptionRepository *db.OptionRepository
}

// NewSettings creates a new Settings module instance.
func NewSettings(optionRepository *db.OptionRepository) *Settings {
	return &Settings{OptionRepository: optionRepository}
}

// UpdateSettingsOptions contains the settings fields that can be updated.
type UpdateSettingsOptions struct {
	GatewayName        string
	GatewayURL         string
	GatewayEmail       string
	GatewayDescription string
	SMTPServer         string
	SMTPPort           int
	SMTPFromEmail      string
	SMTPUsername       string
	SMTPPassword       string
	SMTPUseTLS         bool
	MaintenanceMode    bool
}

// Update persists the provided settings to the options repository.
func (s *Settings) Update(options *UpdateSettingsOptions) error {
	updates := map[string]string{
		"gateway_name":        strings.TrimSpace(options.GatewayName),
		"gateway_url":         strings.TrimSpace(options.GatewayURL),
		"gateway_email":       strings.TrimSpace(options.GatewayEmail),
		"gateway_description": strings.TrimSpace(options.GatewayDescription),
		"smtp_server":         strings.TrimSpace(options.SMTPServer),
		"smtp_port":           strconv.Itoa(options.SMTPPort),
		"smtp_from_email":     strings.TrimSpace(options.SMTPFromEmail),
		"smtp_username":       strings.TrimSpace(options.SMTPUsername),
		"smtp_password":       options.SMTPPassword,
		"smtp_use_tls":        boolToOption(options.SMTPUseTLS),
		"maintenance_mode":    boolToOption(options.MaintenanceMode),
	}

	for key, value := range updates {
		if err := s.OptionRepository.Update(key, value); err != nil {
			return fmt.Errorf("failed to update option %s: %w", key, err)
		}
	}

	return nil
}

func boolToOption(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

// GetSettings retrieves the current settings from the options repository.
func (s *Settings) GetSettings() (map[string]string, error) {
	settings := make(map[string]string)

	gatewayName, err := s.OptionRepository.Get("gateway_name")
	if err != nil {
		return nil, err
	}
	settings["gateway_name"] = gatewayName.Value

	gatewayURL, err := s.OptionRepository.Get("gateway_url")
	if err != nil {
		return nil, err
	}
	settings["gateway_url"] = gatewayURL.Value

	gatewayEmail, err := s.OptionRepository.Get("gateway_email")
	if err != nil {
		return nil, err
	}
	settings["gateway_email"] = gatewayEmail.Value

	gatewayDescription, err := s.OptionRepository.Get("gateway_description")
	if err != nil {
		return nil, err
	}
	settings["gateway_description"] = gatewayDescription.Value

	smtpServer, err := s.OptionRepository.Get("smtp_server")
	if err != nil {
		return nil, err
	}
	settings["smtp_server"] = smtpServer.Value

	smtpPort, err := s.OptionRepository.Get("smtp_port")
	if err != nil {
		return nil, err
	}
	settings["smtp_port"] = smtpPort.Value

	smtpFromEmail, err := s.OptionRepository.Get("smtp_from_email")
	if err != nil {
		return nil, err
	}
	settings["smtp_from_email"] = smtpFromEmail.Value

	smtpUsername, err := s.OptionRepository.Get("smtp_username")
	if err != nil {
		return nil, err
	}
	settings["smtp_username"] = smtpUsername.Value

	smtpPassword, err := s.OptionRepository.Get("smtp_password")
	if err != nil {
		return nil, err
	}
	settings["smtp_password"] = smtpPassword.Value

	smtpUseTLS, err := s.OptionRepository.Get("smtp_use_tls")
	if err != nil {
		return nil, err
	}
	settings["smtp_use_tls"] = smtpUseTLS.Value

	maintenanceMode, err := s.OptionRepository.Get("maintenance_mode")
	if err != nil {
		return nil, err
	}
	settings["maintenance_mode"] = maintenanceMode.Value

	return settings, nil
}
