package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"psql-mcp-registry/internal/factory"
	"psql-mcp-registry/internal/instance_manager"
	"psql-mcp-registry/internal/pg"
	"psql-mcp-registry/internal/registry"
	"psql-mcp-registry/internal/router"
	"psql-mcp-registry/internal/storage/instances"
	"psql-mcp-registry/migrations"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load PostgreSQL configuration from environment variables
	config := pg.LoadConfigFromEnv()
	log.Printf("Connecting to PostgreSQL at %s:%d/%s", config.Host, config.Port, config.Database)

	// Create PostgreSQL client for registry database
	client, err := pg.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create PostgreSQL client: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Error closing client: %v", err)
		}
	}()

	// Connect to database
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL")

	// Apply migrations
	if err := migrations.ApplyMigrations(client.DB()); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// Create instances storage
	instanceStorage := instances.NewPostgresStorage(client.DB())
	log.Println("Initialized instance storage")

	// Create client factory
	configLoader := factory.NewEnvConfigLoader()
	clientFactory := factory.NewPGClientFactory(configLoader)
	log.Println("Initialized client factory")

	// Create instance registry
	instanceRegistry, err := registry.NewRegistry(ctx, instanceStorage, clientFactory)
	if err != nil {
		log.Fatalf("Failed to create instance registry: %v", err)
	}
	log.Println("Initialized instance registry")

	// Create instance manager
	instanceManager := instance_manager.NewManager(instanceStorage, instanceRegistry)
	log.Println("Initialized instance manager")

	// Create router
	queryRouter := router.New(instanceRegistry)
	log.Println("Initialized query router")

	// Log successful initialization
	log.Println("Application initialized successfully")
	log.Println("Instance manager is ready to use")
	log.Println("Query router is ready to use")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Println("Application is running. Press Ctrl+C to exit...")
	<-sigChan

	log.Println("Shutting down gracefully...")

	// Use variables to avoid unused variable warning
	_ = instanceManager
	_ = queryRouter
}
