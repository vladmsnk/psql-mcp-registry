package mcp

import (
	"context"
	"fmt"
	"net/http"

	"psql-mcp-registry/internal/model"
	"psql-mcp-registry/internal/router"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type InstanceManager interface {
	GetInstance(ctx context.Context, name string) (*model.Instance, error)
	ListInstances(ctx context.Context) ([]model.Instance, error)
}

type MCPServer struct {
	server  *mcp.Server
	router  *router.Router
	manager InstanceManager
}

func NewMCPServer(router *router.Router, manager InstanceManager) *MCPServer {
	impl := &mcp.Implementation{
		Name:    "psql-mcp-registry",
		Version: "v1.0.0",
	}

	server := mcp.NewServer(impl, nil)
	mcpServer := &MCPServer{
		server:  server,
		router:  router,
		manager: manager,
	}

	mcpServer.registerTools()
	mcpServer.registerResources()

	return mcpServer
}

func (s *MCPServer) registerResources() {
	// Instances List Resource
	s.server.AddResource(&mcp.Resource{
		URI:         "instances://list",
		Name:        "instances_list",
		Description: "List of all registered PostgreSQL instances",
		MIMEType:    "application/json",
	}, s.handleListInstancesResource)
}

func (s *MCPServer) registerTools() {
	// Database Overview
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "database_overview",
		Description: "Get overview statistics for a PostgreSQL database including transactions, blocks, tuples, and other metrics",
	}, s.handleDatabaseOverview)

	// Cache Hit Rate
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "cache_hit_rate",
		Description: "Get cache hit rate statistics (global or per database) to monitor buffer cache efficiency",
	}, s.handleCacheHitRate)

	// Checkpoints Stats
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "checkpoints_stats",
		Description: "Get checkpoint statistics including timed and requested checkpoints, buffers written, and sync times",
	}, s.handleCheckpointsStats)

	// WAL Activity
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "wal_activity",
		Description: "Get Write-Ahead Log activity statistics including WAL records, bytes, and FPI",
	}, s.handleWalActivity)

	// Tables Info
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "tables_info",
		Description: "Get information about tables including size, row count, and access patterns",
	}, s.handleTablesInfo)

	// Locking Info
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "locking_info",
		Description: "Get locking information for a database to identify blocking queries and lock conflicts",
	}, s.handleLockingInfo)

	// Changed Settings
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "changed_settings",
		Description: "Get PostgreSQL settings that differ from defaults to review configuration changes",
	}, s.handleChangedSettings)

	// Version
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "version",
		Description: "Get PostgreSQL version information",
	}, s.handleVersion)

	// Index Stats
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "index_stats",
		Description: "Get index usage statistics to identify unused or inefficient indexes",
	}, s.handleIndexStats)

	// Active Queries
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "active_queries",
		Description: "Get currently running queries with duration exceeding threshold for real-time performance diagnostics",
	}, s.handleActiveQueries)

	// Connection Stats
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "connection_stats",
		Description: "Get connection pool statistics including active, idle, and waiting connections",
	}, s.handleConnectionStats)

	// Slow Queries
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "slow_queries",
		Description: "Get top slow queries from pg_stat_statements with execution time and cache hit rate metrics",
	}, s.handleSlowQueries)

	// Database Sizes
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "database_sizes",
		Description: "Get sizes of all databases to monitor disk space usage and data growth",
	}, s.handleDatabaseSizes)
}

// Run starts the MCP server over stdio transport
func (s *MCPServer) Run(ctx context.Context) error {
	return s.server.Run(ctx, &mcp.StdioTransport{})
}

// RunWithTransport starts the MCP server with a custom transport
func (s *MCPServer) RunWithTransport(ctx context.Context, transport mcp.Transport) error {
	return s.server.Run(ctx, transport)
}

// RunWithSSE starts the MCP server with SSE transport over HTTP
func (s *MCPServer) RunWithSSE(ctx context.Context, port string) error {
	handler := mcp.NewSSEHandler(func(*http.Request) *mcp.Server {
		return s.server
	}, nil)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		httpServer.Shutdown(context.Background())
	}()

	return httpServer.ListenAndServe()
}

func createRouterRequest(instanceName string, action model.ActionName, params map[string]interface{}) router.QueryRequest {
	return router.QueryRequest{
		InstanceName: instanceName,
		Action:       action,
		Parameters:   params,
	}
}

func (s *MCPServer) executeRouterQuery(ctx context.Context, instanceName string, action model.ActionName, params map[string]interface{}) (interface{}, error) {
	instance, err := s.manager.GetInstance(ctx, instanceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	req := createRouterRequest(instanceName, action, params)
	response, err := s.router.RouteQuery(ctx, req, *instance)
	if err != nil {
		return nil, fmt.Errorf("router query failed: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("query failed: %s", response.Error)
	}

	return response.Data, nil
}
