// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package discovery

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitParseToolsListResponse(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		expected  []Tool
		expectErr bool
	}{
		{
			name: "Valid tools list",
			input: map[string]interface{}{
				"tools": []interface{}{
					map[string]interface{}{
						"name":        "calculator",
						"description": "A simple calculator",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"operation": map[string]interface{}{
									"type": "string",
								},
							},
						},
					},
					map[string]interface{}{
						"name":        "text_formatter",
						"description": "Format text",
						"inputSchema": map[string]interface{}{
							"type": "object",
						},
					},
				},
			},
			expected: []Tool{
				{
					Name:        "calculator",
					Description: "A simple calculator",
					InputSchema: map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"operation": map[string]interface{}{
								"type": "string",
							},
						},
					},
				},
				{
					Name:        "text_formatter",
					Description: "Format text",
					InputSchema: map[string]interface{}{
						"type": "object",
					},
				},
			},
			expectErr: false,
		},
		{
			name: "Empty tools list",
			input: map[string]interface{}{
				"tools": []interface{}{},
			},
			expected:  []Tool{},
			expectErr: false,
		},
		{
			name:      "Nil input",
			input:     nil,
			expected:  []Tool(nil),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseToolsListResponse(tt.input)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, len(tt.expected))

				for i, tool := range result {
					assert.Equal(t, tt.expected[i].Name, tool.Name)
					assert.Equal(t, tt.expected[i].Description, tool.Description)
				}
			}
		})
	}
}

func TestUnitParsePromptsListResponse(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		expected  []Prompt
		expectErr bool
	}{
		{
			name: "Valid prompts list",
			input: map[string]interface{}{
				"prompts": []interface{}{
					map[string]interface{}{
						"name":        "greeting",
						"description": "Generate a greeting",
						"arguments": []interface{}{
							map[string]interface{}{
								"name":        "name",
								"description": "Name to greet",
								"required":    true,
							},
						},
					},
					map[string]interface{}{
						"name":        "summary",
						"description": "Generate a summary",
					},
				},
			},
			expected: []Prompt{
				{
					Name:        "greeting",
					Description: "Generate a greeting",
					Arguments: []PromptArgument{
						{
							Name:        "name",
							Description: "Name to greet",
							Required:    true,
						},
					},
				},
				{
					Name:        "summary",
					Description: "Generate a summary",
					Arguments:   nil,
				},
			},
			expectErr: false,
		},
		{
			name: "Empty prompts list",
			input: map[string]interface{}{
				"prompts": []interface{}{},
			},
			expected:  []Prompt{},
			expectErr: false,
		},
		{
			name:      "Nil input",
			input:     nil,
			expected:  []Prompt(nil),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePromptsListResponse(tt.input)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, len(tt.expected))

				for i, prompt := range result {
					assert.Equal(t, tt.expected[i].Name, prompt.Name)
					assert.Equal(t, tt.expected[i].Description, prompt.Description)
				}
			}
		})
	}
}

func TestUnitParseResourcesListResponse(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		expected  []Resource
		expectErr bool
	}{
		{
			name: "Valid resources list",
			input: map[string]interface{}{
				"resources": []interface{}{
					map[string]interface{}{
						"uri":         "file:///tmp/data.txt",
						"name":        "data.txt",
						"description": "Data file",
						"mimeType":    "text/plain",
					},
					map[string]interface{}{
						"uri":         "http://example.com/api",
						"name":        "API",
						"description": "REST API",
						"mimeType":    "application/json",
					},
				},
			},
			expected: []Resource{
				{
					URI:         "file:///tmp/data.txt",
					Name:        "data.txt",
					Description: "Data file",
					MimeType:    "text/plain",
				},
				{
					URI:         "http://example.com/api",
					Name:        "API",
					Description: "REST API",
					MimeType:    "application/json",
				},
			},
			expectErr: false,
		},
		{
			name: "Empty resources list",
			input: map[string]interface{}{
				"resources": []interface{}{},
			},
			expected:  []Resource{},
			expectErr: false,
		},
		{
			name:      "Nil input",
			input:     nil,
			expected:  []Resource(nil),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseResourcesListResponse(tt.input)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, len(tt.expected))

				for i, resource := range result {
					assert.Equal(t, tt.expected[i].URI, resource.URI)
					assert.Equal(t, tt.expected[i].Name, resource.Name)
				}
			}
		})
	}
}

