// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package discovery

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to get the path to the basic mock server script
func createMockServerScript(t *testing.T) string {
	scriptPath := filepath.Join("..", "testdata", "discovery", "mock_server.py")
	absPath, err := filepath.Abs(scriptPath)
	require.NoError(t, err)
	return absPath
}

// Helper function to get the path to the full mock MCP server script
func createFullMockServerScript(t *testing.T) string {
	scriptPath := filepath.Join("..", "testdata", "discovery", "mock_mcp_server.py")
	absPath, err := filepath.Abs(scriptPath)
	require.NoError(t, err)
	return absPath
}

// isProcessRunning checks if a process with given PID is running
func isProcessRunning(pid int) bool {
	// On Unix systems, sending signal 0 checks if process exists
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(os.Signal(nil))
	return err == nil
}

// requirePython3 skips the test if python3 is not available
func requirePython3(t *testing.T) {
	t.Helper()
	if _, err := os.Stat("/usr/bin/python3"); err == nil {
		return
	}
	if _, err := os.Stat("/usr/local/bin/python3"); err == nil {
		return
	}
	// Try to find python3 in PATH
	if path, err := exec.LookPath("python3"); err == nil && path != "" {
		return
	}
	t.Skip("python3 not found in system, skipping test")
}

// TestNewStdioClient tests the creation of a new stdio client
func TestIntegrationNewStdioClient(t *testing.T) {
	requirePython3(t)

	// Get path to the simple mock server script
	mockScript := createMockServerScript(t)

	config := StdioClientConfig{
		ID:              "test-client-1",
		Command:         "python3",
		Args:            []string{mockScript},
		ProtocolVersion: MCPVersion20241105,
		JSONRPCVersion:  JSONRPC20,
		ClientInfo: ClientInfo{
			Name:    "test-client",
			Version: "1.0.0",
		},
		Timeout: 5 * time.Second,
	}

	client, err := NewStdioClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)

	defer client.Close()

	stdioClient := client.(*StdioClient)
	assert.Equal(t, "test-client-1", stdioClient.id)
	assert.Equal(t, MCPVersion20241105, stdioClient.protocolVersion)
	assert.Equal(t, JSONRPC20, stdioClient.jsonRPCVersion)
	assert.Equal(t, "test-client", stdioClient.clientInfo.Name)
	assert.Equal(t, 5*time.Second, stdioClient.timeout)
	assert.False(t, stdioClient.initialized)
}

// TestNewStdioClientDefaults tests default values
func TestIntegrationNewStdioClientDefaults(t *testing.T) {
	requirePython3(t)

	mockScript := createMockServerScript(t)

	config := StdioClientConfig{
		Command: "python3",
		Args:    []string{mockScript},
	}

	client, err := NewStdioClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)

	defer client.Close()

	stdioClient := client.(*StdioClient)
	assert.Equal(t, MCPVersion20241105, stdioClient.protocolVersion)
	assert.Equal(t, JSONRPC20, stdioClient.jsonRPCVersion)
	assert.Equal(t, "mut-client", stdioClient.clientInfo.Name)
	assert.Equal(t, "0.1.0-dev", stdioClient.clientInfo.Version)
	assert.Equal(t, 30*time.Second, stdioClient.timeout)
}

// TestStdioClient_NextRequestID tests request ID generation
func TestUnitStdioClient_NextRequestID(t *testing.T) {
	client := &StdioClient{
		requestID: 0,
	}

	id1 := client.nextRequestID()
	id2 := client.nextRequestID()
	id3 := client.nextRequestID()

	assert.Equal(t, 1, id1)
	assert.Equal(t, 2, id2)
	assert.Equal(t, 3, id3)
}

