// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package db provides database access layer and repository implementations.
package db

import (
	"database/sql"
	"time"
)

// Activity represents an activity log entry in the database.
type Activity struct {
	ID         int64
	UserID     *int64
	UserEmail  *string
	Action     string
	EntityType string
	EntityID   *int64
	Details    *string
	IPAddress  *string
	UserAgent  *string
	CreatedAt  time.Time
}

// ActivityRepository handles database operations for activity logs.
type ActivityRepository struct {
	db *sql.DB
}

// NewActivityRepository creates a new activity repository.
func NewActivityRepository(db *sql.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

// Create inserts a new activity log entry into the database.
func (r *ActivityRepository) Create(activity *Activity) error {
	result, err := r.db.Exec(
		`INSERT INTO activities (
			user_id, user_email, action, entity_type, entity_id, details, ip_address, user_agent
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		activity.UserID,
		activity.UserEmail,
		activity.Action,
		activity.EntityType,
		activity.EntityID,
		activity.Details,
		activity.IPAddress,
		activity.UserAgent,
	)
	if err != nil {
		return err
	}

	activity.ID, err = result.LastInsertId()
	return err
}

// GetByID retrieves an activity log entry by ID.
func (r *ActivityRepository) GetByID(id int64) (*Activity, error) {
	activity := &Activity{}
	err := r.db.QueryRow(
		`SELECT
			id, user_id, user_email, action, entity_type, entity_id, details, ip_address, user_agent, created_at
		FROM activities
		WHERE id = ?`,
		id,
	).Scan(
		&activity.ID,
		&activity.UserID,
		&activity.UserEmail,
		&activity.Action,
		&activity.EntityType,
		&activity.EntityID,
		&activity.Details,
		&activity.IPAddress,
		&activity.UserAgent,
		&activity.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return activity, nil
}

// List retrieves all activity logs with pagination.
func (r *ActivityRepository) List(limit, offset int) ([]*Activity, error) {
	rows, err := r.db.Query(
		`SELECT
			id, user_id, user_email, action, entity_type, entity_id, details, ip_address, user_agent, created_at
		FROM activities
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActivities(rows)
}

// ListByUser retrieves all activity logs for a specific user.
func (r *ActivityRepository) ListByUser(userID int64, limit, offset int) ([]*Activity, error) {
	rows, err := r.db.Query(
		`SELECT
			id, user_id, user_email, action, entity_type, entity_id, details, ip_address, user_agent, created_at
		FROM activities
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActivities(rows)
}

// ListByAction retrieves all activity logs for a specific action.
func (r *ActivityRepository) ListByAction(action string, limit, offset int) ([]*Activity, error) {
	rows, err := r.db.Query(
		`SELECT
			id, user_id, user_email, action, entity_type, entity_id, details, ip_address, user_agent, created_at
		FROM activities
		WHERE action = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`,
		action,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActivities(rows)
}

// ListByEntity retrieves all activity logs for a specific entity.
func (r *ActivityRepository) ListByEntity(entityType string, entityID int64, limit, offset int) ([]*Activity, error) {
	rows, err := r.db.Query(
		`SELECT
			id, user_id, user_email, action, entity_type, entity_id, details, ip_address, user_agent, created_at
		FROM activities
		WHERE entity_type = ? AND entity_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`,
		entityType,
		entityID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActivities(rows)
}

// ListByDateRange retrieves activity logs within a date range.
func (r *ActivityRepository) ListByDateRange(startDate, endDate time.Time, limit, offset int) ([]*Activity, error) {
	rows, err := r.db.Query(
		`SELECT
			id, user_id, user_email, action, entity_type, entity_id, details, ip_address, user_agent, created_at
		FROM activities
		WHERE created_at >= ? AND created_at <= ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`,
		startDate,
		endDate,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanActivities(rows)
}

// Count returns the total number of activity logs.
func (r *ActivityRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM activities").Scan(&count)
	return count, err
}

// CountByUser returns the total number of activity logs for a specific user.
func (r *ActivityRepository) CountByUser(userID int64) (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM activities WHERE user_id = ?", userID).Scan(&count)
	return count, err
}

// DeleteOlderThan removes activity logs older than a specific date (for cleanup).
func (r *ActivityRepository) DeleteOlderThan(date time.Time) (int64, error) {
	result, err := r.db.Exec("DELETE FROM activities WHERE created_at < ?", date)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *ActivityRepository) scanActivities(rows *sql.Rows) ([]*Activity, error) {
	var activities []*Activity
	for rows.Next() {
		activity := &Activity{}
		if err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.UserEmail,
			&activity.Action,
			&activity.EntityType,
			&activity.EntityID,
			&activity.Details,
			&activity.IPAddress,
			&activity.UserAgent,
			&activity.CreatedAt,
		); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, rows.Err()
}
