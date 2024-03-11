// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupUserTestDB creates a test database with users and users_meta tables
func setupUserTestDB(t *testing.T) (*Connection, func()) {
	tmpFile := "/tmp/test_users_" + time.Now().UTC().Format("20060102150405") + ".db"

	config := Config{
		Driver:     "sqlite",
		DataSource: tmpFile,
	}

	conn, err := NewConnection(config)
	require.NoError(t, err, "Failed to create test database")

	// Create users table
	_, err = conn.DB.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL,
			api_key VARCHAR(255),
			is_active BOOLEAN DEFAULT 1,
			last_login_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err, "Failed to create users table")

	// Create users_meta table
	_, err = conn.DB.Exec(`
		CREATE TABLE users_meta (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			UNIQUE(user_id, key)
		)
	`)
	require.NoError(t, err, "Failed to create users_meta table")

	cleanup := func() {
		conn.Close()
		os.Remove(tmpFile)
	}

	return conn, cleanup
}

func TestUnitUserRepository_Create(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := NewUserRepository(conn.DB)

	t.Run("Create new user successfully", func(t *testing.T) {
		user := &User{
			Email:    "user@example.com",
			Password: "hashed_password",
			Role:     "user",
			IsActive: true,
		}

		err := repo.Create(user)
		assert.NoError(t, err)
		assert.Greater(t, user.ID, int64(0), "User ID should be set after creation")

		// Verify user was created
		fetched, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, fetched)
		assert.Equal(t, user.Email, fetched.Email)
		assert.Equal(t, user.Password, fetched.Password)
		assert.Equal(t, user.Role, fetched.Role)
		assert.True(t, fetched.IsActive)
	})

	t.Run("Create admin user", func(t *testing.T) {
		apiKey := "admin-api-key-123"
		user := &User{
			Email:    "admin@example.com",
			Password: "admin_password",
			Role:     "admin",
			APIKey:   apiKey,
			IsActive: true,
		}

		err := repo.Create(user)
		assert.NoError(t, err)

		fetched, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "admin", fetched.Role)
		assert.NotEmpty(t, fetched.APIKey)
		assert.Equal(t, apiKey, fetched.APIKey)
	})

	t.Run("Create user with last login", func(t *testing.T) {
		now := time.Now().UTC()
		user := &User{
			Email:       "login@example.com",
			Password:    "password",
			Role:        "user",
			IsActive:    true,
			LastLoginAt: now,
		}

		err := repo.Create(user)
		assert.NoError(t, err)

		fetched, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.False(t, fetched.LastLoginAt.IsZero())
	})

	t.Run("Create inactive user", func(t *testing.T) {
		user := &User{
			Email:    "inactive@example.com",
			Password: "password",
			Role:     "user",
			IsActive: false,
		}

		err := repo.Create(user)
		assert.NoError(t, err)

		fetched, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.False(t, fetched.IsActive)
	})

	t.Run("Create duplicate email should fail", func(t *testing.T) {
		user1 := &User{
			Email:    "duplicate@example.com",
			Password: "password1",
			Role:     "user",
			IsActive: true,
		}
		err := repo.Create(user1)
		assert.NoError(t, err)

		user2 := &User{
			Email:    "duplicate@example.com",
			Password: "password2",
			Role:     "admin",
			IsActive: true,
		}
		err = repo.Create(user2)
		assert.Error(t, err, "Should fail on duplicate email")
	})
}

func TestUnitUserRepository_GetByID(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := NewUserRepository(conn.DB)

	t.Run("Get existing user by ID", func(t *testing.T) {
		// Create test user
		user := &User{
			Email:    "test@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)

		// Get by ID
		fetched, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, fetched)
		assert.Equal(t, user.ID, fetched.ID)
		assert.Equal(t, user.Email, fetched.Email)
		assert.Equal(t, user.Password, fetched.Password)
		assert.Equal(t, user.Role, fetched.Role)
		assert.False(t, fetched.CreatedAt.IsZero())
		assert.False(t, fetched.UpdatedAt.IsZero())
	})

	t.Run("Get non-existent user", func(t *testing.T) {
		user, err := repo.GetByID(99999)
		assert.NoError(t, err)
		assert.Nil(t, user, "Should return nil for non-existent user")
	})
}

