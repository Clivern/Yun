// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Prompt represents a prompt discovered from MCP servers.
type Prompt struct {
	ID   int64
	Name string
	// TODO: Add fields for prompt
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PromptRepository handles database operations for prompts.
type PromptRepository struct {
	db *sql.DB
}

// NewPromptRepository creates a new prompt repository.
func NewPromptRepository(db *sql.DB) *PromptRepository {
	return &PromptRepository{db: db}
}

// PromptMeta represents metadata associated with a prompt.
type PromptMeta struct {
	ID        int64
	Key       string
	Value     string
	PromptID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PromptMetaRepository handles database operations for prompt metadata.
type PromptMetaRepository struct {
	db *sql.DB
}

// NewPromptMetaRepository creates a new prompt meta repository.
func NewPromptMetaRepository(db *sql.DB) *PromptMetaRepository {
	return &PromptMetaRepository{db: db}
}

// Create inserts new metadata for a prompt.
func (r *PromptMetaRepository) Create(promptID int64, key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO prompts_meta (prompt_id, key, value) VALUES (?, ?, ?)",
		promptID,
		key,
		value,
	)
	return err
}

// Get retrieves metadata for a prompt by key.
func (r *PromptMetaRepository) Get(promptID int64, key string) (*PromptMeta, error) {
	meta := &PromptMeta{}
	err := r.db.QueryRow(
		`SELECT id, key, value, prompt_id, created_at, updated_at
		FROM prompts_meta
		WHERE prompt_id = ? AND key = ?`,
		promptID,
		key,
	).Scan(
		&meta.ID,
		&meta.Key,
		&meta.Value,
		&meta.PromptID,
		&meta.CreatedAt,
		&meta.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return meta, nil
}

// Update updates metadata for a prompt.
func (r *PromptMetaRepository) Update(promptID int64, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE prompts_meta SET
			value = ?, updated_at = ?
		WHERE prompt_id = ? AND key = ?`,
		value,
		time.Now().UTC(),
		promptID,
		key,
	)
	return err
}

// Delete removes metadata for a prompt.
func (r *PromptMetaRepository) Delete(promptID int64, key string) error {
	_, err := r.db.Exec(
		"DELETE FROM prompts_meta WHERE prompt_id = ? AND key = ?",
		promptID,
		key,
	)
	return err
}

// ListByPrompt retrieves all metadata for a prompt.
func (r *PromptMetaRepository) ListByPrompt(promptID int64) ([]*PromptMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, key, value, prompt_id, created_at, updated_at
		FROM prompts_meta
		WHERE prompt_id = ?
		ORDER BY key`,
		promptID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*PromptMeta
	for rows.Next() {
		meta := &PromptMeta{}
		if err := rows.Scan(
			&meta.ID,
			&meta.Key,
			&meta.Value,
			&meta.PromptID,
			&meta.CreatedAt,
			&meta.UpdatedAt,
		); err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}

	return metadata, rows.Err()
}

// Upsert inserts or updates metadata for a prompt.
func (r *PromptMetaRepository) Upsert(promptID int64, key, value string) error {
	existing, err := r.Get(promptID, key)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(promptID, key, value)
	}

	return r.Update(promptID, key, value)
}
