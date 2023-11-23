// Copyright 2023 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/clivern/badger/internal/api"
	mid "github.com/clivern/badger/internal/middleware"
	"github.com/clivern/badger/internal/service"

	"github.com/drone/envsubst"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the badger backend server",
	Run: func(cmd *cobra.Command, args []string) {
		configUnparsed, err := ioutil.ReadFile(config)

		if err != nil {
			panic(fmt.Sprintf(
				"Error while reading config file [%s]: %s",
				config,
				err.Error(),
			))
		}

		configParsed, err := envsubst.EvalEnv(string(configUnparsed))

		if err != nil {
			panic(fmt.Sprintf(
				"Error while parsing config file [%s]: %s",
				config,
				err.Error(),
			))
		}

		viper.SetConfigType("yaml")
		err = viper.ReadConfig(bytes.NewBuffer([]byte(configParsed)))

		if err != nil {
			panic(fmt.Sprintf(
				"Error while loading configs [%s]: %s",
				config,
				err.Error(),
			))
		}

		fys := service.NewFileSystem()

		if viper.GetString("server.log.output") != "stdout" {
			dir, _ := filepath.Split(viper.GetString("server.log.output"))

			// Create dir
			if !fys.DirExists(dir) {
				if err := fys.EnsureDir(dir, 0775); err != nil {
					panic(fmt.Sprintf(
						"Directory [%s] creation failed with error: %s",
						dir,
						err.Error(),
					))
				}
			}

			// Create log file if not exists
			if !fys.FileExists(viper.GetString("server.log.output")) {
				f, err := os.Create(viper.GetString("server.log.output"))
				if err != nil {
					panic(fmt.Sprintf(
						"Error while creating log file [%s]: %s",
						viper.GetString("server.log.output"),
						err.Error(),
					))
				}
				defer f.Close()
			}
		}

		if viper.GetString("server.log.output") == "stdout" {
			log.SetOutput(os.Stdout)
			gin.DefaultWriter = os.Stdout
		} else {
			f, _ := os.OpenFile(
				viper.GetString("server.log.output"),
				os.O_APPEND|os.O_CREATE|os.O_WRONLY,
				0775,
			)
			log.SetOutput(f)
			gin.DefaultWriter = f
		}

		lvl := strings.ToLower(viper.GetString("server.log.level"))
		level, err := log.ParseLevel(lvl)

		if err != nil {
			level = log.InfoLevel
		}

		log.SetLevel(level)

		if viper.GetString("server.log.format") == "json" {
			log.SetFormatter(&log.JSONFormatter{})
		} else {
			log.SetFormatter(&log.TextFormatter{})
		}

		viper.SetDefault("config", config)

		// Set Gin mode
		if viper.GetString("server.mode") == "dev" {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		r := gin.New()

		// Recover middleware
		r.Use(gin.Recovery())

		// Prometheus middleware for HTTP metrics
		r.Use(mid.PrometheusMiddleware())

		// Custom logger middleware
		r.Use(mid.Logger())

		// Request timeout middleware
		r.Use(func(c *gin.Context) {
			timeoutCtx, _ := c.Request.Context(), func() {}
			if viper.GetInt("server.timeout") > 0 {
				var cancelFn func()
				timeoutCtx, cancelFn = c.Request.Context(), func() {}
				_ = cancelFn
			}
			c.Request = c.Request.WithContext(timeoutCtx)
			c.Next()
		})

		// Routes
		r.GET("/favicon.ico", func(c *gin.Context) {
			c.String(http.StatusNoContent, "")
		})

		r.GET("/", api.HealthAction)
		r.GET("/health", api.HealthAction)

		// Metrics endpoint with basic auth
		r.GET("/metrics", gin.BasicAuth(gin.Accounts{
			viper.GetString("server.metrics.username"): viper.GetString("server.metrics.secret"),
		}), gin.WrapH(promhttp.Handler()))

		var runerr error

		if viper.GetBool("server.tls.status") {
			runerr = r.RunTLS(
				fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("server.port"))),
				viper.GetString("server.tls.crt_path"),
				viper.GetString("server.tls.key_path"),
			)
		} else {
			runerr = r.Run(
				fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("server.port"))),
			)
		}

		if runerr != nil && runerr != http.ErrServerClosed {
			panic(runerr.Error())
		}
	},
}

func init() {
	serverCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.prod.yml",
		"Absolute path to config file (required)",
	)
	serverCmd.MarkFlagRequired("config")
	rootCmd.AddCommand(serverCmd)
}
