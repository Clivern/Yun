// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package core

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/clivern/mut/api"
	"github.com/clivern/mut/db"
	"github.com/clivern/mut/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Setup creates and configures the HTTP server
func Setup(Static embed.FS) http.Handler {
	r := chi.NewRouter()

	// Recover middleware
	r.Use(chimiddleware.Recoverer)

	// Request timeout middleware
	if viper.GetInt("app.timeout") > 0 {
		timeout := time.Duration(viper.GetInt("app.timeout")) * time.Second
		r.Use(chimiddleware.Timeout(timeout))
	}

	// Prometheus middleware for HTTP metrics
	r.Use(middleware.PrometheusMiddleware)

	// Custom logger middleware
	r.Use(middleware.Logger)

	// Request size limit middleware
	r.Use(middleware.RequestSizeLimit(int64(10 * 1024 * 1024)))

	// Session authentication middleware
	r.Use(middleware.SessionAuth())

	// Routes
	r.Get("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	r.Get("/api/v1/public/_health", api.HealthAction)
	r.Get("/api/v1/public/_ready", api.ReadyAction)

	// Setup routes (no authentication required)
	r.Post("/api/v1/public/action/setup", api.SetupAction)
	r.Get("/api/v1/public/action/setup/status", api.SetupStatusAction)

	// Authentication routes (no authentication required for login)
	r.Post("/api/v1/public/action/login", api.LoginAction)

	// Protected routes (authentication required via session cookie)
	r.Post("/api/v1/public/action/logout", api.LogoutAction)

	r.Get("/api/v1/action/profile", api.GetProfileAction)
	r.Put("/api/v1/action/profile", api.UpdateProfileAction)
	r.Put("/api/v1/action/settings", api.UpdateSettingsAction)

	// Metrics endpoint with basic auth
	r.With(middleware.BasicAuth(
		viper.GetString("app.metrics.username"),
		viper.GetString("app.metrics.secret"),
	)).Get("/api/v1/public/_metrics", promhttp.Handler().ServeHTTP)

	// Serve static files from embedded web/dist
	dist, err := fs.Sub(Static, "web/dist")

	if err != nil {
		panic(fmt.Sprintf(
			"Error while accessing dist files: %s",
			err.Error(),
		))
	}

	// Serve static assets (CSS, JS, images, etc.)
	r.Handle("/assets/*", http.StripPrefix("/", http.FileServer(http.FS(dist))))

	// SPA fallback: serve index.html for all other routes
	// This allows Vue Router to handle client-side routing
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// Read index.html from embedded filesystem
		indexFile, err := dist.Open("index.html")
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		// Get file info to set proper headers
		stat, err := indexFile.Stat()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Serve index.html
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
	})

	return r
}

// InitDatabase initializes the database connection from configuration
func InitDatabase() error {
	dbConfig := db.Config{
		Driver:          viper.GetString("app.database.driver"),
		Host:            viper.GetString("app.database.host"),
		Port:            viper.GetInt("app.database.port"),
		Username:        viper.GetString("app.database.username"),
		Password:        viper.GetString("app.database.password"),
		Database:        viper.GetString("app.database.name"),
		MaxOpenConns:    viper.GetInt("app.database.max_open_conns"),
		MaxIdleConns:    viper.GetInt("app.database.max_idle_conns"),
		ConnMaxLifetime: viper.GetInt("app.database.conn_max_lifetime"),
		DataSource:      viper.GetString("app.database.datasource"),
	}

	return db.InitDB(dbConfig)
}

// Run starts the HTTP server with graceful shutdown support
func Run(handler http.Handler) error {
	// Initialize database connection
	if err := InitDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Ensure database connection is closed on exit
	defer func() {
		if err := db.CloseDB(); err != nil {
			log.Error().Err(err).Msg("Error closing database connection")
		}
	}()

	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("app.port"))),
		Handler: handler,
	}

	// Channel to listen for errors from the server
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Info().
			Int("port", viper.GetInt("app.port")).
			Bool("tls", viper.GetBool("app.tls.status")).
			Msg("Starting HTTP server")

		var err error
		if viper.GetBool("app.tls.status") {
			err = srv.ListenAndServeTLS(
				viper.GetString("app.tls.crt_path"),
				viper.GetString("app.tls.key_path"),
			)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Info().
			Str("signal", sig.String()).
			Msg("Received shutdown signal")

		// Create a deadline for graceful shutdown
		shutdownTimeout := 30 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		log.Info().
			Dur("timeout", shutdownTimeout).
			Msg("Gracefully shutting down server")

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}

		log.Info().Msg("Server shutdown complete")
	}

	return nil
}
