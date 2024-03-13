// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package discovery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockMCPServer creates an HTTP server that responds to MCP JSON-RPC requests
func mockMCPServer(_ *testing.T) *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Check Content-Type
		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Parse request
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Handle different methods
		var response JSONRPCResponse
		response.JSONRPC = "2.0"

		switch req.Method {
		case "initialize":
			if req.ID == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			response.ID = *req.ID
			response.Result = map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]interface{}{},
				"serverInfo": map[string]interface{}{
					"name":    "mock-server",
					"version": "1.0.0",
				},
			}

		case "notifications/initialized":
			// Notifications don't have IDs and don't return responses
			w.WriteHeader(http.StatusNoContent)
			return

		case "tools/list":
			if req.ID == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			response.ID = *req.ID
			response.Result = map[string]interface{}{
				"tools": []map[string]interface{}{
					{
						"name":        "test_tool",
						"description": "A test tool",
						"inputSchema": map[string]interface{}{
							"type": "object",
						},
					},
				},
			}

		case "tools/call":
			if req.ID == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			response.ID = *req.ID
			response.Result = map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": "Tool result: " + req.Params["name"].(string),
					},
				},
				"isError": false,
			}

		case "prompts/list":
			if req.ID == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			response.ID = *req.ID
			response.Result = map[string]interface{}{
				"prompts": []map[string]interface{}{
					{
						"name":        "test_prompt",
						"description": "A test prompt",
					},
				},
			}

		case "prompts/get":
			if req.ID == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			response.ID = *req.ID
			response.Result = map[string]interface{}{
				"description": "A test prompt",
				"messages": []map[string]interface{}{
					{
						"role": "user",
						"content": map[string]interface{}{
							"type": "text",
							"text": "Test prompt message",
						},
					},
				},
			}

		case "resources/list":
			if req.ID == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			response.ID = *req.ID
			response.Result = map[string]interface{}{
				"resources": []map[string]interface{}{
					{
						"uri":         "test://resource",
						"name":        "test_resource",
						"description": "A test resource",
						"mimeType":    "text/plain",
					},
				},
			}

		case "resources/read":
			if req.ID == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			response.ID = *req.ID
			response.Result = map[string]interface{}{
				"contents": []map[string]interface{}{
					{
						"uri":      "test://resource",
						"mimeType": "text/plain",
						"text":     "Test resource content",
					},
				},
			}

		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Send JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	return httptest.NewServer(http.HandlerFunc(handler))
}

// mockSSEMCPServer creates an HTTP server that responds with SSE format
func mockSSEMCPServer(_ *testing.T) *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Parse request
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Create response
		var response JSONRPCResponse
		response.JSONRPC = "2.0"

		switch req.Method {
		case "initialize":
			if req.ID == nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			response.ID = *req.ID
			response.Result = map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]interface{}{},
				"serverInfo": map[string]interface{}{
					"name":    "mock-sse-server",
					"version": "1.0.0",
				},
			}

		case "notifications/initialized":
			w.WriteHeader(http.StatusNoContent)
			return

		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Send SSE response
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("event: message\n"))
		data, _ := json.Marshal(response)
		w.Write([]byte("data: "))
		w.Write(data)
		w.Write([]byte("\n\n"))
	}

	return httptest.NewServer(http.HandlerFunc(handler))
}

// TestNewStreamableHTTPClient tests the creation of a new streamable HTTP client
func TestUnitNewStreamableHTTPClient(t *testing.T) {
	config := StreamableHTTPClientConfig{
		ID:              "test-client-1",
		URL:             "https://example.com/mcp",
		ProtocolVersion: MCPVersion20241105,
		JSONRPCVersion:  JSONRPC20,
		ClientInfo: ClientInfo{
			Name:    "test-client",
			Version: "1.0.0",
		},
		Timeout: 5 * time.Second,
	}

	client, err := NewStreamableHTTPClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)

	defer client.Close()

	httpClient := client.(*StreamableHTTPClient)
	assert.Equal(t, "test-client-1", httpClient.id)
	assert.Equal(t, "https://example.com/mcp", httpClient.url)
	assert.Equal(t, MCPVersion20241105, httpClient.protocolVersion)
	assert.Equal(t, JSONRPC20, httpClient.jsonRPCVersion)
	assert.Equal(t, "test-client", httpClient.clientInfo.Name)
	assert.Equal(t, 5*time.Second, httpClient.timeout)
	assert.False(t, httpClient.initialized)
}

