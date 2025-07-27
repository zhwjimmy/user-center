// Package main is the entry point for the UserCenter migration tool
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// Version is the application version
var Version = "dev"

func main() {
	// Parse command line flags
	var (
		configPath = flag.String("config", "configs/config.yaml", "Path to configuration file")
		action     = flag.String("action", "up", "Migration action: up, down, status, create")
		name       = flag.String("name", "", "Migration name (required for create action)")
		steps      = flag.Int("steps", 1, "Number of steps for down migration")
		version    = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	// Show version and exit
	if *version {
		fmt.Printf("UserCenter Migration Tool v%s\n", Version)
		os.Exit(0)
	}

	// Validate action
	if *action != "up" && *action != "down" && *action != "status" && *action != "create" {
		fmt.Fprintf(os.Stderr, "Invalid action: %s. Must be one of: up, down, status, create\n", *action)
		os.Exit(1)
	}

	// Validate name for create action
	if *action == "create" && *name == "" {
		fmt.Fprintf(os.Stderr, "Migration name is required for create action\n")
		os.Exit(1)
	}

	// Initialize application using wire
	app, err := InitializeMigrationApp(*configPath)
	if err != nil {
		fmt.Printf("Failed to initialize migration app: %v\n", err)
		os.Exit(1)
	}

	// Get logger
	log := app.GetLogger()

	log.Info("Starting UserCenter migration tool",
		zap.String("version", Version),
		zap.String("action", *action),
		zap.String("config", *configPath),
	)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Info("Received shutdown signal, cancelling migration...")
		cancel()
	}()

	// Execute migration action
	var err2 error
	switch *action {
	case "up":
		err2 = app.MigrateUp(ctx)
	case "down":
		err2 = app.MigrateDown(ctx, *steps)
	case "status":
		err2 = app.MigrateStatus(ctx)
	case "create":
		err2 = app.CreateMigration(ctx, *name)
	}

	if err2 != nil {
		log.Error("Migration failed", zap.Error(err2))
		os.Exit(1)
	}

	log.Info("Migration completed successfully")
}