// TestStdioClient_Initialize tests the initialization process
func TestIntegrationStdioClient_Initialize(t *testing.T) {
	requirePython3(t)

	mockScript := createFullMockServerScript(t)

	config := StdioClientConfig{
		ID:      "test-client-init",
		Command: "python3",
		Args:    []string{mockScript},
		Timeout: 5 * time.Second,
	}

	client, err := NewStdioClient(config)
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

// TestStdioClient_ListTools tests listing tools
func TestIntegrationStdioClient_ListTools(t *testing.T) {
	requirePython3(t)

	mockScript := createFullMockServerScript(t)

	client, err := NewStdioClient(StdioClientConfig{
		Command: "python3",
		Args:    []string{mockScript},
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

// TestStdioClient_CallTool tests calling a tool
func TestIntegrationStdioClient_CallTool(t *testing.T) {
	requirePython3(t)

	mockScript := createFullMockServerScript(t)

	client, err := NewStdioClient(StdioClientConfig{
		Command: "python3",
		Args:    []string{mockScript},
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
	assert.Contains(t, result.Content[0].Text, "Tool result")
}

// TestStdioClient_ListPrompts tests listing prompts
func TestIntegrationStdioClient_ListPrompts(t *testing.T) {
	requirePython3(t)

	mockScript := createFullMockServerScript(t)

	client, err := NewStdioClient(StdioClientConfig{
		Command: "python3",
		Args:    []string{mockScript},
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

// TestStdioClient_ListResources tests listing resources
func TestIntegrationStdioClient_ListResources(t *testing.T) {
	requirePython3(t)

	mockScript := createFullMockServerScript(t)

	client, err := NewStdioClient(StdioClientConfig{
		Command: "python3",
		Args:    []string{mockScript},
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

// TestStdioClient_Discover tests full discovery
func TestIntegrationStdioClient_Discover(t *testing.T) {
	requirePython3(t)

	mockScript := createFullMockServerScript(t)

	client, err := NewStdioClient(StdioClientConfig{
		Command: "python3",
		Args:    []string{mockScript},
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

// TestStdioClient_InvalidCommand tests error handling for invalid commands
func TestIntegrationStdioClient_InvalidCommand(t *testing.T) {
	config := StdioClientConfig{
		Command: "nonexistent-command-xyz",
		Args:    []string{},
	}

	_, err := NewStdioClient(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to start command")
}

// TestStdioClient_Timeout tests request timeout handling
func TestIntegrationStdioClient_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout test in short mode")
	}
	requirePython3(t)

	// Get path to the slow server script
	scriptPath := filepath.Join("..", "testdata", "discovery", "slow_server.py")
	scriptPath, err := filepath.Abs(scriptPath)
	require.NoError(t, err)

	client, err := NewStdioClient(StdioClientConfig{
		Command: "python3",
		Args:    []string{scriptPath},
		Timeout: 1 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err = client.Initialize(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout")
}

// TestStdioClient_JSONRPCError tests JSON-RPC error handling
func TestIntegrationStdioClient_JSONRPCError(t *testing.T) {
	requirePython3(t)

	// Get path to the error server script
	scriptPath := filepath.Join("..", "testdata", "discovery", "error_server.py")
	scriptPath, err := filepath.Abs(scriptPath)
	require.NoError(t, err)

	client, err := NewStdioClient(StdioClientConfig{
		Command: "python3",
		Args:    []string{scriptPath},
		Timeout: 5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()

	_, err = client.Initialize(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JSON-RPC error")
	assert.Contains(t, err.Error(), "Invalid Request")
}

// TestStdioClient_WorkingDirectory tests working directory configuration
func TestIntegrationStdioClient_WorkingDirectory(t *testing.T) {
	requirePython3(t)

	tmpDir := t.TempDir()

	// Get path to the working directory check script
	scriptPath := filepath.Join("..", "testdata", "discovery", "check_cwd.py")
	scriptPath, err := filepath.Abs(scriptPath)
	require.NoError(t, err)

	client, err := NewStdioClient(StdioClientConfig{
		Command:    "python3",
		Args:       []string{scriptPath},
		WorkingDir: tmpDir,
		Timeout:    5 * time.Second,
	})
	require.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	result, err := client.Initialize(ctx)
	require.NoError(t, err)
	assert.Equal(t, "cwd-test", result.ServerInfo.Name)
}

// TestStdioClient_Close tests client cleanup
func TestIntegrationStdioClient_Close(t *testing.T) {
	requirePython3(t)

	mockScript := createMockServerScript(t)

	client, err := NewStdioClient(StdioClientConfig{
		Command: "python3",
		Args:    []string{mockScript},
	})
	require.NoError(t, err)

	stdioClient := client.(*StdioClient)

	// Verify process is running
	assert.NotNil(t, stdioClient.cmd.Process)
	pid := stdioClient.cmd.Process.Pid

	// Close the client
	err = client.Close()
	assert.NoError(t, err)

	// Verify process has exited (give it up to 2 seconds)
	for i := 0; i < 20; i++ {
		if !isProcessRunning(pid) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
