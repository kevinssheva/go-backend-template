package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kevinssheva/go-backend-template/internal/config"
	"github.com/kevinssheva/go-backend-template/internal/registry"
	"github.com/kevinssheva/go-backend-template/internal/router"
	"github.com/kevinssheva/go-backend-template/pkg/database"
	"github.com/kevinssheva/go-backend-template/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(cfg.App.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("starting application",
		zap.String("env", cfg.App.Env),
		zap.String("log_level", cfg.App.LogLevel),
	)

	dbCfg := database.Config{
		DSN: cfg.Database.GetDSN(),
	}
	db, err := database.NewPostgresDB(dbCfg, log)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer database.Close(db, log)

	repos := registry.NewRepos(db, log)
	services := registry.NewServices(repos, log)
	handlers := registry.NewHandlers(services, log)

	router := router.NewRouter(handlers)

	srv := &http.Server{
		Addr:         cfg.Server.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("server starting", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server exited")
}