// TestNewStreamableHTTPClientDefaults tests default values
func TestUnitNewStreamableHTTPClientDefaults(t *testing.T) {
	config := StreamableHTTPClientConfig{
		URL: "https://example.com/mcp",
	}

	client, err := NewStreamableHTTPClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)

	defer client.Close()

	httpClient := client.(*StreamableHTTPClient)
	assert.Equal(t, MCPVersion20241105, httpClient.protocolVersion)
	assert.Equal(t, JSONRPC20, httpClient.jsonRPCVersion)
	assert.Equal(t, "mut-client", httpClient.clientInfo.Name)
	assert.Equal(t, "0.1.0-dev", httpClient.clientInfo.Version)
	assert.Equal(t, 30*time.Second, httpClient.timeout)
}

// TestNewStreamableHTTPClientMissingURL tests error for missing URL
func TestUnitNewStreamableHTTPClientMissingURL(t *testing.T) {
	config := StreamableHTTPClientConfig{}

	_, err := NewStreamableHTTPClient(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP URL is required")
}

// TestStreamableHTTPClient_Initialize tests the initialization process
func TestIntegrationStreamableHTTPClient_Initialize(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	config := StreamableHTTPClientConfig{
		ID:      "test-client-init",
		URL:     server.URL + "/mcp",
		Timeout: 5 * time.Second,
	}

	client, err := NewStreamableHTTPClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Initialize(ctx)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "2024-11-05", result.ProtocolVersion)
	assert.Equal(t, "mock-server", result.ServerInfo.Name)
	assert.Equal(t, "1.0.0", result.ServerInfo.Version)

	// Test double initialization
	_, err = client.Initialize(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already initialized")
}

// TestStreamableHTTPClient_InitializeSSE tests initialization with SSE response
func TestIntegrationStreamableHTTPClient_InitializeSSE(t *testing.T) {
	server := mockSSEMCPServer(t)
	defer server.Close()

	config := StreamableHTTPClientConfig{
		ID:      "test-client-sse",
		URL:     server.URL + "/mcp",
		Timeout: 5 * time.Second,
	}

	client, err := NewStreamableHTTPClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Initialize(ctx)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "2024-11-05", result.ProtocolVersion)
	assert.Equal(t, "mock-sse-server", result.ServerInfo.Name)
	assert.Equal(t, "1.0.0", result.ServerInfo.Version)
}

