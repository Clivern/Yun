// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package v1

import (
	"fmt"
	"net/http"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// Health checks the health status of the Mut API
func (c *Client) Health() (*HealthResponse, error) {
	resp, err := c.doRequest(http.MethodGet, "/api/v1/public/_health", nil)
	if err != nil {
		return nil, err
	}

	var healthResp HealthResponse
	if err := c.parseJSONResponse(resp, &healthResp); err != nil {
		return nil, err
	}

	return &healthResp, nil
}

// Health checks the health status of the Mut API at the given base URL
func Health(baseURL string) (*HealthResponse, error) {
	client, err := NewClient(ClientConfig{BaseURL: baseURL})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client.Health()
}

