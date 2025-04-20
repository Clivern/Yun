// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupSessionTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	// Create users table
	_, err = db.Exec(`
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

	// Create sessions table
	_, err = db.Exec(`
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

	return db
}

func TestUnitSessionRepository_Create(t *testing.T) {
	t.Run("Create session successfully", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		ipAddress := "192.168.1.1"
		userAgent := "Mozilla/5.0"
		session := &Session{
			Token:     "test-token-123",
			UserID:    user.ID,
			IPAddress: &ipAddress,
			UserAgent: &userAgent,
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		}

		err = sessionRepo.Create(session)
		assert.NoError(t, err)
		assert.NotZero(t, session.ID)
	})

	t.Run("Create session without optional fields", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		session := &Session{
			Token:     "test-token-456",
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		}

		err = sessionRepo.Create(session)
		assert.NoError(t, err)
		assert.NotZero(t, session.ID)
	})
}

func TestUnitSessionRepository_GetByToken(t *testing.T) {
	t.Run("Get existing session by token", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		session := &Session{
			Token:     "unique-token",
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		}
		err = sessionRepo.Create(session)
		assert.NoError(t, err)

		retrieved, err := sessionRepo.GetByToken("unique-token")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, session.Token, retrieved.Token)
		assert.Equal(t, session.UserID, retrieved.UserID)
	})

	t.Run("Get non-existent session returns nil", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		sessionRepo := NewSessionRepository(db)

		retrieved, err := sessionRepo.GetByToken("non-existent-token")
		assert.NoError(t, err)
		assert.Nil(t, retrieved)
	})
}

func TestUnitSessionRepository_GetByID(t *testing.T) {
	t.Run("Get existing session by ID", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		session := &Session{
			Token:     "test-token",
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		}
		err = sessionRepo.Create(session)
		assert.NoError(t, err)

		retrieved, err := sessionRepo.GetByID(session.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, session.ID, retrieved.ID)
		assert.Equal(t, session.Token, retrieved.Token)
	})

	t.Run("Get non-existent session by ID returns nil", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		sessionRepo := NewSessionRepository(db)

		retrieved, err := sessionRepo.GetByID(999)
		assert.NoError(t, err)
		assert.Nil(t, retrieved)
	})
}

func TestUnitSessionRepository_GetByUserID(t *testing.T) {
	t.Run("Get all sessions for a user", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		for i := 0; i < 3; i++ {
			session := &Session{
				Token:     "token-" + string(rune(i)),
				UserID:    user.ID,
				ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
			}
			err = sessionRepo.Create(session)
			assert.NoError(t, err)
		}

		sessions, err := sessionRepo.GetByUserID(user.ID)
		assert.NoError(t, err)
		assert.Len(t, sessions, 3)
	})

	t.Run("Get sessions for user with no sessions", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		sessions, err := sessionRepo.GetByUserID(user.ID)
		assert.NoError(t, err)
		assert.Empty(t, sessions)
	})
}

func TestUnitSessionRepository_Delete(t *testing.T) {
	t.Run("Delete existing session", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		session := &Session{
			Token:     "test-token",
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		}
		err = sessionRepo.Create(session)
		assert.NoError(t, err)

		err = sessionRepo.Delete(session.ID)
		assert.NoError(t, err)

		retrieved, err := sessionRepo.GetByID(session.ID)
		assert.NoError(t, err)
		assert.Nil(t, retrieved)
	})
}

func TestUnitSessionRepository_DeleteByToken(t *testing.T) {
	t.Run("Delete session by token", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		session := &Session{
			Token:     "test-token",
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		}
		err = sessionRepo.Create(session)
		assert.NoError(t, err)

		err = sessionRepo.DeleteByToken("test-token")
		assert.NoError(t, err)

		retrieved, err := sessionRepo.GetByToken("test-token")
		assert.NoError(t, err)
		assert.Nil(t, retrieved)
	})
}

func TestUnitSessionRepository_DeleteByUserID(t *testing.T) {
	t.Run("Delete all sessions for a user", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		for i := 0; i < 3; i++ {
			session := &Session{
				Token:     "token-" + string(rune(i)),
				UserID:    user.ID,
				ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
			}
			err = sessionRepo.Create(session)
			assert.NoError(t, err)
		}

		err = sessionRepo.DeleteByUserID(user.ID)
		assert.NoError(t, err)

		sessions, err := sessionRepo.GetByUserID(user.ID)
		assert.NoError(t, err)
		assert.Empty(t, sessions)
	})
}