func TestUnitParseInitializeResponse(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		expected  *InitializeResult
		expectErr bool
	}{
		{
			name: "Valid initialize response",
			input: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools":     map[string]interface{}{"enabled": true},
					"prompts":   map[string]interface{}{"enabled": true},
					"resources": map[string]interface{}{"enabled": false},
				},
				"serverInfo": map[string]interface{}{
					"name":            "test-server",
					"version":         "1.0.0",
					"protocolVersion": "2024-11-05",
				},
			},
			expected: &InitializeResult{
				ProtocolVersion: "2024-11-05",
				Capabilities: map[string]interface{}{
					"tools":     map[string]interface{}{"enabled": true},
					"prompts":   map[string]interface{}{"enabled": true},
					"resources": map[string]interface{}{"enabled": false},
				},
				ServerInfo: ServerInfo{
					Name:            "test-server",
					Version:         "1.0.0",
					ProtocolVersion: "2024-11-05",
				},
			},
			expectErr: false,
		},
		{
			name:      "Nil input",
			input:     nil,
			expected:  &InitializeResult{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseInitializeResponse(tt.input)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.ProtocolVersion, result.ProtocolVersion)
				assert.Equal(t, tt.expected.ServerInfo.Name, result.ServerInfo.Name)
			}
		})
	}
}

func TestUnitJSONRPCRequestSerialization(t *testing.T) {
	id := 1
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      &id,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2024-11-05",
		},
	}

	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var decoded JSONRPCRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, req.JSONRPC, decoded.JSONRPC)
	assert.Equal(t, req.Method, decoded.Method)
	assert.Equal(t, *req.ID, *decoded.ID)
}

func TestUnitJSONRPCResponseSerialization(t *testing.T) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result: map[string]interface{}{
			"status": "ok",
		},
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	var decoded JSONRPCResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, resp.JSONRPC, decoded.JSONRPC)
	assert.Equal(t, resp.ID, decoded.ID)
}

func TestUnitJSONRPCErrorSerialization(t *testing.T) {
	rpcErr := JSONRPCError{
		Code:    -32600,
		Message: "Invalid Request",
		Data:    "Additional error data",
	}

	data, err := json.Marshal(rpcErr)
	assert.NoError(t, err)

	var decoded JSONRPCError
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, rpcErr.Code, decoded.Code)
	assert.Equal(t, rpcErr.Message, decoded.Message)
	assert.Equal(t, rpcErr.Data, decoded.Data)
}

func TestUnitToolSerialization(t *testing.T) {
	tool := Tool{
		Name:        "calculator",
		Description: "A calculator tool",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"operation": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}

	data, err := json.Marshal(tool)
	assert.NoError(t, err)

	var decoded Tool
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, tool.Name, decoded.Name)
	assert.Equal(t, tool.Description, decoded.Description)
}

func TestUnitPromptSerialization(t *testing.T) {
	prompt := Prompt{
		Name:        "greeting",
		Description: "Generate a greeting",
		Arguments: []PromptArgument{
			{
				Name:        "name",
				Description: "Name to greet",
				Required:    true,
			},
		},
	}

	data, err := json.Marshal(prompt)
	assert.NoError(t, err)

	var decoded Prompt
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, prompt.Name, decoded.Name)
	assert.Equal(t, prompt.Description, decoded.Description)
	assert.Len(t, decoded.Arguments, len(prompt.Arguments))
}

func TestUnitResourceSerialization(t *testing.T) {
	resource := Resource{
		URI:         "file:///tmp/data.txt",
		Name:        "data.txt",
		Description: "Data file",
		MimeType:    "text/plain",
	}

	data, err := json.Marshal(resource)
	assert.NoError(t, err)

	var decoded Resource
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, resource.URI, decoded.URI)
	assert.Equal(t, resource.Name, decoded.Name)
	assert.Equal(t, resource.MimeType, decoded.MimeType)
}

func TestUnitToolCallResultSerialization(t *testing.T) {
	result := ToolCallResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: "Result of calculation: 42",
			},
		},
		IsError: false,
	}

	data, err := json.Marshal(result)
	assert.NoError(t, err)

	var decoded ToolCallResult
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Len(t, decoded.Content, len(result.Content))
	assert.Equal(t, result.IsError, decoded.IsError)
}

