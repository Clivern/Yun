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
				"description": "Test prompt description",
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
						"uri":      req.Params["uri"].(string),
						"mimeType": "text/plain",
						"text":     "Resource content",
					},
				},
			}

		default:
			if req.ID != nil {
				response.ID = *req.ID
				response.Error = &JSONRPCError{
					Code:    -32601,
					Message: "Method not found",
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	return httptest.NewServer(http.HandlerFunc(handler))
}

// TestIntegrationNewSSEClient tests the creation of a new SSE client
func TestIntegrationNewSSEClient(t *testing.T) {
	config := SSEClientConfig{
		ID:              "test-client-1",
		URL:             "https://api.example.com/mcp/sse",
		ProtocolVersion: MCPVersion20241105,
		JSONRPCVersion:  JSONRPC20,
		ClientInfo: ClientInfo{
			Name:    "test-client",
			Version: "1.0.0",
		},
		Timeout: 5 * time.Second,
	}

	client, err := NewSSEClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)

	defer client.Close()

	sseClient := client.(*SSEClient)
	assert.Equal(t, "test-client-1", sseClient.id)
	assert.Equal(t, MCPVersion20241105, sseClient.protocolVersion)
	assert.Equal(t, JSONRPC20, sseClient.jsonRPCVersion)
	assert.Equal(t, "test-client", sseClient.clientInfo.Name)
	assert.Equal(t, 5*time.Second, sseClient.timeout)
	assert.False(t, sseClient.initialized)
	assert.NotNil(t, sseClient.httpClient)
}

// TestIntegrationNewSSEClientDefaults tests default values
func TestIntegrationNewSSEClientDefaults(t *testing.T) {
	config := SSEClientConfig{
		URL: "https://api.example.com/mcp/sse",
	}

	client, err := NewSSEClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)

	defer client.Close()

	sseClient := client.(*SSEClient)
	assert.Equal(t, MCPVersion20241105, sseClient.protocolVersion)
	assert.Equal(t, JSONRPC20, sseClient.jsonRPCVersion)
	assert.Equal(t, "mut-client", sseClient.clientInfo.Name)
	assert.Equal(t, "0.1.0-dev", sseClient.clientInfo.Version)
	assert.Equal(t, 30*time.Second, sseClient.timeout)
	assert.NotNil(t, sseClient.headers)
}

