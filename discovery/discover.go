// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package discovery provides Gateway discovery functions for the Mut application.
package discovery

import (
	"context"

	"github.com/rs/zerolog/log"
)

// DiscoverStdio discovers the MCP server using stdio
func DiscoverStdio(ctx context.Context, config StdioClientConfig) (*Result, error) {
	client, err := NewStdioClient(config)
	if err != nil {
		log.Error().
			Err(err).
			Str("stdio_id", config.ID).
			Msg("Failed to create stdio client")
		return nil, err
	}

	log.Info().
		Str("stdio_id", config.ID).
		Msg("Created stdio client")

	result, err := client.Discover(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Str("stdio_id", config.ID).
			Msg("Failed to discover MCP server")
		return nil, err
	}

	log.Info().
		Str("stdio_id", config.ID).
		Msg("Discovered MCP server")

	return result, nil
}

// DiscoverSSE discovers the MCP server using SSE (Server-Sent Events)
func DiscoverSSE(ctx context.Context, config SSEClientConfig) (*Result, error) {
	client, err := NewSSEClient(config)
	if err != nil {
		log.Error().
			Err(err).
			Str("sse_id", config.ID).
			Str("url", config.URL).
			Msg("Failed to create SSE client")
		return nil, err
	}

	log.Info().
		Str("sse_id", config.ID).
		Str("url", config.URL).
		Msg("Created SSE client")

	result, err := client.Discover(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Str("sse_id", config.ID).
			Str("url", config.URL).
			Msg("Failed to discover MCP server")
		return nil, err
	}

	log.Info().
		Str("sse_id", config.ID).
		Str("url", config.URL).
		Msg("Discovered MCP server")

	return result, nil
}
