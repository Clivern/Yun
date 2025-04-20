// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// Context key for request ID
type contextKey string

// RequestIDKey is the context key for storing request IDs
const RequestIDKey contextKey = "request_id"

// RequestID middleware adds a unique request ID to each request
// It checks for existing X-Request-ID header, or generates a new UUID
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")

		if requestID == "" {
			requestID = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", requestID)

		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}
