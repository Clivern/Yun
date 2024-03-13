// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package discovery

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// StreamableHTTPClient implements Client using streamable HTTP communication
type StreamableHTTPClient struct {
	id              string
	url             string
	headers         map[string]string
	httpClient      *http.Client
	requestID       int
	requestIDMutex  sync.Mutex
	sessionID       string
	sessionMutex    sync.Mutex
	protocolVersion MCPProtocolVersion
	jsonRPCVersion  JSONRPCVersion
	clientInfo      ClientInfo
	serverInfo      *ServerInfo
	initialized     bool
	timeout         time.Duration
}

// StreamableHTTPClientConfig represents configuration for StreamableHTTPClient
type StreamableHTTPClientConfig struct {
	// ID is an optional database identifier for logging (e.g., "gateway:123")
	ID string

	// URL is the HTTP endpoint URL (e.g., "https://api.example.com/mcp")
	URL string

	// Headers for authentication/authorization
	// Common examples:
	//   - Authorization: "Bearer <token>"
	//   - X-API-Key: "<api-key>"
	//   - X-Custom-Auth: "<custom-value>"
	Headers map[string]string

	// ProtocolVersion is the MCP protocol version to use
	ProtocolVersion MCPProtocolVersion

	// JSONRPCVersion is the JSON-RPC version to use
	JSONRPCVersion JSONRPCVersion

	// ClientInfo contains client identification
	ClientInfo ClientInfo

	// Timeout for operations (default: 30 seconds)
	Timeout time.Duration
}

// NewStreamableHTTPClient creates a new streamable HTTP client
func NewStreamableHTTPClient(config StreamableHTTPClientConfig) (Client, error) {
	// Set defaults
	if config.ProtocolVersion == "" {
		config.ProtocolVersion = DefaultMCPVersion
	}
	if config.JSONRPCVersion == "" {
		config.JSONRPCVersion = JSONRPC20
	}
	if config.ClientInfo.Name == "" {
		config.ClientInfo.Name = "mut-client"
	}
	if config.ClientInfo.Version == "" {
		config.ClientInfo.Version = "0.1.0-dev"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.URL == "" {
		return nil, fmt.Errorf("HTTP URL is required")
	}

	// Initialize headers map if nil
	headers := config.Headers
	if headers == nil {
		headers = make(map[string]string)
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: config.Timeout,
	}

	return &StreamableHTTPClient{
		id:              config.ID,
		url:             config.URL,
		headers:         headers,
		httpClient:      httpClient,
		requestID:       0,
		protocolVersion: config.ProtocolVersion,
		jsonRPCVersion:  config.JSONRPCVersion,
		clientInfo:      config.ClientInfo,
		timeout:         config.Timeout,
	}, nil
}

// nextRequestID returns the next request ID
func (c *StreamableHTTPClient) nextRequestID() int {
	c.requestIDMutex.Lock()
	defer c.requestIDMutex.Unlock()
	c.requestID++
	return c.requestID
}

// sendRequest sends a JSON-RPC request via HTTP POST and returns the response
func (c *StreamableHTTPClient) sendRequest(ctx context.Context, method string, params map[string]interface{}) (*JSONRPCResponse, error) {
	reqID := c.nextRequestID()
	request := JSONRPCRequest{
		JSONRPC: string(c.jsonRPCVersion),
		ID:      &reqID,
		Method:  method,
		Params:  params,
	}

	// Marshal request
	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Info().
		Str("method", method).
		Int("id", reqID).
		Str("http_id", c.id).
		Msg("Sending MCP request via HTTP")

	// Create HTTP request with context timeout
	reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, c.url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Add session ID if we have one
	c.sessionMutex.Lock()
	if c.sessionID != "" {
		req.Header.Set("mcp-session-id", c.sessionID)
	}
	c.sessionMutex.Unlock()

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	// Read and unmarshal response
	// Check if response is SSE format
	contentType := resp.Header.Get("Content-Type")
	var jsonData []byte

	if strings.Contains(contentType, "text/event-stream") {
		// Parse SSE format: event: message\ndata: {...}\n
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				jsonData = []byte(strings.TrimPrefix(line, "data: "))
				break
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to parse SSE response: %w", err)
		}
		if jsonData == nil {
			return nil, fmt.Errorf("no data field found in SSE response")
		}
	} else {
		// Plain JSON response
		var err error
		jsonData, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
	}

	var response JSONRPCResponse
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Validate response ID matches request ID
	if response.ID != reqID {
		return nil, fmt.Errorf("response ID mismatch: expected %d, got %d", reqID, response.ID)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error %d: %s", response.Error.Code, response.Error.Message)
	}

	// Extract and store session ID if present
	if sessionID := resp.Header.Get("mcp-session-id"); sessionID != "" {
		c.sessionMutex.Lock()
		c.sessionID = sessionID
		c.sessionMutex.Unlock()
	}

	log.Info().
		Str("method", method).
		Int("id", reqID).
		Str("http_id", c.id).
		Msg("Received MCP response via HTTP")

	return &response, nil
}