// TestIntegrationNewSSEClient_MissingURL tests error handling for missing URL
func TestIntegrationNewSSEClient_MissingURL(t *testing.T) {
	config := SSEClientConfig{}

	_, err := NewSSEClient(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SSE URL is required")
}

// TestIntegrationNewSSEClient_Headers tests header initialization
func TestIntegrationNewSSEClient_Headers(t *testing.T) {
	config := SSEClientConfig{
		URL: "https://api.example.com/mcp/sse",
		Headers: map[string]string{
			"Authorization": "Bearer token123",
			"X-API-Key":     "key456",
		},
	}

	client, err := NewSSEClient(config)
	require.NoError(t, err)
	defer client.Close()

	sseClient := client.(*SSEClient)
	assert.Equal(t, "Bearer token123", sseClient.headers["Authorization"])
	assert.Equal(t, "key456", sseClient.headers["X-API-Key"])
}

// TestIntegrationNewSSEClient_NilHeaders tests nil headers are initialized
func TestIntegrationNewSSEClient_NilHeaders(t *testing.T) {
	config := SSEClientConfig{
		URL:     "https://api.example.com/mcp/sse",
		Headers: nil,
	}

	client, err := NewSSEClient(config)
	require.NoError(t, err)
	defer client.Close()

	sseClient := client.(*SSEClient)
	assert.NotNil(t, sseClient.headers)
	assert.Empty(t, sseClient.headers)
}

// TestIntegrationSSEClient_NextRequestID tests request ID generation
func TestIntegrationSSEClient_NextRequestID(t *testing.T) {
	client := &SSEClient{
		requestID: 0,
	}

	id1 := client.nextRequestID()
	id2 := client.nextRequestID()
	id3 := client.nextRequestID()

	assert.Equal(t, 1, id1)
	assert.Equal(t, 2, id2)
	assert.Equal(t, 3, id3)
}

// TestIntegrationSSEClient_Initialize tests the initialization process
func TestIntegrationSSEClient_Initialize(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		ID:      "test-client-init",
		URL:     server.URL,
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.Initialize(ctx)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "2024-11-05", result.ProtocolVersion)
	assert.Equal(t, "mock-server", result.ServerInfo.Name)
	assert.Equal(t, "1.0.0", result.ServerInfo.Version)

	sseClient := client.(*SSEClient)
	assert.True(t, sseClient.initialized)
	assert.NotNil(t, sseClient.serverInfo)

	// Test double initialization
	_, err = client.Initialize(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already initialized")
}

// TestIntegrationSSEClient_InitializeWithHeaders tests initialization with custom headers
func TestIntegrationSSEClient_InitializeWithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Authorization header is present
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method == "initialize" && req.ID != nil {
			response := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      *req.ID,
				Result: map[string]interface{}{
					"protocolVersion": "2024-11-05",
					"capabilities":    map[string]interface{}{},
					"serverInfo": map[string]interface{}{
						"name":    "auth-server",
						"version": "1.0.0",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		URL: server.URL,
		Headers: map[string]string{
			"Authorization": "Bearer test-token",
		},
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	result, err := client.Initialize(ctx)
	require.NoError(t, err)
	assert.Equal(t, "auth-server", result.ServerInfo.Name)
}

// TestIntegrationSSEClient_ListTools tests listing tools
func TestIntegrationSSEClient_ListTools(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		URL:     server.URL,
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

// TestIntegrationSSEClient_CallTool tests calling a tool
func TestIntegrationSSEClient_CallTool(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Should fail before initialization
	_, err = client.CallTool(ctx, "test_tool", ToolArgument{"input": "test"})
	assert.Error(t, err)

	// Initialize first
	_, err = client.Initialize(ctx)
	require.NoError(t, err)

	// Now call tool
	result, err := client.CallTool(ctx, "test_tool", ToolArgument{"input": "test"})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Content, 1)
	assert.Equal(t, "text", result.Content[0].Type)
	assert.Contains(t, result.Content[0].Text, "test_tool")
}

// TestIntegrationSSEClient_ListPrompts tests listing prompts
func TestIntegrationSSEClient_ListPrompts(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		URL:     server.URL,
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

// TestIntegrationSSEClient_GetPrompt tests getting a prompt
func TestIntegrationSSEClient_GetPrompt(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Initialize
	_, err = client.Initialize(ctx)
	require.NoError(t, err)

	result, err := client.GetPrompt(ctx, "test_prompt", map[string]string{"arg1": "value1"})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Test prompt description", result.Description)
	assert.Len(t, result.Messages, 1)
	assert.Equal(t, "user", result.Messages[0].Role)
}

// TestIntegrationSSEClient_ListResources tests listing resources
func TestIntegrationSSEClient_ListResources(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		URL:     server.URL,
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
	assert.Equal(t, "test://resource", resources[0].URI)
}

// TestIntegrationSSEClient_ReadResource tests reading a resource
func TestIntegrationSSEClient_ReadResource(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	// Initialize
	_, err = client.Initialize(ctx)
	require.NoError(t, err)

	result, err := client.ReadResource(ctx, "test://resource")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Contents, 1)
	assert.Equal(t, "test://resource", result.Contents[0].URI)
	assert.Equal(t, "text/plain", result.Contents[0].MimeType)
	assert.Equal(t, "Resource content", result.Contents[0].Text)
}

// TestIntegrationSSEClient_Discover tests full discovery
func TestIntegrationSSEClient_Discover(t *testing.T) {
	server := mockMCPServer(t)
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		URL:     server.URL,
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

// TestIntegrationSSEClient_HTTPError tests HTTP error handling
func TestIntegrationSSEClient_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	_, err = client.Initialize(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP error 500")
}

// TestIntegrationSSEClient_JSONRPCError tests JSON-RPC error handling
func TestIntegrationSSEClient_JSONRPCError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		response := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      *req.ID,
			Error: &JSONRPCError{
				Code:    -32601,
				Message: "Method not found",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	_, err = client.Initialize(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JSON-RPC error")
	assert.Contains(t, err.Error(), "Method not found")
}

// TestIntegrationSSEClient_Timeout tests request timeout handling
func TestIntegrationSSEClient_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout test in short mode")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		URL:     server.URL,
		Timeout: 500 * time.Millisecond,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = client.Initialize(ctx)
	assert.Error(t, err)
	// Either context timeout or HTTP client timeout
	assert.True(t, strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "context deadline exceeded"))
}

// TestIntegrationSSEClient_InvalidJSON tests invalid JSON response handling
func TestIntegrationSSEClient_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client, err := NewSSEClient(SSEClientConfig{
		URL:     server.URL,
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	_, err = client.Initialize(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal")
}

// TestIntegrationSSEClient_Close tests client cleanup
func TestIntegrationSSEClient_Close(t *testing.T) {
	client, err := NewSSEClient(SSEClientConfig{
		URL: "https://api.example.com/mcp/sse",
	})
	require.NoError(t, err)

	sseClient := client.(*SSEClient)
	sseClient.initialized = true
	sseClient.serverInfo = &ServerInfo{Name: "test"}

	err = client.Close()
	assert.NoError(t, err)

	assert.False(t, sseClient.initialized)
	assert.Nil(t, sseClient.serverInfo)
}
