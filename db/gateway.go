// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Gateway represents a gateway in the database.
type Gateway struct {
	ID   int64
	Name string
	// TODO: Add fields for gateway
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GatewayRepository handles database operations for gateways.
type GatewayRepository struct {
	db *sql.DB
}

// NewGatewayRepository creates a new gateway repository.
func NewGatewayRepository(db *sql.DB) *GatewayRepository {
	return &GatewayRepository{db: db}
}

// Count returns the total number of gateways.
func (r *GatewayRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM gateways").Scan(&count)
	return count, err
}

// GatewayMeta represents metadata associated with a gateway.
type GatewayMeta struct {
	ID        int64
	Key       string
	Value     string
	GatewayID int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GatewayMetaRepository handles database operations for gateway metadata.
type GatewayMetaRepository struct {
	db *sql.DB
}

// NewGatewayMetaRepository creates a new gateway meta repository.
func NewGatewayMetaRepository(db *sql.DB) *GatewayMetaRepository {
	return &GatewayMetaRepository{db: db}
}

// Create inserts new metadata for a gateway.
func (r *GatewayMetaRepository) Create(gatewayID int64, key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO gateways_meta (gateway_id, key, value) VALUES (?, ?, ?)",
		gatewayID,
		key,
		value,
	)
	return err
}

// Get retrieves metadata for a gateway by key.
func (r *GatewayMetaRepository) Get(gatewayID int64, key string) (*GatewayMeta, error) {
	meta := &GatewayMeta{}
	err := r.db.QueryRow(
		`SELECT id, key, value, gateway_id, created_at, updated_at
		FROM gateways_meta
		WHERE gateway_id = ? AND key = ?`,
		gatewayID,
		key,
	).Scan(
		&meta.ID,
		&meta.Key,
		&meta.Value,
		&meta.GatewayID,
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

// Update updates metadata for a gateway.
func (r *GatewayMetaRepository) Update(gatewayID int64, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE gateways_meta SET
			value = ?, updated_at = ?
		WHERE gateway_id = ? AND key = ?`,
		value,
		time.Now().UTC(),
		gatewayID,
		key,
	)
	return err
}

// Delete removes metadata for a gateway.
func (r *GatewayMetaRepository) Delete(gatewayID int64, key string) error {
	_, err := r.db.Exec(
		"DELETE FROM gateways_meta WHERE gateway_id = ? AND key = ?",
		gatewayID,
		key,
	)
	return err
}

// ListByGateway retrieves all metadata for a gateway.
func (r *GatewayMetaRepository) ListByGateway(gatewayID int64) ([]*GatewayMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, key, value, gateway_id, created_at, updated_at
		FROM gateways_meta
		WHERE gateway_id = ?
		ORDER BY key`,
		gatewayID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*GatewayMeta
	for rows.Next() {
		meta := &GatewayMeta{}
		if err := rows.Scan(
			&meta.ID,
			&meta.Key,
			&meta.Value,
			&meta.GatewayID,
			&meta.CreatedAt,
			&meta.UpdatedAt,
		); err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}

	return metadata, rows.Err()
}

// Upsert inserts or updates metadata for a gateway.
func (r *GatewayMetaRepository) Upsert(gatewayID int64, key, value string) error {
	existing, err := r.Get(gatewayID, key)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(gatewayID, key, value)
	}

	return r.Update(gatewayID, key, value)
}