// sendNotification sends a JSON-RPC notification (no response expected)
func (c *StreamableHTTPClient) sendNotification(ctx context.Context, method string, params map[string]interface{}) error {
	request := JSONRPCRequest{
		JSONRPC: string(c.jsonRPCVersion),
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	log.Info().
		Str("method", method).
		Str("http_id", c.id).
		Msg("Sending MCP notification via HTTP")

	// Create HTTP request with context timeout
	reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, c.url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Add session ID if we have one
	c.sessionMutex.Lock()
	if c.sessionID != "" {
		req.Header.Set("mcp-session-id", c.sessionID)
	}
	c.sessionMutex.Unlock()

	// Send request (we don't need to read the response for notifications)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code (202 Accepted is also a valid response for notifications)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Initialize initializes the MCP connection
func (c *StreamableHTTPClient) Initialize(ctx context.Context) (*InitializeResult, error) {
	if c.initialized {
		return nil, fmt.Errorf("client already initialized")
	}

	params := map[string]interface{}{
		"protocolVersion": string(c.protocolVersion),
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name":    c.clientInfo.Name,
			"version": c.clientInfo.Version,
		},
	}

	response, err := c.sendRequest(ctx, "initialize", params)
	if err != nil {
		return nil, fmt.Errorf("initialize request failed: %w", err)
	}

	result, err := ParseInitializeResponse(response.Result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse initialize response: %w", err)
	}

	c.serverInfo = &result.ServerInfo
	c.initialized = true

	// Send initialized notification
	if err := c.sendNotification(ctx, "notifications/initialized", nil); err != nil {
		return nil, fmt.Errorf("failed to send initialized notification: %w", err)
	}

	log.Info().
		Str("server", result.ServerInfo.Name).
		Str("version", result.ServerInfo.Version).
		Str("http_id", c.id).
		Msg("MCP client initialized via HTTP")

	return result, nil
}

// ListTools lists all available tools
func (c *StreamableHTTPClient) ListTools(ctx context.Context) ([]Tool, error) {
	if !c.initialized {
		return nil, fmt.Errorf("client not initialized")
	}

	response, err := c.sendRequest(ctx, "tools/list", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("tools/list request failed: %w", err)
	}

	return ParseToolsListResponse(response.Result)
}

// CallTool calls a tool with given arguments
func (c *StreamableHTTPClient) CallTool(ctx context.Context, name string, arguments ToolArgument) (*ToolCallResult, error) {
	if !c.initialized {
		return nil, fmt.Errorf("client not initialized")
	}

	params := map[string]interface{}{
		"name":      name,
		"arguments": arguments,
	}

	response, err := c.sendRequest(ctx, "tools/call", params)
	if err != nil {
		return nil, fmt.Errorf("tools/call request failed: %w", err)
	}

	// Parse result
	data, err := json.Marshal(response.Result)
	if err != nil {
		return nil, err
	}

	var result ToolCallResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListPrompts lists all available prompts
func (c *StreamableHTTPClient) ListPrompts(ctx context.Context) ([]Prompt, error) {
	if !c.initialized {
		return nil, fmt.Errorf("client not initialized")
	}

	response, err := c.sendRequest(ctx, "prompts/list", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("prompts/list request failed: %w", err)
	}

	return ParsePromptsListResponse(response.Result)
}

// GetPrompt gets a prompt with given arguments
func (c *StreamableHTTPClient) GetPrompt(ctx context.Context, name string, arguments map[string]string) (*PromptResult, error) {
	if !c.initialized {
		return nil, fmt.Errorf("client not initialized")
	}

	params := map[string]interface{}{
		"name":      name,
		"arguments": arguments,
	}

	response, err := c.sendRequest(ctx, "prompts/get", params)
	if err != nil {
		return nil, fmt.Errorf("prompts/get request failed: %w", err)
	}

	// Parse result
	data, err := json.Marshal(response.Result)
	if err != nil {
		return nil, err
	}

	var result PromptResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListResources lists all available resources
func (c *StreamableHTTPClient) ListResources(ctx context.Context) ([]Resource, error) {
	if !c.initialized {
		return nil, fmt.Errorf("client not initialized")
	}

	response, err := c.sendRequest(ctx, "resources/list", map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("resources/list request failed: %w", err)
	}

	return ParseResourcesListResponse(response.Result)
}

// ReadResource reads a resource by URI
func (c *StreamableHTTPClient) ReadResource(ctx context.Context, uri string) (*ResourceReadResult, error) {
	if !c.initialized {
		return nil, fmt.Errorf("client not initialized")
	}

	params := map[string]interface{}{
		"uri": uri,
	}

	response, err := c.sendRequest(ctx, "resources/read", params)
	if err != nil {
		return nil, fmt.Errorf("resources/read request failed: %w", err)
	}

	// Parse result
	data, err := json.Marshal(response.Result)
	if err != nil {
		return nil, err
	}

	var result ResourceReadResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Discover performs full discovery of server capabilities
func (c *StreamableHTTPClient) Discover(ctx context.Context) (*Result, error) {
	// Initialize if not already done
	if !c.initialized {
		initResult, err := c.Initialize(ctx)
		if err != nil {
			return nil, fmt.Errorf("initialization failed: %w", err)
		}
		c.serverInfo = &initResult.ServerInfo
	}

	result := &Result{
		ServerInfo: *c.serverInfo,
	}

	// Discover tools
	tools, err := c.ListTools(ctx)
	if err != nil {
		log.Warn().
			Err(err).
			Str("http_id", c.id).
			Msg("Failed to list tools")
	} else {
		result.Tools = tools
	}

	// Discover prompts
	prompts, err := c.ListPrompts(ctx)
	if err != nil {
		log.Warn().
			Err(err).
			Str("http_id", c.id).
			Msg("Failed to list prompts")
	} else {
		result.Prompts = prompts
	}

	// Discover resources
	resources, err := c.ListResources(ctx)
	if err != nil {
		log.Warn().
			Err(err).
			Str("http_id", c.id).
			Msg("Failed to list resources")
	} else {
		result.Resources = resources
	}

	log.Info().
		Int("tools", len(result.Tools)).
		Int("prompts", len(result.Prompts)).
		Int("resources", len(result.Resources)).
		Str("http_id", c.id).
		Msg("MCP discovery completed via HTTP")

	return result, nil
}

// Close closes the client connection
func (c *StreamableHTTPClient) Close() error {
	log.Info().
		Str("http_id", c.id).
		Msg("Closing MCP HTTP client")

	// HTTP client doesn't need explicit cleanup, but we reset the state
	c.initialized = false
	c.serverInfo = nil

	return nil
}