func TestUnitDiscoveryResultSerialization(t *testing.T) {
	discoveryResult := Result{
		ServerInfo: ServerInfo{
			Name:            "test-server",
			Version:         "1.0.0",
			ProtocolVersion: "2024-11-05",
		},
		Tools: []Tool{
			{
				Name:        "calculator",
				Description: "A calculator",
				InputSchema: map[string]interface{}{
					"type": "object",
				},
			},
		},
		Prompts: []Prompt{
			{
				Name:        "greeting",
				Description: "Generate greeting",
			},
		},
		Resources: []Resource{
			{
				URI:      "file:///data.txt",
				Name:     "data.txt",
				MimeType: "text/plain",
			},
		},
	}

	data, err := json.Marshal(discoveryResult)
	assert.NoError(t, err)

	var decoded Result
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, discoveryResult.ServerInfo.Name, decoded.ServerInfo.Name)
	assert.Len(t, decoded.Tools, len(discoveryResult.Tools))
	assert.Len(t, decoded.Prompts, len(discoveryResult.Prompts))
	assert.Len(t, decoded.Resources, len(discoveryResult.Resources))
}

func TestUnitInitializeParamsSerialization(t *testing.T) {
	params := InitializeParams{
		ProtocolVersion: "2024-11-05",
		Capabilities: map[string]interface{}{
			"tools": map[string]interface{}{
				"enabled": true,
			},
		},
		ClientInfo: ClientInfo{
			Name:    "test-client",
			Version: "1.0.0",
		},
	}

	data, err := json.Marshal(params)
	assert.NoError(t, err)

	var decoded InitializeParams
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, params.ProtocolVersion, decoded.ProtocolVersion)
	assert.Equal(t, params.ClientInfo.Name, decoded.ClientInfo.Name)
}

func TestUnitResourceContentSerialization(t *testing.T) {
	content := ResourceContent{
		URI:      "file:///data.txt",
		MimeType: "text/plain",
		Text:     "Hello, World!",
	}

	data, err := json.Marshal(content)
	assert.NoError(t, err)

	var decoded ResourceContent
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, content.URI, decoded.URI)
	assert.Equal(t, content.Text, decoded.Text)
}

func TestUnitPromptMessageSerialization(t *testing.T) {
	message := PromptMessage{
		Role: "assistant",
		Content: ToolContent{
			Type: "text",
			Text: "Hello!",
		},
	}

	data, err := json.Marshal(message)
	assert.NoError(t, err)

	var decoded PromptMessage
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, message.Role, decoded.Role)
	assert.Equal(t, message.Content.Text, decoded.Content.Text)
}

func TestUnitConstants(t *testing.T) {
	assert.Equal(t, JSONRPCVersion("2.0"), JSONRPC20)
	assert.Equal(t, MCPProtocolVersion("2024-11-05"), MCPVersion20241105)
	assert.Equal(t, MCPProtocolVersion("2024-11-25"), MCPVersion20241125)
	assert.Equal(t, MCPProtocolVersion("2025-06-18"), MCPVersion20250618)
	assert.Equal(t, MCPVersion20241105, DefaultMCPVersion)
	assert.Equal(t, MCPVersion20250618, LatestMCPVersion)
}

func TestUnitSupportedVersions(t *testing.T) {
	versions := SupportedVersions()
	assert.Len(t, versions, 3)
	assert.Equal(t, MCPVersion20241105, versions[0])
	assert.Equal(t, MCPVersion20241125, versions[1])
	assert.Equal(t, MCPVersion20250618, versions[2])
}

func TestUnitIsVersionSupported(t *testing.T) {
	tests := []struct {
		name     string
		version  MCPProtocolVersion
		expected bool
	}{
		{
			name:     "2024-11-05 is supported",
			version:  MCPVersion20241105,
			expected: true,
		},
		{
			name:     "2024-11-25 is supported",
			version:  MCPVersion20241125,
			expected: true,
		},
		{
			name:     "2025-06-18 is supported",
			version:  MCPVersion20250618,
			expected: true,
		},
		{
			name:     "2025-06-18 is supported",
			version:  MCPProtocolVersion("2025-06-18"),
			expected: true,
		},
		{
			name:     "Unknown version is not supported",
			version:  MCPProtocolVersion("2099-12-31"),
			expected: false,
		},
		{
			name:     "Empty version is not supported",
			version:  MCPProtocolVersion(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsVersionSupported(tt.version)
			assert.Equal(t, tt.expected, result)
		})
	}
}
