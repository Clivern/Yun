// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Session represents a user session in the database.
type Session struct {
	ID        int64
	Token     string
	UserID    int64
	IPAddress *string
	UserAgent *string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SessionRepository handles database operations for sessions.
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new session repository.
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create inserts a new session into the database.
//
// Example:
//
//	session := &Session{
//		Token:     "secure-random-token",
//		UserID:    1,
//		ExpiresAt: time.Now().Add(24 * time.Hour),
//	}
//	err := repo.Create(session)
func (r *SessionRepository) Create(session *Session) error {
	result, err := r.db.Exec(
		`INSERT INTO sessions (token, user_id, ip_address, user_agent, expires_at)
		VALUES (?, ?, ?, ?, ?)`,
		session.Token,
		session.UserID,
		session.IPAddress,
		session.UserAgent,
		session.ExpiresAt,
	)
	if err != nil {
		return err
	}

	session.ID, err = result.LastInsertId()
	return err
}

// GetByToken retrieves a session by token.
func (r *SessionRepository) GetByToken(token string) (*Session, error) {
	session := &Session{}
	err := r.db.QueryRow(
		`SELECT id, token, user_id, ip_address, user_agent, expires_at, created_at, updated_at
		FROM sessions
		WHERE token = ?`,
		token,
	).Scan(
		&session.ID,
		&session.Token,
		&session.UserID,
		&session.IPAddress,
		&session.UserAgent,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetByID retrieves a session by ID.
func (r *SessionRepository) GetByID(id int64) (*Session, error) {
	session := &Session{}
	err := r.db.QueryRow(
		`SELECT id, token, user_id, ip_address, user_agent, expires_at, created_at, updated_at
		FROM sessions
		WHERE id = ?`,
		id,
	).Scan(
		&session.ID,
		&session.Token,
		&session.UserID,
		&session.IPAddress,
		&session.UserAgent,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetByUserID retrieves all sessions for a user.
func (r *SessionRepository) GetByUserID(userID int64) ([]*Session, error) {
	rows, err := r.db.Query(
		`SELECT id, token, user_id, ip_address, user_agent, expires_at, created_at, updated_at
		FROM sessions
		WHERE user_id = ?
		ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next() {
		session := &Session{}
		if err := rows.Scan(
			&session.ID,
			&session.Token,
			&session.UserID,
			&session.IPAddress,
			&session.UserAgent,
			&session.ExpiresAt,
			&session.CreatedAt,
			&session.UpdatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

// Delete removes a session from the database.
func (r *SessionRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE id = ?", id)
	return err
}

// DeleteByToken removes a session by token.
func (r *SessionRepository) DeleteByToken(token string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err
}

// DeleteByUserID removes all sessions for a user.
func (r *SessionRepository) DeleteByUserID(userID int64) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	return err
}

// DeleteExpired removes all expired sessions.
func (r *SessionRepository) DeleteExpired() (int64, error) {
	result, err := r.db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now().UTC())
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// IsValid checks if a session exists and is not expired.
func (r *SessionRepository) IsValid(token string) (bool, error) {
	session, err := r.GetByToken(token)
	if err != nil {
		return false, err
	}

	if session == nil {
		return false, nil
	}

	return session.ExpiresAt.After(time.Now().UTC()), nil
}

// UpdateExpiration updates the expiration time of a session.
func (r *SessionRepository) UpdateExpiration(id int64, expiresAt time.Time) error {
	_, err := r.db.Exec(
		`UPDATE sessions SET
			expires_at = ?, updated_at = ?
		WHERE id = ?`,
		expiresAt,
		time.Now().UTC(),
		id,
	)
	return err
}

// Count returns the total number of active (non-expired) sessions.
func (r *SessionRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM sessions WHERE expires_at > ?", time.Now().UTC()).Scan(&count)
	return count, err
}

// CountByUserID returns the number of active sessions for a user.
func (r *SessionRepository) CountByUserID(userID int64) (int64, error) {
	var count int64
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM sessions WHERE user_id = ? AND expires_at > ?`,
		userID,
		time.Now().UTC(),
	).Scan(&count)
	return count, err
}
