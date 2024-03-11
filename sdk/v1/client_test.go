// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  ClientConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  ClientConfig{BaseURL: "https://api.example.com"},
			wantErr: false,
		},
		{
			name:    "missing baseURL",
			config:  ClientConfig{BaseURL: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}
