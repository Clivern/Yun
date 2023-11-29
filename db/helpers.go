// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	// globalConnection holds the singleton database connection
	globalConnection *Connection
	// mu protects globalConnection during initialization
	mu sync.RWMutex
)

// InitDB initializes the global database connection
// This should be called once at application startup
func InitDB(config Config) error {
	mu.Lock()
	defer mu.Unlock()

	if globalConnection != nil {
		log.Warn().Msg("Database connection already initialized")
		return nil
	}

	conn, err := NewConnection(config)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	globalConnection = conn
	log.Info().Msg("Global database connection initialized")
	return nil
}

// GetDB returns the global database connection
// Panics if the database has not been initialized
func GetDB() *sql.DB {
	mu.RLock()
	defer mu.RUnlock()

	if globalConnection == nil {
		log.Error().Msg("Database not initialized")
		return nil
	}

	return globalConnection.DB
}

// CloseDB closes the global database connection
// This should be called at application shutdown
func CloseDB() error {
	mu.Lock()
	defer mu.Unlock()

	if globalConnection == nil {
		return nil
	}

	err := globalConnection.Close()
	globalConnection = nil
	return err
}
