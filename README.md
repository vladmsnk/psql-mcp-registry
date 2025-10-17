# PostgreSQL MCP Registry

A PostgreSQL-based Model Context Protocol (MCP) registry that manages multiple PostgreSQL instances.

## Features

- **Dynamic Instance Management**: Add and remove PostgreSQL instances at runtime
- **Thread-safe Operations**: Concurrent access to registry with proper locking
- **Lazy Connection**: Registry creation doesn't fail if some instances are unavailable
- **Docker Support**: Easily spin up test PostgreSQL instances

## Architecture

The registry implements a fourth variant pattern that:
1. Creates the registry even if some instances fail to connect
2. Allows lazy loading of instances via `AddInstance()`
3. Provides thread-safe access to clients
4. Supports graceful shutdown with `Close()`

## HTTP API

The service provides an HTTP API for managing PostgreSQL instances alongside the MCP protocol interface. Both servers run in parallel.

### Environment Variables

- `HTTP_API_PORT` - Port for the HTTP API server (default: `8080`)
- `GIN_MODE` - Gin framework mode: `release` or `debug` (default: `release`)

### API Endpoints

#### Register Instance
```bash
POST /api/v1/instances
Content-Type: application/json

{
  "name": "prod_db",
  "database_name": "production",
  "description": "Production PostgreSQL instance",
  "creator_username": "admin"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "prod_db",
  "database_name": "production",
  "description": "Production PostgreSQL instance",
  "creator_username": "admin",
  "status": "active",
  "created_at": "2025-10-15T10:30:00Z",
  "updated_at": "2025-10-15T10:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request body or missing required fields
- `409 Conflict` - Instance with this name already exists
- `500 Internal Server Error` - Registration failed

#### Health Check
```bash
GET /health
```

**Response (200 OK):**
```json
{
  "status": "healthy"
}
```

### Usage Examples

**Register a new instance:**
```bash
curl -X POST http://localhost:8080/api/v1/instances \
  -H "Content-Type: application/json" \
  -d '{
    "name": "prod_db",
    "database_name": "production",
    "description": "Production PostgreSQL instance",
    "creator_username": "admin"
  }'
```

**Check API health:**
```bash
curl http://localhost:8080/health
```

**Note:** Instance connection details must be configured via environment variables following the `PSQL_INSTANCE_{NAME}_*` pattern (see Environment Variable Format section below).

## Quick Start

### 1. Start Test PostgreSQL Instances

```bash
docker-compose up -d
```

This starts two PostgreSQL instances:
- **postgres-test**: `localhost:5432` (user: testuser, password: testpass, db: testdb)
- **postgres-test-dev**: `localhost:5433` (user: devuser, password: devpass, db: devdb)

### 2. Configure Environment Variables

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` to match your instances:
```bash
# Primary instance
PSQL_INSTANCE_PROD_HOST=localhost
PSQL_INSTANCE_PROD_PORT=5432
PSQL_INSTANCE_PROD_USER=testuser
PSQL_INSTANCE_PROD_PASSWORD=testpass
PSQL_INSTANCE_PROD_DATABASE=testdb
PSQL_INSTANCE_PROD_SSLMODE=disable

# Development instance
PSQL_INSTANCE_DEV_HOST=localhost
PSQL_INSTANCE_DEV_PORT=5433
PSQL_INSTANCE_DEV_USER=devuser
PSQL_INSTANCE_DEV_PASSWORD=devpass
PSQL_INSTANCE_DEV_DATABASE=devdb
PSQL_INSTANCE_DEV_SSLMODE=disable
```

### 3. Run Tests

```bash
# Run all tests
go test ./...

# Run registry tests with verbose output
go test -v ./internal/registry/

# Check health of docker instances
docker-compose ps
```

## Usage Example

```go
package main

import (
    "context"
    "log"
    "psql-mcp-registry/internal/registry"
    "psql-mcp-registry/internal/storage/instances"
)

func main() {
    ctx := context.Background()
    
    // Create storage backend
    storage := instances.NewPostgresStorage()
    
    // Create registry (won't fail if some instances are down)
    reg, err := registry.NewRegistry(ctx, storage)
    if err != nil {
        log.Fatal(err)
    }
    defer reg.Close()
    
    // Get a client for a specific instance
    client, ok := reg.GetClient("prod")
    if !ok {
        log.Println("prod instance not available")
        return
    }
    
    // Use the client...
    // client.Query(...)
    
    // Add a new instance dynamically
    newInstance := registry.Instance{Name: "staging"}
    if err := reg.AddInstance(ctx, newInstance); err != nil {
        log.Printf("Failed to add staging: %v", err)
    }
    
    // List all available instances
    instances := reg.ListInstances()
    log.Printf("Available instances: %v", instances)
}
```

## Registry API

### `NewRegistry(ctx, storage) (*Registry, error)`
Creates a new registry. Will succeed even if some instances fail to connect.

### `AddInstance(ctx, instance) error`
Adds a new instance to the registry (lazy loading).

### `GetClient(name) (*pg.Client, bool)`
Returns the client for a given instance name.

### `RemoveInstance(name) error`
Removes an instance and closes its connection.

### `ListInstances() []string`
Returns all registered instance names.

### `Close() error`
Closes all client connections.

## Testing

The test suite uses mocks to avoid requiring real database connections:

- **TestNewRegistry_Success**: Verifies registry creation with failing instances
- **TestNewRegistry_ListInstancesError**: Tests storage error handling
- **TestRegistry_GetClient**: Tests client retrieval
- **TestRegistry_ListInstances**: Tests instance listing
- **TestRegistry_AddInstance_NoConfig**: Tests adding instances without config
- **TestInstance_LoadConfigByInstanceName_Success**: Tests config loading
- **TestInstance_LoadConfigByInstanceName_MissingHost**: Tests missing config

## Docker Commands

```bash
# Start instances
docker-compose up -d

# View logs
docker-compose logs -f

# Stop instances
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v

# Check instance health
docker-compose ps
```

## Environment Variable Format

Instance configuration follows this pattern:
```
PSQL_INSTANCE_<NAME>_HOST=hostname
PSQL_INSTANCE_<NAME>_PORT=5432
PSQL_INSTANCE_<NAME>_USER=username
PSQL_INSTANCE_<NAME>_PASSWORD=password
PSQL_INSTANCE_<NAME>_DATABASE=dbname
PSQL_INSTANCE_<NAME>_SSLMODE=disable
PSQL_INSTANCE_<NAME>_MAX_OPEN_CONNS=25
PSQL_INSTANCE_<NAME>_MAX_IDLE_CONNS=10
```

Where `<NAME>` is the uppercase instance name (e.g., `PROD`, `DEV`, `STAGING`).

## License

MIT

