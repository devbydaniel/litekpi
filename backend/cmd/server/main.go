package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/trackable/trackable/internal/platform/config"
	"github.com/trackable/trackable/internal/platform/database"
	"github.com/trackable/trackable/internal/platform/router"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Create context that listens for shutdown signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to database
	log.Println("Connecting to database...")
	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	// Create router
	r := router.New(db, cfg)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		serverErrors <- server.ListenAndServe()
	}()

	// Wait for shutdown signal or server error
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("Received signal %v, starting graceful shutdown", sig)

		// Create shutdown context with timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
		defer shutdownCancel()

		// Gracefully shutdown the server
		if err := server.Shutdown(shutdownCtx); err != nil {
			// Force shutdown if graceful shutdown fails
			server.Close()
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		log.Println("Server stopped gracefully")
	}

	return nil
}
