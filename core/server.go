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

	r.Use(chimiddleware.Recoverer)
	if viper.GetInt("app.timeout") > 0 {
		timeout := time.Duration(viper.GetInt("app.timeout")) * time.Second
		r.Use(chimiddleware.Timeout(timeout))
	}
	r.Use(middleware.PrometheusMiddleware)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestSizeLimit(int64(10 * 1024 * 1024)))
	r.Use(middleware.SessionAuth())

	// Routes
	r.Get("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Public Actions
	r.Group(func(r chi.Router) {
		r.Get("/api/v1/public/_health", api.HealthAction)
		r.Get("/api/v1/public/_ready", api.ReadyAction)
		r.Post("/api/v1/public/action/setup", api.SetupAction)
		r.Get("/api/v1/public/action/setup/status", api.SetupStatusAction)
		r.Post("/api/v1/public/action/login", api.LoginAction)
		r.Post("/api/v1/public/action/logout", api.LogoutAction)
	})
	// Private Actions
	r.Group(func(r chi.Router) {
		r.Get("/api/v1/action/profile", api.GetProfileAction)
		r.Put("/api/v1/action/profile", api.UpdateProfileAction)
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireRole(db.UserRoleAdmin))
		r.Put("/api/v1/action/settings", api.UpdateSettingsAction)
		r.Get("/api/v1/action/settings", api.GetSettingsAction)
	})
	// Users routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireRole(db.UserRoleAdmin))
		r.Post("/api/v1/users", api.CreateUserAction)
		r.Get("/api/v1/users", api.ListUsersAction)
		r.Get("/api/v1/users/{id}", api.GetUserAction)
		r.Put("/api/v1/users/{id}", api.UpdateUserAction)
		r.Delete("/api/v1/users/{id}", api.DeleteUserAction)
	})
	r.With(middleware.BasicAuth(
		viper.GetString("app.metrics.username"),
		viper.GetString("app.metrics.secret"),
	)).Get("/api/v1/public/_metrics", promhttp.Handler().ServeHTTP)

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
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := dist.Open("index.html")
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		stat, err := indexFile.Stat()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

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
	if err := InitDatabase(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	defer func() {
		if err := db.CloseDB(); err != nil {
			log.Error().Err(err).Msg("Error closing database connection")
		}
	}()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("app.port"))),
		Handler: handler,
	}

	serverErrors := make(chan error, 1)

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

		// Ignore ErrServerClosed as it's expected during graceful shutdown
		if err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Info().
			Str("signal", sig.String()).
			Msg("Received shutdown signal")

		shutdownTimeout := 30 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		log.Info().
			Dur("timeout", shutdownTimeout).
			Msg("Gracefully shutting down server")

		// Shutdown with timeout to allow in-flight requests to complete
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}

		log.Info().Msg("Server shutdown complete")
	}

	return nil
}
