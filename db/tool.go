// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Tool represents a tool discovered from MCP servers.
type Tool struct {
	ID   int64
	Name string
	// TODO: Add fields for tool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToolRepository handles database operations for tools.
type ToolRepository struct {
	db *sql.DB
}

// NewToolRepository creates a new tool repository.
func NewToolRepository(db *sql.DB) *ToolRepository {
	return &ToolRepository{db: db}
}

// ToolMeta represents metadata associated with a tool.
type ToolMeta struct {
	ID        int64
	Key       string
	Value     string
	ToolID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToolMetaRepository handles database operations for tool metadata.
type ToolMetaRepository struct {
	db *sql.DB
}

// NewToolMetaRepository creates a new tool meta repository.
func NewToolMetaRepository(db *sql.DB) *ToolMetaRepository {
	return &ToolMetaRepository{db: db}
}

// Create inserts new metadata for a tool.
func (r *ToolMetaRepository) Create(toolID int64, key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO tools_meta (tool_id, key, value) VALUES (?, ?, ?)",
		toolID,
		key,
		value,
	)
	return err
}

// Get retrieves metadata for a tool by key.
func (r *ToolMetaRepository) Get(toolID int64, key string) (*ToolMeta, error) {
	meta := &ToolMeta{}
	err := r.db.QueryRow(
		`SELECT id, key, value, tool_id, created_at, updated_at
		FROM tools_meta
		WHERE tool_id = ? AND key = ?`,
		toolID,
		key,
	).Scan(
		&meta.ID,
		&meta.Key,
		&meta.Value,
		&meta.ToolID,
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

// Update updates metadata for a tool.
func (r *ToolMetaRepository) Update(toolID int64, key, value string) error {
	_, err := r.db.Exec(
		"UPDATE tools_meta SET value = ?, updated_at = ? WHERE tool_id = ? AND key = ?",
		value,
		time.Now().UTC(),
		toolID,
		key,
	)
	return err
}

// Delete removes metadata for a tool.
func (r *ToolMetaRepository) Delete(toolID int64, key string) error {
	_, err := r.db.Exec(
		"DELETE FROM tools_meta WHERE tool_id = ? AND key = ?",
		toolID,
		key,
	)
	return err
}

// ListByTool retrieves all metadata for a tool.
func (r *ToolMetaRepository) ListByTool(toolID int64) ([]*ToolMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, key, value, tool_id, created_at, updated_at
		FROM tools_meta
		WHERE tool_id = ?
		ORDER BY key`,
		toolID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*ToolMeta
	for rows.Next() {
		meta := &ToolMeta{}
		if err := rows.Scan(
			&meta.ID,
			&meta.Key,
			&meta.Value,
			&meta.ToolID,
			&meta.CreatedAt,
			&meta.UpdatedAt,
		); err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}

	return metadata, rows.Err()
}

// Upsert inserts or updates metadata for a tool.
func (r *ToolMetaRepository) Upsert(toolID int64, key, value string) error {
	existing, err := r.Get(toolID, key)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(toolID, key, value)
	}

	return r.Update(toolID, key, value)
}
