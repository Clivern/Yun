// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Server represents a virtual MCP server in the database.
type Server struct {
	ID   int64
	Name string
	// TODO: Add fields for server
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ServerRepository handles database operations for servers.
type ServerRepository struct {
	db *sql.DB
}

// NewServerRepository creates a new server repository.
func NewServerRepository(db *sql.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

// ServerMeta represents metadata associated with a server.
type ServerMeta struct {
	ID        int64
	Key       string
	Value     string
	ServerID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ServerMetaRepository handles database operations for server metadata.
type ServerMetaRepository struct {
	db *sql.DB
}

// NewServerMetaRepository creates a new server meta repository.
func NewServerMetaRepository(db *sql.DB) *ServerMetaRepository {
	return &ServerMetaRepository{db: db}
}

// Create inserts new metadata for a server.
func (r *ServerMetaRepository) Create(serverID int64, key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO servers_meta (server_id, key, value) VALUES (?, ?, ?)",
		serverID,
		key,
		value,
	)
	return err
}

// Get retrieves metadata for a server by key.
func (r *ServerMetaRepository) Get(serverID int64, key string) (*ServerMeta, error) {
	meta := &ServerMeta{}
	err := r.db.QueryRow(
		`SELECT id, key, value, server_id, created_at, updated_at
		FROM servers_meta
		WHERE server_id = ? AND key = ?`,
		serverID,
		key,
	).Scan(
		&meta.ID,
		&meta.Key,
		&meta.Value,
		&meta.ServerID,
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

// Update updates metadata for a server.
func (r *ServerMetaRepository) Update(serverID int64, key, value string) error {
	_, err := r.db.Exec(
		`UPDATE servers_meta SET
			value = ?, updated_at = ?
		WHERE server_id = ? AND key = ?`,
		value,
		time.Now().UTC(),
		serverID,
		key,
	)
	return err
}

// Delete removes metadata for a server.
func (r *ServerMetaRepository) Delete(serverID int64, key string) error {
	_, err := r.db.Exec(
		"DELETE FROM servers_meta WHERE server_id = ? AND key = ?",
		serverID,
		key,
	)
	return err
}

// ListByServer retrieves all metadata for a server.
func (r *ServerMetaRepository) ListByServer(serverID int64) ([]*ServerMeta, error) {
	rows, err := r.db.Query(
		`SELECT id, key, value, server_id, created_at, updated_at
		FROM servers_meta
		WHERE server_id = ?
		ORDER BY key`,
		serverID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*ServerMeta
	for rows.Next() {
		meta := &ServerMeta{}
		if err := rows.Scan(
			&meta.ID,
			&meta.Key,
			&meta.Value,
			&meta.ServerID,
			&meta.CreatedAt,
			&meta.UpdatedAt,
		); err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}

	return metadata, rows.Err()
}

// Upsert inserts or updates metadata for a server.
func (r *ServerMetaRepository) Upsert(serverID int64, key, value string) error {
	existing, err := r.Get(serverID, key)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(serverID, key, value)
	}

	return r.Update(serverID, key, value)
}

// ServerTool represents a many-to-many relationship between servers and tools.
type ServerTool struct {
	ServerID  int64
	ToolID    int64
	CreatedAt time.Time
}

// ServerToolRepository handles the server-tool relationship operations.
type ServerToolRepository struct {
	db *sql.DB
}

// NewServerToolRepository creates a new server-tool repository.
func NewServerToolRepository(db *sql.DB) *ServerToolRepository {
	return &ServerToolRepository{db: db}
}

// AddTool associates a tool with a server.
func (r *ServerToolRepository) AddTool(serverID, toolID int64) error {
	_, err := r.db.Exec(
		"INSERT INTO server_tools (server_id, tool_id) VALUES (?, ?)",
		serverID,
		toolID,
	)
	return err
}

// RemoveTool removes a tool association from a server.
func (r *ServerToolRepository) RemoveTool(serverID, toolID int64) error {
	_, err := r.db.Exec(
		"DELETE FROM server_tools WHERE server_id = ? AND tool_id = ?",
		serverID,
		toolID,
	)
	return err
}

// GetToolsByServer retrieves all tool IDs associated with a server.
func (r *ServerToolRepository) GetToolsByServer(serverID int64) ([]int64, error) {
	rows, err := r.db.Query(
		"SELECT tool_id FROM server_tools WHERE server_id = ?",
		serverID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var toolIDs []int64
	for rows.Next() {
		var toolID int64
		if err := rows.Scan(
			&toolID,
		); err != nil {
			return nil, err
		}
		toolIDs = append(toolIDs, toolID)
	}

	return toolIDs, rows.Err()
}

// GetServersByTool retrieves all server IDs that use a specific tool.
func (r *ServerToolRepository) GetServersByTool(toolID int64) ([]int64, error) {
	rows, err := r.db.Query(
		"SELECT server_id FROM server_tools WHERE tool_id = ?",
		toolID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var serverIDs []int64
	for rows.Next() {
		var serverID int64
		if err := rows.Scan(
			&serverID,
		); err != nil {
			return nil, err
		}
		serverIDs = append(serverIDs, serverID)
	}

	return serverIDs, rows.Err()
}

// RemoveAllToolsFromServer removes all tool associations from a server.
func (r *ServerToolRepository) RemoveAllToolsFromServer(serverID int64) error {
	_, err := r.db.Exec(
		"DELETE FROM server_tools WHERE server_id = ?",
		serverID,
	)
	return err
}

// ServerResource represents a many-to-many relationship between servers and resources.
type ServerResource struct {
	ServerID   int64
	ResourceID int64
	CreatedAt  time.Time
}

// ServerResourceRepository handles the server-resource relationship operations.
type ServerResourceRepository struct {
	db *sql.DB
}

// NewServerResourceRepository creates a new server-resource repository.
func NewServerResourceRepository(db *sql.DB) *ServerResourceRepository {
	return &ServerResourceRepository{db: db}
}

// AddResource associates a resource with a server.
func (r *ServerResourceRepository) AddResource(serverID, resourceID int64) error {
	_, err := r.db.Exec(
		"INSERT INTO server_resources (server_id, resource_id) VALUES (?, ?)",
		serverID,
		resourceID,
	)
	return err
}

// RemoveResource removes a resource association from a server.
func (r *ServerResourceRepository) RemoveResource(serverID, resourceID int64) error {
	_, err := r.db.Exec(
		"DELETE FROM server_resources WHERE server_id = ? AND resource_id = ?",
		serverID,
		resourceID,
	)
	return err
}

// GetResourcesByServer retrieves all resource IDs associated with a server.
func (r *ServerResourceRepository) GetResourcesByServer(serverID int64) ([]int64, error) {
	rows, err := r.db.Query(
		"SELECT resource_id FROM server_resources WHERE server_id = ?",
		serverID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resourceIDs []int64
	for rows.Next() {
		var resourceID int64
		if err := rows.Scan(
			&resourceID,
		); err != nil {
			return nil, err
		}
		resourceIDs = append(resourceIDs, resourceID)
	}

	return resourceIDs, rows.Err()
}

// GetServersByResource retrieves all server IDs that use a specific resource.
func (r *ServerResourceRepository) GetServersByResource(resourceID int64) ([]int64, error) {
	rows, err := r.db.Query(
		"SELECT server_id FROM server_resources WHERE resource_id = ?",
		resourceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var serverIDs []int64
	for rows.Next() {
		var serverID int64
		if err := rows.Scan(
			&serverID,
		); err != nil {
			return nil, err
		}
		serverIDs = append(serverIDs, serverID)
	}

	return serverIDs, rows.Err()
}

// RemoveAllResourcesFromServer removes all resource associations from a server.
func (r *ServerResourceRepository) RemoveAllResourcesFromServer(serverID int64) error {
	_, err := r.db.Exec(
		"DELETE FROM server_resources WHERE server_id = ?",
		serverID,
	)
	return err
}

// ServerPrompt represents a many-to-many relationship between servers and prompts.
type ServerPrompt struct {
	ServerID  int64
	PromptID  int64
	CreatedAt time.Time
}

// ServerPromptRepository handles the server-prompt relationship operations.
type ServerPromptRepository struct {
	db *sql.DB
}

// NewServerPromptRepository creates a new server-prompt repository.
func NewServerPromptRepository(db *sql.DB) *ServerPromptRepository {
	return &ServerPromptRepository{db: db}
}

// AddPrompt associates a prompt with a server.
func (r *ServerPromptRepository) AddPrompt(serverID, promptID int64) error {
	_, err := r.db.Exec(
		"INSERT INTO server_prompts (server_id, prompt_id) VALUES (?, ?)",
		serverID,
		promptID,
	)
	return err
}

// RemovePrompt removes a prompt association from a server.
func (r *ServerPromptRepository) RemovePrompt(serverID, promptID int64) error {
	_, err := r.db.Exec(
		"DELETE FROM server_prompts WHERE server_id = ? AND prompt_id = ?",
		serverID,
		promptID,
	)
	return err
}

// GetPromptsByServer retrieves all prompt IDs associated with a server.
func (r *ServerPromptRepository) GetPromptsByServer(serverID int64) ([]int64, error) {
	rows, err := r.db.Query(
		"SELECT prompt_id FROM server_prompts WHERE server_id = ?",
		serverID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var promptIDs []int64
	for rows.Next() {
		var promptID int64
		if err := rows.Scan(
			&promptID,
		); err != nil {
			return nil, err
		}
		promptIDs = append(promptIDs, promptID)
	}

	return promptIDs, rows.Err()
}

// GetServersByPrompt retrieves all server IDs that use a specific prompt.
func (r *ServerPromptRepository) GetServersByPrompt(promptID int64) ([]int64, error) {
	rows, err := r.db.Query(
		"SELECT server_id FROM server_prompts WHERE prompt_id = ?",
		promptID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var serverIDs []int64
	for rows.Next() {
		var serverID int64
		if err := rows.Scan(
			&serverID,
		); err != nil {
			return nil, err
		}
		serverIDs = append(serverIDs, serverID)
	}

	return serverIDs, rows.Err()
}

// RemoveAllPromptsFromServer removes all prompt associations from a server.
func (r *ServerPromptRepository) RemoveAllPromptsFromServer(serverID int64) error {
	_, err := r.db.Exec(
		"DELETE FROM server_prompts WHERE server_id = ?",
		serverID,
	)
	return err
}
