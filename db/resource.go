// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Resource represents a resource discovered from MCP servers.
type Resource struct {
	ID   int64
	Name string
	// TODO: Add fields for resource
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ResourceRepository handles database operations for resources.
type ResourceRepository struct {
	db *sql.DB
}

// NewResourceRepository creates a new resource repository.
func NewResourceRepository(db *sql.DB) *ResourceRepository {
	return &ResourceRepository{db: db}
}

// ResourceMeta represents metadata associated with a resource.
type ResourceMeta struct {
	ID         int64
	Key        string
	Value      string
	ResourceID int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ResourceMetaRepository handles database operations for resource metadata.
type ResourceMetaRepository struct {
	db *sql.DB
}

// NewResourceMetaRepository creates a new resource meta repository.
func NewResourceMetaRepository(db *sql.DB) *ResourceMetaRepository {
	return &ResourceMetaRepository{db: db}
}

// Create inserts new metadata for a resource.
func (r *ResourceMetaRepository) Create(resourceID int64, key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO resources_meta (resource_id, key, value) VALUES (?, ?, ?)",
		resourceID,
		key,
		value,
	)
	return err
}

// Get retrieves metadata for a resource by key.
func (r *ResourceMetaRepository) Get(resourceID int64, key string) (*ResourceMeta, error) {
	meta := &ResourceMeta{}
	err := r.db.QueryRow(
		`SELECT id, key, value, resource_id, created_at, updated_at
		FROM resources_meta
		WHERE resource_id = ? AND key = ?`,
		resourceID,
		key,
	).Scan(
		&meta.ID,
		&meta.Key,
		&meta.Value,
		&meta.ResourceID,
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

// Update updates metadata for a resource.
func (r *ResourceMetaRepository) Update(resourceID int64, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE resources_meta SET
			value = ?, updated_at = ?
		WHERE resource_id = ? AND key = ?`,
		value,
		time.Now().UTC(),
		resourceID,
		key,
	)
	return err
}

// Delete removes metadata for a resource.
func (r *ResourceMetaRepository) Delete(resourceID int64, key string) error {
	_, err := r.db.Exec(
		"DELETE FROM resources_meta WHERE resource_id = ? AND key = ?",
		resourceID,
		key,
	)
	return err
}

// ListByResource retrieves all metadata for a resource.
func (r *ResourceMetaRepository) ListByResource(resourceID int64) ([]*ResourceMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, key, value, resource_id, created_at, updated_at
		FROM resources_meta
		WHERE resource_id = ?
		ORDER BY key`,
		resourceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*ResourceMeta
	for rows.Next() {
		meta := &ResourceMeta{}
		if err := rows.Scan(
			&meta.ID,
			&meta.Key,
			&meta.Value,
			&meta.ResourceID,
			&meta.CreatedAt,
			&meta.UpdatedAt,
		); err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}

	return metadata, rows.Err()
}

// Upsert inserts or updates metadata for a resource.
func (r *ResourceMetaRepository) Upsert(resourceID int64, key, value string) error {
	existing, err := r.Get(resourceID, key)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(resourceID, key, value)
	}

	return r.Update(resourceID, key, value)
}