// TestStreamableHTTPClient_ListTools tests listing tools
func TestIntegrationStreamableHTTPClient_ListTools(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewStreamableHTTPClient(StreamableHTTPClientConfig{
		URL:     server.URL + "/mcp",
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Should fail before initialization
	_, err = client.ListTools(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")

	// Initialize first
	_, err = client.Initialize(ctx)
	require.NoError(t, err)

	// Now list tools
	tools, err := client.ListTools(ctx)
	require.NoError(t, err)
	require.Len(t, tools, 1)
	assert.Equal(t, "test_tool", tools[0].Name)
	assert.Equal(t, "A test tool", tools[0].Description)
}

// TestStreamableHTTPClient_CallTool tests calling a tool
func TestIntegrationStreamableHTTPClient_CallTool(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewStreamableHTTPClient(StreamableHTTPClientConfig{
		URL:     server.URL + "/mcp",
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Initialize first
	_, err = client.Initialize(ctx)
	require.NoError(t, err)

	// Now call tool
	result, err := client.CallTool(ctx, "test_tool", ToolArgument{"name": "test_tool"})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Content, 1)
	assert.Equal(t, "text", result.Content[0].Type)
	assert.Contains(t, result.Content[0].Text, "Tool result")
}

// TestStreamableHTTPClient_ListPrompts tests listing prompts
func TestIntegrationStreamableHTTPClient_ListPrompts(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewStreamableHTTPClient(StreamableHTTPClientConfig{
		URL:     server.URL + "/mcp",
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Initialize
	_, err = client.Initialize(ctx)
	require.NoError(t, err)

	prompts, err := client.ListPrompts(ctx)
	require.NoError(t, err)
	require.Len(t, prompts, 1)
	assert.Equal(t, "test_prompt", prompts[0].Name)
}

// TestStreamableHTTPClient_ListResources tests listing resources
func TestIntegrationStreamableHTTPClient_ListResources(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewStreamableHTTPClient(StreamableHTTPClientConfig{
		URL:     server.URL + "/mcp",
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Initialize
	_, err = client.Initialize(ctx)
	require.NoError(t, err)

	resources, err := client.ListResources(ctx)
	require.NoError(t, err)
	require.Len(t, resources, 1)
	assert.Equal(t, "test_resource", resources[0].Name)
}

// TestStreamableHTTPClient_Discover tests full discovery
func TestIntegrationStreamableHTTPClient_Discover(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewStreamableHTTPClient(StreamableHTTPClientConfig{
		URL:     server.URL + "/mcp",
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	result, err := client.Discover(ctx)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "mock-server", result.ServerInfo.Name)
	assert.Len(t, result.Tools, 1)
	assert.Len(t, result.Prompts, 1)
	assert.Len(t, result.Resources, 1)
}

// TestStreamableHTTPClient_Headers tests custom headers
func TestIntegrationStreamableHTTPClient_Headers(t *testing.T) {
	headerReceived := false
	expectedHeader := "Bearer test-token"
	expectedValue := "custom-value"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		customHeader := r.Header.Get("X-Custom-Header")

		if authHeader == expectedHeader && customHeader == expectedValue {
			headerReceived = true
		}

		// Parse and respond
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check if this is a notification (no ID)
		if req.ID == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      *req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]interface{}{},
				"serverInfo": map[string]interface{}{
					"name":    "mock-server",
					"version": "1.0.0",
				},
			},
		})
	}))
	defer server.Close()

	config := StreamableHTTPClientConfig{
		URL: server.URL + "/mcp",
		Headers: map[string]string{
			"Authorization":   expectedHeader,
			"X-Custom-Header": expectedValue,
		},
		Timeout: 5 * time.Second,
	}

	client, err := NewStreamableHTTPClient(config)
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	_, err = client.Initialize(ctx)
	require.NoError(t, err)

	assert.True(t, headerReceived, "Custom headers should be sent")
}

// TestStreamableHTTPClient_SessionID tests session ID handling
func TestIntegrationStreamableHTTPClient_SessionID(t *testing.T) {
	sessionIDs := []string{}
	mu := &sync.Mutex{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Store session ID from request
		receivedSessionID := r.Header.Get("mcp-session-id")

		// Parse and respond
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check if this is a notification (no ID)
		if req.ID == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		mu.Lock()
		sessionIDs = append(sessionIDs, receivedSessionID)
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("mcp-session-id", "test-session-123")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      *req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]interface{}{},
				"serverInfo": map[string]interface{}{
					"name":    "mock-server",
					"version": "1.0.0",
				},
			},
		})
	}))
	defer server.Close()

	config := StreamableHTTPClientConfig{
		URL:     server.URL + "/mcp",
		Timeout: 5 * time.Second,
	}

	client, err := NewStreamableHTTPClient(config)
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// First request - no session ID yet
	_, err = client.Initialize(ctx)
	require.NoError(t, err)

	// First request should not have session ID
	mu.Lock()
	assert.Equal(t, "", sessionIDs[0], "First request should not have session ID")
	mu.Unlock()

	// Get the session ID that was stored
	httpClient := client.(*StreamableHTTPClient)
	httpClient.sessionMutex.Lock()
	storedSessionID := httpClient.sessionID
	httpClient.sessionMutex.Unlock()
	assert.Equal(t, "test-session-123", storedSessionID)

	// Second request - should now include session ID
	_, err = client.ListTools(ctx)
	require.NoError(t, err)

	// Second request should include session ID
	mu.Lock()
	assert.Equal(t, "test-session-123", sessionIDs[1], "Second request should include session ID")
	mu.Unlock()
}

// TestStreamableHTTPClient_Close tests client cleanup
func TestUnitStreamableHTTPClient_Close(t *testing.T) {
	config := StreamableHTTPClientConfig{
		URL:     "https://example.com/mcp",
		Timeout: 5 * time.Second,
	}

	client, err := NewStreamableHTTPClient(config)
	require.NoError(t, err)

	err = client.Close()
	assert.NoError(t, err)

	// Verify state is reset
	httpClient := client.(*StreamableHTTPClient)
	assert.False(t, httpClient.initialized)
	assert.Nil(t, httpClient.serverInfo)
}

// TestStreamableHTTPClient_Timeout tests request timeout handling
func TestIntegrationStreamableHTTPClient_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout test in short mode")
	}

	// Create a slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(2 * time.Second) // Longer than client timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := StreamableHTTPClientConfig{
		URL:     server.URL + "/mcp",
		Timeout: 500 * time.Millisecond,
	}

	client, err := NewStreamableHTTPClient(config)
	require.NoError(t, err)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err = client.Initialize(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}