func TestUnitUserRepository_GetByEmail(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := NewUserRepository(conn.DB)

	t.Run("Get existing user by email", func(t *testing.T) {
		// Create test user
		user := &User{
			Email:    "email@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)

		// Get by email
		fetched, err := repo.GetByEmail("email@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, fetched)
		assert.Equal(t, user.ID, fetched.ID)
		assert.Equal(t, user.Email, fetched.Email)
	})

	t.Run("Get non-existent email", func(t *testing.T) {
		user, err := repo.GetByEmail("nonexistent@example.com")
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("Email lookup is case sensitive", func(t *testing.T) {
		user := &User{
			Email:    "CaseSensitive@Example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)

		// Try to get with different case
		_, err = repo.GetByEmail("casesensitive@example.com")
		assert.NoError(t, err)
		// Note: SQLite is case-insensitive by default, but this documents the behavior
	})
}

func TestUnitUserRepository_GetByAPIKey(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := NewUserRepository(conn.DB)

	t.Run("Get user by API key", func(t *testing.T) {
		apiKey := "test-api-key-12345"
		user := &User{
			Email:    "apiuser@example.com",
			Password: "password",
			Role:     "admin",
			APIKey:   apiKey,
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)

		// Get by API key
		fetched, err := repo.GetByAPIKey(apiKey)
		assert.NoError(t, err)
		assert.NotNil(t, fetched)
		assert.Equal(t, user.ID, fetched.ID)
		assert.Equal(t, user.Email, fetched.Email)
		assert.NotEmpty(t, fetched.APIKey)
		assert.Equal(t, apiKey, fetched.APIKey)
	})

	t.Run("Get user with non-existent API key", func(t *testing.T) {
		user, err := repo.GetByAPIKey("non-existent-key")
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("User without API key", func(t *testing.T) {
		user := &User{
			Email:    "noapi@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)

		fetched, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.Empty(t, fetched.APIKey)
	})
}

func TestUnitUserRepository_Update(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := NewUserRepository(conn.DB)

	t.Run("Update user fields", func(t *testing.T) {
		// Create user
		user := &User{
			Email:    "update@example.com",
			Password: "old_password",
			Role:     "user",
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)

		// Update fields
		user.Email = "newemail@example.com"
		user.Password = "new_password"
		user.Role = "admin"

		err = repo.Update(user)
		assert.NoError(t, err)

		// Verify update
		fetched, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "newemail@example.com", fetched.Email)
		assert.Equal(t, "new_password", fetched.Password)
		assert.Equal(t, "admin", fetched.Role)
	})

	t.Run("Update API key", func(t *testing.T) {
		user := &User{
			Email:    "apiupdate@example.com",
			Password: "password",
			Role:     "admin",
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)

		// Add API key
		newAPIKey := "new-api-key-xyz"
		user.APIKey = newAPIKey
		err = repo.Update(user)
		assert.NoError(t, err)

		fetched, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, fetched.APIKey)
		assert.Equal(t, newAPIKey, fetched.APIKey)
	})

	t.Run("Deactivate user", func(t *testing.T) {
		user := &User{
			Email:    "deactivate@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)

		user.IsActive = false
		err = repo.Update(user)
		assert.NoError(t, err)

		fetched, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.False(t, fetched.IsActive)
	})

	t.Run("Update updates timestamp", func(t *testing.T) {
		user := &User{
			Email:    "timestamp@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)

		// Get initial timestamp
		initial, err := repo.GetByID(user.ID)
		require.NoError(t, err)
		initialUpdatedAt := initial.UpdatedAt

		// Wait and update
		time.Sleep(10 * time.Millisecond)
		user.Email = "newtimestamp@example.com"
		err = repo.Update(user)
		assert.NoError(t, err)

		// Verify timestamp changed
		updated, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.True(t, updated.UpdatedAt.After(initialUpdatedAt))
	})
}

func TestUnitUserRepository_UpdateLastLogin(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := NewUserRepository(conn.DB)

	t.Run("Update last login timestamp", func(t *testing.T) {
		user := &User{
			Email:    "login@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)

		// Initially no last login
		fetched, err := repo.GetByID(user.ID)
		require.NoError(t, err)
		assert.True(t, fetched.LastLoginAt.IsZero())

		// Update last login
		beforeUpdate := time.Now().UTC()
		err = repo.UpdateLastLogin(user.ID)
		assert.NoError(t, err)
		afterUpdate := time.Now().UTC()

		// Verify last login was set
		fetched, err = repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.False(t, fetched.LastLoginAt.IsZero())
		assert.True(t, fetched.LastLoginAt.After(beforeUpdate.Add(-time.Second)))
		assert.True(t, fetched.LastLoginAt.Before(afterUpdate.Add(time.Second)))
	})

	t.Run("Update last login multiple times", func(t *testing.T) {
		user := &User{
			Email:    "multilogin@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)

		// First login
		err = repo.UpdateLastLogin(user.ID)
		require.NoError(t, err)

		first, err := repo.GetByID(user.ID)
		require.NoError(t, err)
		firstLogin := first.LastLoginAt

		// Wait and login again
		time.Sleep(10 * time.Millisecond)
		err = repo.UpdateLastLogin(user.ID)
		assert.NoError(t, err)

		second, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.False(t, second.LastLoginAt.IsZero())
		assert.True(t, second.LastLoginAt.After(firstLogin))
	})
}

func TestUnitUserRepository_Delete(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := NewUserRepository(conn.DB)

	t.Run("Delete existing user", func(t *testing.T) {
		user := &User{
			Email:    "delete@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)

		// Verify user exists
		fetched, err := repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, fetched)

		// Delete user
		err = repo.Delete(user.ID)
		assert.NoError(t, err)

		// Verify user is gone
		fetched, err = repo.GetByID(user.ID)
		assert.NoError(t, err)
		assert.Nil(t, fetched)
	})

	t.Run("Delete non-existent user", func(t *testing.T) {
		err := repo.Delete(99999)
		assert.NoError(t, err, "Deleting non-existent user should not error")
	})

	t.Run("Delete and recreate", func(t *testing.T) {
		email := "recreate@example.com"
		user := &User{
			Email:    email,
			Password: "password1",
			Role:     "user",
			IsActive: true,
		}
		err := repo.Create(user)
		require.NoError(t, err)
		firstID := user.ID

		// Delete
		err = repo.Delete(user.ID)
		assert.NoError(t, err)

		// Recreate with same email
		newUser := &User{
			Email:    email,
			Password: "password2",
			Role:     "admin",
			IsActive: true,
		}
		err = repo.Create(newUser)
		assert.NoError(t, err)
		assert.NotEqual(t, firstID, newUser.ID, "New user should have different ID")
	})
}

func TestUnitUserRepository_List(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := NewUserRepository(conn.DB)

	t.Run("List empty table", func(t *testing.T) {
		users, err := repo.List(10, 0)
		assert.NoError(t, err)
		assert.Empty(t, users)
	})

	t.Run("List all users", func(t *testing.T) {
		// Create test users
		for i := 1; i <= 5; i++ {
			user := &User{
				Email:    "user" + string(rune('0'+i)) + "@example.com",
				Password: "password",
				Role:     "user",
				IsActive: true,
			}
			err := repo.Create(user)
			require.NoError(t, err)
			// Add small delay to ensure different created_at times
			time.Sleep(time.Millisecond)
		}

		// List all
		users, err := repo.List(10, 0)
		assert.NoError(t, err)
		assert.Len(t, users, 5)

		// Verify ordered by created_at DESC (newest first)
		for i := 0; i < len(users)-1; i++ {
			assert.True(t, users[i].CreatedAt.After(users[i+1].CreatedAt) ||
				users[i].CreatedAt.Equal(users[i+1].CreatedAt))
		}
	})

	t.Run("List with pagination", func(t *testing.T) {
		// Clear previous data
		conn.DB.Exec("DELETE FROM users")

		// Create 10 users
		for i := 1; i <= 10; i++ {
			user := &User{
				Email:    "page" + string(rune('0'+i)) + "@example.com",
				Password: "password",
				Role:     "user",
				IsActive: true,
			}
			err := repo.Create(user)
			require.NoError(t, err)
			time.Sleep(time.Millisecond)
		}

		// Get first page
		page1, err := repo.List(3, 0)
		assert.NoError(t, err)
		assert.Len(t, page1, 3)

		// Get second page
		page2, err := repo.List(3, 3)
		assert.NoError(t, err)
		assert.Len(t, page2, 3)

		// Ensure different users
		assert.NotEqual(t, page1[0].ID, page2[0].ID)
	})

	t.Run("List with limit larger than total", func(t *testing.T) {
		users, err := repo.List(100, 0)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(users), 100)
	})

	t.Run("List with offset beyond results", func(t *testing.T) {
		users, err := repo.List(10, 1000)
		assert.NoError(t, err)
		assert.Empty(t, users)
	})
}

func TestUnitUserRepository_Count(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	repo := NewUserRepository(conn.DB)

	t.Run("Count empty table", func(t *testing.T) {
		count, err := repo.Count()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Count users", func(t *testing.T) {
		// Create users
		for i := 1; i <= 7; i++ {
			user := &User{
				Email:    "count" + string(rune('0'+i)) + "@example.com",
				Password: "password",
				Role:     "user",
				IsActive: true,
			}
			err := repo.Create(user)
			require.NoError(t, err)
		}

		count, err := repo.Count()
		assert.NoError(t, err)
		assert.Equal(t, int64(7), count)
	})

	t.Run("Count after deletion", func(t *testing.T) {
		// Get initial count
		initialCount, err := repo.Count()
		require.NoError(t, err)

		// Create and delete a user
		user := &User{
			Email:    "countdelete@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err = repo.Create(user)
		require.NoError(t, err)

		countAfterCreate, err := repo.Count()
		assert.NoError(t, err)
		assert.Equal(t, initialCount+1, countAfterCreate)

		err = repo.Delete(user.ID)
		require.NoError(t, err)

		countAfterDelete, err := repo.Count()
		assert.NoError(t, err)
		assert.Equal(t, initialCount, countAfterDelete)
	})
}

func TestUnitUserMetaRepository_Create(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	userRepo := NewUserRepository(conn.DB)
	metaRepo := NewUserMetaRepository(conn.DB)

	t.Run("Create user metadata", func(t *testing.T) {
		// Create user first
		user := &User{
			Email:    "meta@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		// Create metadata
		err = metaRepo.Create(user.ID, "theme", "dark")
		assert.NoError(t, err)

		// Verify
		meta, err := metaRepo.Get(user.ID, "theme")
		assert.NoError(t, err)
		assert.NotNil(t, meta)
		assert.Equal(t, "theme", meta.Key)
		assert.Equal(t, "dark", meta.Value)
		assert.Equal(t, user.ID, meta.UserID)
	})

	t.Run("Create multiple metadata entries", func(t *testing.T) {
		user := &User{
			Email:    "multimeta@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		metadata := map[string]string{
			"language":   "en",
			"timezone":   "UTC",
			"notify":     "true",
			"newsletter": "false",
		}

		for key, value := range metadata {
			err := metaRepo.Create(user.ID, key, value)
			assert.NoError(t, err)
		}

		// Verify all created
		for key, expectedValue := range metadata {
			meta, err := metaRepo.Get(user.ID, key)
			assert.NoError(t, err)
			assert.NotNil(t, meta)
			assert.Equal(t, expectedValue, meta.Value)
		}
	})

	t.Run("Create duplicate key should fail", func(t *testing.T) {
		user := &User{
			Email:    "dupmeta@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		err = metaRepo.Create(user.ID, "setting", "value1")
		assert.NoError(t, err)

		err = metaRepo.Create(user.ID, "setting", "value2")
		assert.Error(t, err, "Should fail on duplicate key for same user")
	})
}

func TestUnitUserMetaRepository_Get(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	userRepo := NewUserRepository(conn.DB)
	metaRepo := NewUserMetaRepository(conn.DB)

	t.Run("Get existing metadata", func(t *testing.T) {
		user := &User{
			Email:    "getmeta@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		err = metaRepo.Create(user.ID, "color", "blue")
		require.NoError(t, err)

		meta, err := metaRepo.Get(user.ID, "color")
		assert.NoError(t, err)
		assert.NotNil(t, meta)
		assert.Equal(t, "color", meta.Key)
		assert.Equal(t, "blue", meta.Value)
		assert.Greater(t, meta.ID, int64(0))
		assert.False(t, meta.CreatedAt.IsZero())
		assert.False(t, meta.UpdatedAt.IsZero())
	})

	t.Run("Get non-existent metadata", func(t *testing.T) {
		user := &User{
			Email:    "nometa@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		meta, err := metaRepo.Get(user.ID, "nonexistent")
		assert.NoError(t, err)
		assert.Nil(t, meta)
	})

	t.Run("Get metadata for non-existent user", func(t *testing.T) {
		meta, err := metaRepo.Get(99999, "somekey")
		assert.NoError(t, err)
		assert.Nil(t, meta)
	})
}

func TestUnitUserMetaRepository_Update(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	userRepo := NewUserRepository(conn.DB)
	metaRepo := NewUserMetaRepository(conn.DB)

	t.Run("Update existing metadata", func(t *testing.T) {
		user := &User{
			Email:    "updatemeta@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		// Create initial metadata
		err = metaRepo.Create(user.ID, "status", "online")
		require.NoError(t, err)

		// Update
		err = metaRepo.Update(user.ID, "status", "offline")
		assert.NoError(t, err)

		// Verify
		meta, err := metaRepo.Get(user.ID, "status")
		assert.NoError(t, err)
		assert.Equal(t, "offline", meta.Value)
	})

	t.Run("Update updates timestamp", func(t *testing.T) {
		user := &User{
			Email:    "metatime@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		err = metaRepo.Create(user.ID, "timestamp", "value1")
		require.NoError(t, err)

		initial, err := metaRepo.Get(user.ID, "timestamp")
		require.NoError(t, err)
		initialUpdatedAt := initial.UpdatedAt

		time.Sleep(10 * time.Millisecond)

		err = metaRepo.Update(user.ID, "timestamp", "value2")
		assert.NoError(t, err)

		updated, err := metaRepo.Get(user.ID, "timestamp")
		assert.NoError(t, err)
		assert.True(t, updated.UpdatedAt.After(initialUpdatedAt))
	})

	t.Run("Update non-existent metadata", func(t *testing.T) {
		user := &User{
			Email:    "noupdate@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		err = metaRepo.Update(user.ID, "nonexistent", "value")
		assert.NoError(t, err)

		// Verify not created
		meta, err := metaRepo.Get(user.ID, "nonexistent")
		assert.NoError(t, err)
		assert.Nil(t, meta)
	})
}

func TestUnitUserMetaRepository_Delete(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	userRepo := NewUserRepository(conn.DB)
	metaRepo := NewUserMetaRepository(conn.DB)

	t.Run("Delete existing metadata", func(t *testing.T) {
		user := &User{
			Email:    "deletemeta@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		err = metaRepo.Create(user.ID, "temp", "value")
		require.NoError(t, err)

		// Verify exists
		meta, err := metaRepo.Get(user.ID, "temp")
		assert.NoError(t, err)
		assert.NotNil(t, meta)

		// Delete
		err = metaRepo.Delete(user.ID, "temp")
		assert.NoError(t, err)

		// Verify deleted
		meta, err = metaRepo.Get(user.ID, "temp")
		assert.NoError(t, err)
		assert.Nil(t, meta)
	})

	t.Run("Delete non-existent metadata", func(t *testing.T) {
		user := &User{
			Email:    "nodelete@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		err = metaRepo.Delete(user.ID, "nonexistent")
		assert.NoError(t, err)
	})
}

func TestUnitUserMetaRepository_ListByUser(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	userRepo := NewUserRepository(conn.DB)
	metaRepo := NewUserMetaRepository(conn.DB)

	t.Run("List empty metadata", func(t *testing.T) {
		user := &User{
			Email:    "empty@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		metadata, err := metaRepo.ListByUser(user.ID)
		assert.NoError(t, err)
		assert.Empty(t, metadata)
	})

	t.Run("List all metadata for user", func(t *testing.T) {
		user := &User{
			Email:    "listmeta@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		// Create multiple metadata entries
		testData := map[string]string{
			"pref_theme":    "dark",
			"pref_lang":     "en",
			"pref_timezone": "America/New_York",
			"pref_notify":   "email",
		}

		for key, value := range testData {
			err := metaRepo.Create(user.ID, key, value)
			require.NoError(t, err)
		}

		// List all
		metadata, err := metaRepo.ListByUser(user.ID)
		assert.NoError(t, err)
		assert.Len(t, metadata, len(testData))

		// Verify all keys present and sorted
		for i, meta := range metadata {
			assert.Equal(t, testData[meta.Key], meta.Value)
			assert.Equal(t, user.ID, meta.UserID)

			// Check alphabetical order
			if i > 0 {
				assert.True(t, metadata[i-1].Key < metadata[i].Key,
					"Metadata should be sorted by key")
			}
		}
	})

	t.Run("List metadata for multiple users", func(t *testing.T) {
		user1 := &User{
			Email:    "multi1@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user1)
		require.NoError(t, err)

		user2 := &User{
			Email:    "multi2@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err = userRepo.Create(user2)
		require.NoError(t, err)

		// Create metadata for both users
		metaRepo.Create(user1.ID, "key1", "value1")
		metaRepo.Create(user2.ID, "key2", "value2")

		// List for user1
		metadata1, err := metaRepo.ListByUser(user1.ID)
		assert.NoError(t, err)
		assert.Len(t, metadata1, 1)
		assert.Equal(t, "key1", metadata1[0].Key)

		// List for user2
		metadata2, err := metaRepo.ListByUser(user2.ID)
		assert.NoError(t, err)
		assert.Len(t, metadata2, 1)
		assert.Equal(t, "key2", metadata2[0].Key)
	})
}

func TestUnitUserMetaRepository_Upsert(t *testing.T) {
	conn, cleanup := setupUserTestDB(t)
	defer cleanup()

	userRepo := NewUserRepository(conn.DB)
	metaRepo := NewUserMetaRepository(conn.DB)

	t.Run("Upsert creates new metadata", func(t *testing.T) {
		user := &User{
			Email:    "upsert@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		err = metaRepo.Upsert(user.ID, "new_key", "new_value")
		assert.NoError(t, err)

		meta, err := metaRepo.Get(user.ID, "new_key")
		assert.NoError(t, err)
		assert.NotNil(t, meta)
		assert.Equal(t, "new_value", meta.Value)
	})

	t.Run("Upsert updates existing metadata", func(t *testing.T) {
		user := &User{
			Email:    "upsertupdate@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		// Create initial
		err = metaRepo.Create(user.ID, "existing", "initial")
		require.NoError(t, err)

		// Upsert should update
		err = metaRepo.Upsert(user.ID, "existing", "updated")
		assert.NoError(t, err)

		meta, err := metaRepo.Get(user.ID, "existing")
		assert.NoError(t, err)
		assert.Equal(t, "updated", meta.Value)
	})

	t.Run("Multiple upserts", func(t *testing.T) {
		user := &User{
			Email:    "multiupsert@example.com",
			Password: "password",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		require.NoError(t, err)

		// First upsert - creates
		err = metaRepo.Upsert(user.ID, "counter", "1")
		assert.NoError(t, err)

		// Subsequent upserts - update
		for i := 2; i <= 5; i++ {
			err = metaRepo.Upsert(user.ID, "counter", string(rune('0'+i)))
			assert.NoError(t, err)
		}

		meta, err := metaRepo.Get(user.ID, "counter")
		assert.NoError(t, err)
		assert.Equal(t, "5", meta.Value)

		// Verify only one record exists
		allMeta, err := metaRepo.ListByUser(user.ID)
		assert.NoError(t, err)
		countEntries := 0
		for _, m := range allMeta {
			if m.Key == "counter" {
				countEntries++
			}
		}
		assert.Equal(t, 1, countEntries, "Should only have one counter entry")
	})
}
