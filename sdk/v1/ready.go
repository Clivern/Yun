// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package v1

import (
	"fmt"
	"net/http"
)

// ReadyResponse represents the readiness check response
type ReadyResponse struct {
	Status string `json:"status"`
}

// Ready checks if the Mut API is ready to serve traffic
func (c *Client) Ready() (*ReadyResponse, error) {
	resp, err := c.doRequest(http.MethodGet, "/api/v1/public/_ready", nil)
	if err != nil {
		return nil, err
	}

	var readyResp ReadyResponse
	if err := c.parseJSONResponse(resp, &readyResp); err != nil {
		return nil, err
	}

	return &readyResp, nil
}

// Ready checks if the Mut API is ready to serve traffic at the given base URL
func Ready(baseURL string) (*ReadyResponse, error) {
	client, err := NewClient(ClientConfig{BaseURL: baseURL})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client.Ready()
}
