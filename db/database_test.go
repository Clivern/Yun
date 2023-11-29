// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLiteConnection(t *testing.T) {
	// Create a temporary SQLite database
	tmpFile := "/tmp/test_yun.db"
	defer os.Remove(tmpFile)

	config := Config{
		Driver:     "sqlite",
		DataSource: tmpFile,
	}

	conn, err := NewConnection(config)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	// Test ping
	err = conn.Ping()
	assert.NoError(t, err)

	// Test close
	err = conn.Close()
	assert.NoError(t, err)
}

func TestUnsupportedDriver(t *testing.T) {
	config := Config{
		Driver: "postgresql",
	}

	conn, err := NewConnection(config)
	assert.Error(t, err)
	assert.Nil(t, conn)
	assert.Contains(t, err.Error(), "unsupported database driver")
}

func TestInitDBAndGetDB(t *testing.T) {
	// Clean up any existing global connection
	CloseDB()

	// Create a temporary SQLite database
	tmpFile := "/tmp/test_yun_global.db"
	defer os.Remove(tmpFile)

	config := Config{
		Driver:     "sqlite",
		DataSource: tmpFile,
	}

	// Test InitDB
	err := InitDB(config)
	assert.NoError(t, err)

	// Test GetDB
	db := GetDB()
	assert.NotNil(t, db)

	// Test that the connection is working
	err = db.Ping()
	assert.NoError(t, err)

	// Test InitDB called twice (should not error)
	err = InitDB(config)
	assert.NoError(t, err)

	// Test CloseDB
	err = CloseDB()
	assert.NoError(t, err)
}

func TestGetDBBeforeInit(t *testing.T) {
	// Clean up any existing global connection
	CloseDB()

	// GetDB should panic if called before InitDB
	assert.Panics(t, func() {
		GetDB()
	})
}
