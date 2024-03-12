// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package discovery

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// StdioClient implements Client using stdio communication
type StdioClient struct {
	id              string
	cmd             *exec.Cmd
	stdin           io.WriteCloser
	stdout          io.ReadCloser
	stderr          io.ReadCloser
	scanner         *bufio.Scanner
	requestID       int
	requestIDMutex  sync.Mutex
	protocolVersion MCPProtocolVersion
	jsonRPCVersion  JSONRPCVersion
	clientInfo      ClientInfo
	serverInfo      *ServerInfo
	initialized     bool
	timeout         time.Duration
}

// StdioClientConfig represents configuration for StdioClient
type StdioClientConfig struct {
	// ID is an optional database identifier for logging (e.g., "gateway:123" or "mcp:456")
	ID string

	// Command is the command to execute (e.g., "uv")
	Command string

	// Args are the command arguments (e.g., ["run", "python", "-m", "sample_mcp_server.server"])
	Args []string

	// WorkingDir is the working directory for the command
	WorkingDir string

	// ProtocolVersion is the MCP protocol version to use
	ProtocolVersion MCPProtocolVersion

	// JSONRPCVersion is the JSON-RPC version to use
	JSONRPCVersion JSONRPCVersion

	// ClientInfo contains client identification
	ClientInfo ClientInfo

	// Timeout for operations (default: 30 seconds)
	Timeout time.Duration
}

// NewStdioClient creates a new stdio client
func NewStdioClient(config StdioClientConfig) (Client, error) {
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

	// Create command
	cmd := exec.Command(config.Command, config.Args...)
	if config.WorkingDir != "" {
		cmd.Dir = config.WorkingDir
	}

	// Setup pipes
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	log.Info().
		Str("command", config.Command).
		Strs("args", config.Args).
		Str("stdio_id", config.ID).
		Msg("Started MCP server process")

	// Create scanner for reading responses
	scanner := bufio.NewScanner(stdout)

	// Start goroutine to read stderr
	go func() {
		stderrScanner := bufio.NewScanner(stderr)
		for stderrScanner.Scan() {
			log.Info().
				Str("error", stderrScanner.Text()).
				Str("stdio_id", config.ID).
				Msg("MCP server error")
		}
	}()

	return &StdioClient{
		cmd:             cmd,
		stdin:           stdin,
		stdout:          stdout,
		stderr:          stderr,
		scanner:         scanner,
		requestID:       0,
		protocolVersion: config.ProtocolVersion,
		jsonRPCVersion:  config.JSONRPCVersion,
		clientInfo:      config.ClientInfo,
		initialized:     false,
		timeout:         config.Timeout,
		id:              config.ID,
	}, nil
}

// nextRequestID returns the next request ID
func (c *StdioClient) nextRequestID() int {
	c.requestIDMutex.Lock()
	defer c.requestIDMutex.Unlock()
	c.requestID++
	return c.requestID
}

// sendRequest sends a JSON-RPC request and returns the response
func (c *StdioClient) sendRequest(ctx context.Context, method string, params map[string]interface{}) (*JSONRPCResponse, error) {
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
		Str("stdio_id", c.id).
		Msg("Sending MCP request")

	// Write request
	if _, err := c.stdin.Write(append(data, '\n')); err != nil {
		return nil, fmt.Errorf("failed to write request: %w", err)
	}

	// Read response with timeout
	responseChan := make(chan *JSONRPCResponse, 1)
	errorChan := make(chan error, 1)

	go func() {
		if c.scanner.Scan() {
			var response JSONRPCResponse
			if err := json.Unmarshal(c.scanner.Bytes(), &response); err != nil {
				errorChan <- fmt.Errorf("failed to unmarshal response: %w", err)
				return
			}
			responseChan <- &response
		} else {
			if err := c.scanner.Err(); err != nil {
				errorChan <- fmt.Errorf("scanner error: %w", err)
			} else {
				errorChan <- fmt.Errorf("no response received")
			}
		}
	}()

	select {
	case response := <-responseChan:
		// Validate response ID matches request ID
		if response.ID != reqID {
			return nil, fmt.Errorf("response ID mismatch: expected %d, got %d", reqID, response.ID)
		}
		if response.Error != nil {
			return nil, fmt.Errorf("JSON-RPC error %d: %s", response.Error.Code, response.Error.Message)
		}
		log.Info().
			Str("method", method).
			Int("id", reqID).
			Str("stdio_id", c.id).
			Msg("Received MCP response")
		return response, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("request timeout: %w", ctx.Err())
	}
}

