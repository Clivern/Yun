// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Option represents a key-value option in the database.
type Option struct {
	ID        int64
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OptionRepository handles database operations for options.
type OptionRepository struct {
	db *sql.DB
}

// NewOptionRepository creates a new option repository.
func NewOptionRepository(db *sql.DB) *OptionRepository {
	return &OptionRepository{db: db}
}

// Create inserts a new option into the database.
func (r *OptionRepository) Create(key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO options (key, value) VALUES (?, ?)",
		key,
		value,
	)
	return err
}

// Get retrieves an option by key.
func (r *OptionRepository) Get(key string) (*Option, error) {
	option := &Option{}
	err := r.db.QueryRow(
		`SELECT id, key, value, created_at, updated_at
		FROM options
		WHERE key = ?`,
		key,
	).Scan(
		&option.ID,
		&option.Key,
		&option.Value,
		&option.CreatedAt,
		&option.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return option, nil
}

// Update updates an option value.
func (r *OptionRepository) Update(key, value string) error {
	_, err := r.db.Exec(
		`UPDATE options SET
			value = ?, updated_at = ?
		WHERE key = ?`,
		value,
		time.Now().UTC(),
		key,
	)
	return err
}

// Delete removes an option from the database.
func (r *OptionRepository) Delete(key string) error {
	_, err := r.db.Exec("DELETE FROM options WHERE key = ?", key)
	return err
}

// List retrieves all options from the database.
func (r *OptionRepository) List() ([]*Option, error) {
	rows, err := r.db.Query("SELECT id, key, value, created_at, updated_at FROM options ORDER BY key")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []*Option
	for rows.Next() {
		option := &Option{}
		if err := rows.Scan(
			&option.ID,
			&option.Key,
			&option.Value,
			&option.CreatedAt,
			&option.UpdatedAt,
		); err != nil {
			return nil, err
		}
		options = append(options, option)
	}

	return options, rows.Err()
}
