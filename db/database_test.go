// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitSQLiteConnection(t *testing.T) {
	tmpFile := "/tmp/test_mut.db"
	defer os.Remove(tmpFile)

	config := Config{
		Driver:     "sqlite",
		DataSource: tmpFile,
	}

	conn, err := NewConnection(config)
	assert.NoError(t, err)
	assert.NotNil(t, conn)
	err = conn.Ping()
	assert.NoError(t, err)
	err = conn.Close()
	assert.NoError(t, err)
}

func TestUnitUnsupportedDriver(t *testing.T) {
	config := Config{
		Driver: "mysql",
	}

	conn, err := NewConnection(config)
	assert.Error(t, err)
	assert.Nil(t, conn)
	assert.Contains(t, err.Error(), "unsupported database driver")
}

func TestUnitInitDBAndGetDB(t *testing.T) {
	CloseDB()

	tmpFile := "/tmp/test_mut_global.db"
	defer os.Remove(tmpFile)

	config := Config{
		Driver:     "sqlite",
		DataSource: tmpFile,
	}

	err := InitDB(config)
	assert.NoError(t, err)
	db := GetDB()
	assert.NotNil(t, db)
	err = db.Ping()
	assert.NoError(t, err)
	err = InitDB(config)
	assert.NoError(t, err)
	err = CloseDB()
	assert.NoError(t, err)
}

func TestUnitGetDBBeforeInit(t *testing.T) {
	CloseDB()

	db := GetDB()
	assert.Nil(t, db)
}