// sendNotification sends a JSON-RPC notification (no response expected)
func (c *StdioClient) sendNotification(method string, params map[string]interface{}) error {
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
		Str("stdio_id", c.id).
		Msg("Sending MCP notification")

	if _, err := c.stdin.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write notification: %w", err)
	}

	return nil
}

// Initialize initializes the MCP connection
func (c *StdioClient) Initialize(ctx context.Context) (*InitializeResult, error) {
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
	if err := c.sendNotification("notifications/initialized", nil); err != nil {
		return nil, fmt.Errorf("failed to send initialized notification: %w", err)
	}

	log.Info().
		Str("server", result.ServerInfo.Name).
		Str("version", result.ServerInfo.Version).
		Str("stdio_id", c.id).
		Msg("MCP client initialized")

	return result, nil
}

// ListTools lists all available tools
func (c *StdioClient) ListTools(ctx context.Context) ([]Tool, error) {
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
func (c *StdioClient) CallTool(ctx context.Context, name string, arguments ToolArgument) (*ToolCallResult, error) {
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
func (c *StdioClient) ListPrompts(ctx context.Context) ([]Prompt, error) {
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
func (c *StdioClient) GetPrompt(ctx context.Context, name string, arguments map[string]string) (*PromptResult, error) {
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
func (c *StdioClient) ListResources(ctx context.Context) ([]Resource, error) {
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
func (c *StdioClient) ReadResource(ctx context.Context, uri string) (*ResourceReadResult, error) {
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
func (c *StdioClient) Discover(ctx context.Context) (*Result, error) {
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
			Str("stdio_id", c.id).
			Msg("Failed to list tools")
	} else {
		result.Tools = tools
	}

	// Discover prompts
	prompts, err := c.ListPrompts(ctx)
	if err != nil {
		log.Warn().
			Err(err).
			Str("stdio_id", c.id).
			Msg("Failed to list prompts")
	} else {
		result.Prompts = prompts
	}

	// Discover resources
	resources, err := c.ListResources(ctx)
	if err != nil {
		log.Warn().
			Err(err).
			Str("stdio_id", c.id).
			Msg("Failed to list resources")
	} else {
		result.Resources = resources
	}

	log.Info().
		Int("tools", len(result.Tools)).
		Int("prompts", len(result.Prompts)).
		Int("resources", len(result.Resources)).
		Str("stdio_id", c.id).
		Msg("MCP discovery completed")

	return result, nil
}

// Close closes the client connection
func (c *StdioClient) Close() error {
	log.Info().
		Str("stdio_id", c.id).
		Msg("Closing MCP client")

	// Close stdin to signal the process to exit
	if c.stdin != nil {
		c.stdin.Close()
	}

	// Kill the process and wait for it to exit
	if c.cmd != nil && c.cmd.Process != nil {
		// Try to kill the process gracefully first
		if err := c.cmd.Process.Kill(); err != nil {
			log.Debug().
				Err(err).
				Str("stdio_id", c.id).
				Msg("Failed to kill process")
		}

		// Wait for process to exit (with timeout handled by test context)
		if err := c.cmd.Wait(); err != nil {
			// Process.Kill() will cause Wait() to return an error, which is expected
			log.Debug().
				Err(err).
				Str("stdio_id", c.id).
				Msg("Process exited")
		}
	}

	return nil
}
