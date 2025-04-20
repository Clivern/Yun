// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"database/sql"
	"testing"
	"time"

	"github.com/clivern/mut/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupSessionModuleTestDB(t *testing.T) *sql.DB {
	testDB, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	_, err = testDB.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			api_key VARCHAR(255) UNIQUE,
			is_active BOOLEAN DEFAULT 1,
			last_login_at DATETIME NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	assert.NoError(t, err)

	_, err = testDB.Exec(`
		CREATE TABLE sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			token VARCHAR(255) NOT NULL UNIQUE,
			user_id INTEGER NOT NULL,
			ip_address VARCHAR(45),
			user_agent VARCHAR(500),
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)
	assert.NoError(t, err)

	return testDB
}

func TestUnitSessionManager_CreateSession(t *testing.T) {
	t.Run("Create session successfully", func(t *testing.T) {
		testDB := setupSessionModuleTestDB(t)
		defer testDB.Close()

		userRepo := db.NewUserRepository(testDB)
		sessionRepo := db.NewSessionRepository(testDB)
		sessionManager := NewSessionManager(sessionRepo, userRepo)

		user := &db.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		session, err := sessionManager.CreateSession(user.ID, 24*time.Hour, "192.168.1.1", "Mozilla/5.0")

		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.NotEmpty(t, session.Token)
		assert.Equal(t, user.ID, session.UserID)
		assert.NotNil(t, session.IPAddress)
		assert.Equal(t, "192.168.1.1", *session.IPAddress)
		assert.NotNil(t, session.UserAgent)
		assert.Equal(t, "Mozilla/5.0", *session.UserAgent)
		assert.True(t, session.ExpiresAt.After(time.Now().UTC()))
	})

	t.Run("Create session without optional fields", func(t *testing.T) {
		testDB := setupSessionModuleTestDB(t)
		defer testDB.Close()

		userRepo := db.NewUserRepository(testDB)
		sessionRepo := db.NewSessionRepository(testDB)
		sessionManager := NewSessionManager(sessionRepo, userRepo)

		user := &db.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		session, err := sessionManager.CreateSession(user.ID, 24*time.Hour, "", "")

		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.NotEmpty(t, session.Token)
		assert.Nil(t, session.IPAddress)
		assert.Nil(t, session.UserAgent)
	})

	t.Run("Create session for non-existent user fails", func(t *testing.T) {
		testDB := setupSessionModuleTestDB(t)
		defer testDB.Close()

		userRepo := db.NewUserRepository(testDB)
		sessionRepo := db.NewSessionRepository(testDB)
		sessionManager := NewSessionManager(sessionRepo, userRepo)

		session, err := sessionManager.CreateSession(999, 24*time.Hour, "", "")

		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Each session gets unique token", func(t *testing.T) {
		testDB := setupSessionModuleTestDB(t)
		defer testDB.Close()

		userRepo := db.NewUserRepository(testDB)
		sessionRepo := db.NewSessionRepository(testDB)
		sessionManager := NewSessionManager(sessionRepo, userRepo)

		user := &db.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		session1, err1 := sessionManager.CreateSession(user.ID, 24*time.Hour, "", "")
		session2, err2 := sessionManager.CreateSession(user.ID, 24*time.Hour, "", "")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, session1.Token, session2.Token)
	})
}

func TestUnitSessionManager_ValidateSession(t *testing.T) {
	t.Run("Validate valid session", func(t *testing.T) {
		testDB := setupSessionModuleTestDB(t)
		defer testDB.Close()

		userRepo := db.NewUserRepository(testDB)
		sessionRepo := db.NewSessionRepository(testDB)
		sessionManager := NewSessionManager(sessionRepo, userRepo)

		user := &db.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		session, err := sessionManager.CreateSession(user.ID, 24*time.Hour, "", "")
		assert.NoError(t, err)

		validUser, validSession, err := sessionManager.ValidateSession(session.Token)

		assert.NoError(t, err)
		assert.NotNil(t, validUser)
		assert.NotNil(t, validSession)
		assert.Equal(t, user.ID, validUser.ID)
		assert.Equal(t, session.Token, validSession.Token)
	})

	t.Run("Validate expired session fails", func(t *testing.T) {
		testDB := setupSessionModuleTestDB(t)
		defer testDB.Close()

		userRepo := db.NewUserRepository(testDB)
		sessionRepo := db.NewSessionRepository(testDB)
		sessionManager := NewSessionManager(sessionRepo, userRepo)

		user := &db.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		expiredSession := &db.Session{
			Token:     "expired-token",
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(-1 * time.Hour),
		}
		err = sessionRepo.Create(expiredSession)
		assert.NoError(t, err)

		validUser, validSession, err := sessionManager.ValidateSession(expiredSession.Token)

		assert.Error(t, err)
		assert.Nil(t, validUser)
		assert.Nil(t, validSession)
		assert.Contains(t, err.Error(), "expired")
	})

	t.Run("Validate non-existent session fails", func(t *testing.T) {
		testDB := setupSessionModuleTestDB(t)
		defer testDB.Close()

		userRepo := db.NewUserRepository(testDB)
		sessionRepo := db.NewSessionRepository(testDB)
		sessionManager := NewSessionManager(sessionRepo, userRepo)

		validUser, validSession, err := sessionManager.ValidateSession("non-existent-token")

		assert.Error(t, err)
		assert.Nil(t, validUser)
		assert.Nil(t, validSession)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Validate session for inactive user fails", func(t *testing.T) {
		// Arrange
		testDB := setupSessionModuleTestDB(t)
		defer testDB.Close()

		userRepo := db.NewUserRepository(testDB)
		sessionRepo := db.NewSessionRepository(testDB)
		sessionManager := NewSessionManager(sessionRepo, userRepo)

		user := &db.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		session, err := sessionManager.CreateSession(user.ID, 24*time.Hour, "", "")
		assert.NoError(t, err)

		// Deactivate user
		user.IsActive = false
		err = userRepo.Update(user)
		assert.NoError(t, err)

		// Act
		validUser, validSession, err := sessionManager.ValidateSession(session.Token)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, validUser)
		assert.Nil(t, validSession)
		assert.Contains(t, err.Error(), "not active")
	})
}

func TestUnitSessionManager_RevokeUserSessions(t *testing.T) {
	t.Run("Revoke all user sessions", func(t *testing.T) {
		testDB := setupSessionModuleTestDB(t)
		defer testDB.Close()

		userRepo := db.NewUserRepository(testDB)
		sessionRepo := db.NewSessionRepository(testDB)
		sessionManager := NewSessionManager(sessionRepo, userRepo)

		user := &db.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		for i := 0; i < 3; i++ {
			_, err := sessionManager.CreateSession(user.ID, 24*time.Hour, "", "")
			assert.NoError(t, err)
		}

		err = sessionManager.RevokeUserSessions(user.ID)
		assert.NoError(t, err)

		sessions, err := sessionRepo.GetByUserID(user.ID)
		assert.NoError(t, err)
		assert.Empty(t, sessions)
	})
}

func TestUnitSessionManager_GetUserSessions(t *testing.T) {
	t.Run("Get active user sessions", func(t *testing.T) {
		testDB := setupSessionModuleTestDB(t)
		defer testDB.Close()

		userRepo := db.NewUserRepository(testDB)
		sessionRepo := db.NewSessionRepository(testDB)
		sessionManager := NewSessionManager(sessionRepo, userRepo)

		user := &db.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		for i := 0; i < 2; i++ {
			_, err := sessionManager.CreateSession(user.ID, 24*time.Hour, "", "")
			assert.NoError(t, err)
		}

		expiredSession := &db.Session{
			Token:     "expired",
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(-1 * time.Hour),
		}
		err = sessionRepo.Create(expiredSession)
		assert.NoError(t, err)

		sessions, err := sessionManager.GetUserSessions(user.ID)
		assert.NoError(t, err)
		assert.Len(t, sessions, 2)
	})
}

func TestUnitSessionManager_CleanupExpiredSessions(t *testing.T) {
	t.Run("Cleanup expired sessions", func(t *testing.T) {
		// Arrange
		testDB := setupSessionModuleTestDB(t)
		defer testDB.Close()

		userRepo := db.NewUserRepository(testDB)
		sessionRepo := db.NewSessionRepository(testDB)
		sessionManager := NewSessionManager(sessionRepo, userRepo)

		user := &db.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		for i := 0; i < 3; i++ {
			expiredSession := &db.Session{
				Token:     "expired-" + string(rune(i)),
				UserID:    user.ID,
				ExpiresAt: time.Now().UTC().Add(-1 * time.Hour),
			}
			err = sessionRepo.Create(expiredSession)
			assert.NoError(t, err)
		}

		_, err = sessionManager.CreateSession(user.ID, 24*time.Hour, "", "")
		assert.NoError(t, err)

		deleted, err := sessionManager.CleanupExpiredSessions()
		assert.NoError(t, err)
		assert.Equal(t, int64(3), deleted)

		count, err := sessionRepo.Count()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}

func TestUnitGenerateSecureToken(t *testing.T) {
	t.Run("Generate secure token", func(t *testing.T) {
		token1, err1 := generateSecureToken(32)
		token2, err2 := generateSecureToken(32)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEmpty(t, token1)
		assert.NotEmpty(t, token2)
		assert.NotEqual(t, token1, token2, "Tokens should be unique")
	})

	t.Run("Generate token with different lengths", func(t *testing.T) {
		token16, err16 := generateSecureToken(16)
		token32, err32 := generateSecureToken(32)
		token64, err64 := generateSecureToken(64)

		assert.NoError(t, err16)
		assert.NoError(t, err32)
		assert.NoError(t, err64)
		assert.NotEmpty(t, token16)
		assert.NotEmpty(t, token32)
		assert.NotEmpty(t, token64)
		assert.True(t, len(token64) > len(token32))
		assert.True(t, len(token32) > len(token16))
	})
}
