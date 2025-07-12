// Package main is the entry point for the UserCenter application
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// Version is the application version
var Version = "dev"

func main() {
	// Initialize application using wire
	app, err := InitializeApp()
	if err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Get logger from server
	log := app.GetLogger()

	log.Info("Starting UserCenter application",
		zap.String("version", Version),
	)

	// Start server in a goroutine
	go func() {
		if err := app.Start(); err != nil {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), app.GetShutdownTimeout())
	defer cancel()

	// Shutdown server
	if err := app.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}
