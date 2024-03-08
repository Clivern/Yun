// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/clivern/yun/db"
)

// SessionManager handles session operations.
type SessionManager struct {
	SessionRepo *db.SessionRepository
	UserRepo    *db.UserRepository
}

// NewSessionManager creates a new session manager.
func NewSessionManager(sessionRepo *db.SessionRepository, userRepo *db.UserRepository) *SessionManager {
	return &SessionManager{
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
	}
}

// CreateSession creates a new session for a user.
//
// Example:
//
//	session, err := manager.CreateSession(userID, 24*time.Hour, "192.168.1.1", "Mozilla/5.0...")
func (s *SessionManager) CreateSession(userID int64, duration time.Duration, ipAddress, userAgent string) (*db.Session, error) {
	// Verify user exists
	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Generate secure random token
	token, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	// Create session
	session := &db.Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().UTC().Add(duration),
	}

	if ipAddress != "" {
		session.IPAddress = &ipAddress
	}

	if userAgent != "" {
		session.UserAgent = &userAgent
	}

	err = s.SessionRepo.Create(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// ValidateSession validates a session token and returns the associated user.
func (s *SessionManager) ValidateSession(token string) (*db.User, *db.Session, error) {
	session, err := s.SessionRepo.GetByToken(token)
	if err != nil {
		return nil, nil, err
	}

	if session == nil {
		return nil, nil, errors.New("session not found")
	}

	// Check if session is expired
	if session.ExpiresAt.Before(time.Now().UTC()) {
		// Clean up expired session
		s.SessionRepo.Delete(session.ID)
		return nil, nil, errors.New("session expired")
	}

	// Get user
	user, err := s.UserRepo.GetByID(session.UserID)
	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, errors.New("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, nil, errors.New("user is not active")
	}

	return user, session, nil
}

// RefreshSession extends the expiration time of a session.
func (s *SessionManager) RefreshSession(token string, duration time.Duration) error {
	session, err := s.SessionRepo.GetByToken(token)
	if err != nil {
		return err
	}

	if session == nil {
		return errors.New("session not found")
	}

	newExpiration := time.Now().UTC().Add(duration)
	return s.SessionRepo.UpdateExpiration(session.ID, newExpiration)
}

// RevokeSession revokes a session by token.
func (s *SessionManager) RevokeSession(token string) error {
	return s.SessionRepo.DeleteByToken(token)
}

// RevokeUserSessions revokes all sessions for a user.
func (s *SessionManager) RevokeUserSessions(userID int64) error {
	return s.SessionRepo.DeleteByUserID(userID)
}

// GetUserSessions retrieves all active sessions for a user.
func (s *SessionManager) GetUserSessions(userID int64) ([]*db.Session, error) {
	sessions, err := s.SessionRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Filter out expired sessions
	var activeSessions []*db.Session
	now := time.Now().UTC()
	for _, session := range sessions {
		if session.ExpiresAt.After(now) {
			activeSessions = append(activeSessions, session)
		}
	}

	return activeSessions, nil
}

// CleanupExpiredSessions removes all expired sessions from the database.
func (s *SessionManager) CleanupExpiredSessions() (int64, error) {
	return s.SessionRepo.DeleteExpired()
}

// generateSecureToken generates a cryptographically secure random token.
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}