func TestUnitSessionRepository_DeleteExpired(t *testing.T) {
	t.Run("Delete expired sessions", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		expiredSession := &Session{
			Token:     "expired-token",
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(-1 * time.Hour),
		}
		err = sessionRepo.Create(expiredSession)
		assert.NoError(t, err)

		activeSession := &Session{
			Token:     "active-token",
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		}
		err = sessionRepo.Create(activeSession)
		assert.NoError(t, err)

		deleted, err := sessionRepo.DeleteExpired()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), deleted)

		retrieved, err := sessionRepo.GetByToken("expired-token")
		assert.NoError(t, err)
		assert.Nil(t, retrieved)

		retrieved, err = sessionRepo.GetByToken("active-token")
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
	})
}

func TestUnitSessionRepository_IsValid(t *testing.T) {
	t.Run("Valid non-expired session", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		session := &Session{
			Token:     "valid-token",
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		}
		err = sessionRepo.Create(session)
		assert.NoError(t, err)

		isValid, err := sessionRepo.IsValid("valid-token")
		assert.NoError(t, err)
		assert.True(t, isValid)
	})

	t.Run("Expired session is invalid", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		session := &Session{
			Token:     "expired-token",
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(-1 * time.Hour),
		}
		err = sessionRepo.Create(session)
		assert.NoError(t, err)

		isValid, err := sessionRepo.IsValid("expired-token")
		assert.NoError(t, err)
		assert.False(t, isValid)
	})

	t.Run("Non-existent session is invalid", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		sessionRepo := NewSessionRepository(db)
		isValid, err := sessionRepo.IsValid("non-existent")
		assert.NoError(t, err)
		assert.False(t, isValid)
	})
}

func TestUnitSessionRepository_UpdateExpiration(t *testing.T) {
	t.Run("Update session expiration", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		oldExpiration := time.Now().UTC().Add(1 * time.Hour)
		session := &Session{
			Token:     "test-token",
			UserID:    user.ID,
			ExpiresAt: oldExpiration,
		}
		err = sessionRepo.Create(session)
		assert.NoError(t, err)

		newExpiration := time.Now().UTC().Add(48 * time.Hour)

		err = sessionRepo.UpdateExpiration(session.ID, newExpiration)
		assert.NoError(t, err)

		retrieved, err := sessionRepo.GetByID(session.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.True(t, retrieved.ExpiresAt.After(oldExpiration))
	})
}

func TestUnitSessionRepository_Count(t *testing.T) {
	t.Run("Count active sessions", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user := &User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Role:     "user",
			IsActive: true,
		}
		err := userRepo.Create(user)
		assert.NoError(t, err)

		for i := 0; i < 2; i++ {
			session := &Session{
				Token:     "active-token-" + string(rune(i)),
				UserID:    user.ID,
				ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
			}
			err = sessionRepo.Create(session)
			assert.NoError(t, err)
		}

		expiredSession := &Session{
			Token:     "expired-token",
			UserID:    user.ID,
			ExpiresAt: time.Now().UTC().Add(-1 * time.Hour),
		}
		err = sessionRepo.Create(expiredSession)
		assert.NoError(t, err)
		count, err := sessionRepo.Count()
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})
}

func TestUnitSessionRepository_CountByUserID(t *testing.T) {
	t.Run("Count active sessions for user", func(t *testing.T) {
		db := setupSessionTestDB(t)
		defer db.Close()

		userRepo := NewUserRepository(db)
		sessionRepo := NewSessionRepository(db)

		user1 := &User{
			Email:    "user1@example.com",
			Password: "hashedpassword",
			Role:     "user",
			APIKey:   "api-key-user1",
			IsActive: true,
		}
		err := userRepo.Create(user1)
		assert.NoError(t, err)

		user2 := &User{
			Email:    "user2@example.com",
			Password: "hashedpassword",
			Role:     "user",
			APIKey:   "api-key-user2",
			IsActive: true,
		}
		err = userRepo.Create(user2)
		assert.NoError(t, err)

		for i := 0; i < 2; i++ {
			session := &Session{
				Token:     "user1-token-" + string(rune(i)),
				UserID:    user1.ID,
				ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
			}
			err = sessionRepo.Create(session)
			assert.NoError(t, err)
		}

		session := &Session{
			Token:     "user2-token",
			UserID:    user2.ID,
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		}
		err = sessionRepo.Create(session)
		assert.NoError(t, err)

		count, err := sessionRepo.CountByUserID(user1.ID)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})
}
