// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/clivern/mut/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegrationReadyEndpoint tests the readiness check endpoint
func TestIntegrationReadyEndpoint(t *testing.T) {
	t.Run("ReadyAction should return OK when database is available", func(t *testing.T) {
		// Clean up any existing global connection
		db.CloseDB()

		// Create a temporary SQLite database
		tmpFile := "/tmp/test_ready_" + strings.ReplaceAll(t.Name(), "/", "_") + ".db"
		defer os.Remove(tmpFile)

		config := db.Config{
			Driver:     "sqlite",
			DataSource: tmpFile,
		}

		// Initialize the database
		err := db.InitDB(config)
		require.NoError(t, err, "Failed to initialize test database")
		defer db.CloseDB()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/_ready", nil)
		w := httptest.NewRecorder()

		ReadyAction(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, `{"status":"ok"}`, strings.TrimSpace(w.Body.String()))
	})

	t.Run("ReadyAction should return ServiceUnavailable when database ping fails", func(t *testing.T) {
		// Clean up any existing global connection
		db.CloseDB()

		// Create a temporary SQLite database
		tmpFile := "/tmp/test_ready_fail_" + strings.ReplaceAll(t.Name(), "/", "_") + ".db"
		defer os.Remove(tmpFile)

		config := db.Config{
			Driver:     "sqlite",
			DataSource: tmpFile,
		}

		// Initialize the database
		err := db.InitDB(config)
		require.NoError(t, err, "Failed to initialize test database")

		// Close the underlying database connection to make Ping() fail
		// We need to access the connection and close it directly
		database := db.GetDB()
		require.NotNil(t, database, "Database connection should be initialized")
		database.Close()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/_ready", nil)
		w := httptest.NewRecorder()

		ReadyAction(w, req)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.Equal(t, `{"status":"not_ok"}`, strings.TrimSpace(w.Body.String()))

		// Clean up
		db.CloseDB()
	})
}
