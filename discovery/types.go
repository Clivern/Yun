// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package discovery

import (
	"encoding/json"
)

// MCPProtocolVersion represents the MCP protocol version
type MCPProtocolVersion string

// JSONRPCVersion represents the JSON-RPC protocol version
type JSONRPCVersion string

const (
	// JSONRPC20 represents JSON-RPC 2.0
	JSONRPC20 JSONRPCVersion = "2.0"

	// MCP Protocol Versions (date-based: YYYY-MM-DD)
	// Each version represents the last date when backward-incompatible changes were made

	// MCPVersion20241105 represents MCP protocol version 2024-11-05 (pre-release)
	MCPVersion20241105 MCPProtocolVersion = "2024-11-05"

	// MCPVersion20241125 represents MCP protocol version 2024-11-25 (initial official release)
	MCPVersion20241125 MCPProtocolVersion = "2024-11-25"

	// MCPVersion20250618 represents MCP protocol version 2025-06-18 (current stable)
	// Features: structured tool outputs, OAuth authorization, elicitation, enhanced security
	MCPVersion20250618 MCPProtocolVersion = "2025-06-18"

	// DefaultMCPVersion is the default protocol version to use
	// Use 2024-11-05 for backward compatibility with existing code
	DefaultMCPVersion = MCPVersion20241105

	// LatestMCPVersion is the latest supported protocol version
	LatestMCPVersion = MCPVersion20250618
)

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      *int                   `json:"id,omitempty"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      int                    `json:"id"`
	Result  map[string]interface{} `json:"result,omitempty"`
	Error   *JSONRPCError          `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// ClientInfo represents MCP client information
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ServerInfo represents MCP server information
type ServerInfo struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	ProtocolVersion string `json:"protocolVersion,omitempty"`
}

// InitializeParams represents parameters for initialize request
type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
}

// InitializeResult represents the result of initialize request
type InitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ServerInfo      ServerInfo             `json:"serverInfo"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ToolArgument represents arguments for a tool call
type ToolArgument map[string]interface{}

// ToolContent represents the content of a tool result
type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ToolCallResult represents the result of a tool call
type ToolCallResult struct {
	Content []ToolContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// Prompt represents an MCP prompt
type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
}

// PromptArgument represents an argument for a prompt
type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// PromptMessage represents a message in a prompt result
type PromptMessage struct {
	Role    string      `json:"role"`
	Content ToolContent `json:"content"`
}

// PromptResult represents the result of getting a prompt
type PromptResult struct {
	Description string          `json:"description"`
	Messages    []PromptMessage `json:"messages"`
}

// Resource represents an MCP resource
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

// ResourceContent represents the content of a resource
type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text,omitempty"`
	Blob     string `json:"blob,omitempty"`
}

// ResourceReadResult represents the result of reading a resource
type ResourceReadResult struct {
	Contents []ResourceContent `json:"contents"`
}

// Result contains all discovered MCP capabilities
type Result struct {
	ServerInfo ServerInfo `json:"serverInfo"`
	Tools      []Tool     `json:"tools"`
	Prompts    []Prompt   `json:"prompts"`
	Resources  []Resource `json:"resources"`
}

// ParseToolsListResponse parses a tools/list response
func ParseToolsListResponse(result map[string]interface{}) ([]Tool, error) {
	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var response struct {
		Tools []Tool `json:"tools"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Tools, nil
}

// ParsePromptsListResponse parses a prompts/list response
func ParsePromptsListResponse(result map[string]interface{}) ([]Prompt, error) {
	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var response struct {
		Prompts []Prompt `json:"prompts"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Prompts, nil
}

// ParseResourcesListResponse parses a resources/list response
func ParseResourcesListResponse(result map[string]interface{}) ([]Resource, error) {
	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var response struct {
		Resources []Resource `json:"resources"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Resources, nil
}

// ParseInitializeResponse parses an initialize response
func ParseInitializeResponse(result map[string]interface{}) (*InitializeResult, error) {
	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var initResult InitializeResult
	if err := json.Unmarshal(data, &initResult); err != nil {
		return nil, err
	}

	return &initResult, nil
}

// SupportedVersions returns a list of all supported MCP protocol versions in order (oldest to newest)
func SupportedVersions() []MCPProtocolVersion {
	return []MCPProtocolVersion{
		MCPVersion20241105,
		MCPVersion20241125,
		MCPVersion20250618,
	}
}

// IsVersionSupported checks if a protocol version is supported
func IsVersionSupported(version MCPProtocolVersion) bool {
	for _, v := range SupportedVersions() {
		if v == version {
			return true
		}
	}
	return false
}
