package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"psql-mcp-registry/internal/api"
	"psql-mcp-registry/internal/factory"
	"psql-mcp-registry/internal/instance_manager"
	mcpserver "psql-mcp-registry/internal/mcp"
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

	// Create MCP server
	mcpServer := mcpserver.NewMCPServer(queryRouter, instanceManager)
	log.Println("Initialized MCP server")

	// Read HTTP API port from environment variable (default: 8080)
	httpPort := os.Getenv("HTTP_API_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// Read MCP SSE port from environment variable (default: 3000)
	mcpPort := os.Getenv("MCP_PORT")
	if mcpPort == "" {
		mcpPort = "3000"
	}

	// Create HTTP API server
	apiServer := api.NewAPIServer(instanceManager, httpPort)
	log.Printf("Initialized HTTP API server on port %s", httpPort)
	log.Printf("MCP server will use SSE transport on port %s", mcpPort)

	// Log successful initialization
	log.Println("Application initialized successfully")
	log.Println("Starting servers...")

	// Run both servers in goroutines
	errChan := make(chan error, 2)

	// Run MCP server with SSE transport
	go func() {
		log.Printf("Starting MCP server with SSE transport on :%s", mcpPort)
		log.Printf("SSE endpoint: http://localhost:%s/sse", mcpPort)
		log.Printf("Test with: npx @modelcontextprotocol/inspector http://localhost:%s/sse", mcpPort)
		if err := mcpServer.RunWithSSE(ctx, mcpPort); err != nil {
			errChan <- err
		}
	}()

	// Run HTTP API server
	go func() {
		log.Printf("Starting HTTP API server on :%s", httpPort)
		if err := apiServer.Run(ctx); err != nil {
			errChan <- err
		}
	}()

	// Wait for interrupt signal or error
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		log.Println("Received interrupt signal, shutting down gracefully...")
		cancel()
	case err := <-errChan:
		log.Printf("Server error: %v", err)
		cancel()
	}

	log.Println("Shutdown complete")
}
