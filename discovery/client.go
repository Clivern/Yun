// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package discovery provides Gateway discovery functions for the Mut application.
package discovery

import (
	"context"
)

// Client represents a discovery client for MCP servers.
//
// The Client interface supports multiple transport mechanisms:
//
// 1. Stdio Transport (StdioClient):
//   - Communicates via stdin/stdout with a subprocess
//   - Security is handled at the OS process level
//   - Headers are ignored (no protocol-level auth)
//
// 2. SSE/HTTP Transport (SSEClient):
//   - Communicates via Server-Sent Events over HTTP
//   - Supports authentication via HTTP headers
//   - Headers can include: Authorization, X-API-Key, etc.
//
// Example usage with stdio (no auth):
//
//	client, err := NewStdioClient(StdioClientConfig{
//	    Command: "python",
//	    Args: []string{"-m", "mcp_server"},
//	})
//
// Example usage with SSE (with auth):
//
//	client, err := NewSSEClient(SSEClientConfig{
//	    URL: "https://api.example.com/mcp/sse",
//	    Headers: map[string]string{
//	        "Authorization": "Bearer token123",
//	        "X-API-Key": "key456",
//	    },
//	})
type Client interface {
	// Initialize initializes the MCP connection
	Initialize(ctx context.Context) (*InitializeResult, error)

	// ListTools lists all available tools
	ListTools(ctx context.Context) ([]Tool, error)

	// CallTool calls a tool with given arguments
	CallTool(ctx context.Context, name string, arguments ToolArgument) (*ToolCallResult, error)

	// ListPrompts lists all available prompts
	ListPrompts(ctx context.Context) ([]Prompt, error)

	// GetPrompt gets a prompt with given arguments
	GetPrompt(ctx context.Context, name string, arguments map[string]string) (*PromptResult, error)

	// ListResources lists all available resources
	ListResources(ctx context.Context) ([]Resource, error)

	// ReadResource reads a resource by URI
	ReadResource(ctx context.Context, uri string) (*ResourceReadResult, error)

	// Discover performs full discovery of server capabilities
	Discover(ctx context.Context) (*Result, error)

	// Close closes the client connection
	Close() error
}
