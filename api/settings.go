// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/clivern/yun/service"

	"github.com/rs/zerolog/log"
)

// UpdateSettingsAction handles user settings update requests
func UpdateSettingsAction(w http.ResponseWriter, _ *http.Request) {
	log.Debug().Msg("Update settings endpoint called")

	service.WriteJSON(w, http.StatusOK, map[string]interface{}{})
}
