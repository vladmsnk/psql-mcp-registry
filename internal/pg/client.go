package pg

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

//go:generate mockery --case snake --name ClientInterface
type ClientInterface interface {
	GetDatabaseOverview(ctx context.Context, dbName string) (*DatabaseOverview, error)
	GetCacheHitRateDB(ctx context.Context, dbName string) (*CacheHitRate, error)
	GetCacheHitRateGlobal(ctx context.Context) (*CacheHitRate, error)
	GetCheckpointsStats(ctx context.Context) (*CheckpointsStats, error)
	GetWalActivity(ctx context.Context) (*WalActivity, error)
	GetTablesInfo(ctx context.Context, limit int) ([]TableInfo, error)
	GetLockingInfo(ctx context.Context, dbName string) ([]LockInfo, error)
	GetChangedSettings(ctx context.Context) ([]SettingInfo, error)
	GetIndexStats(ctx context.Context, limit int) ([]IndexStats, error)
	GetActiveQueries(ctx context.Context, dbName string, minDuration int) ([]ActiveQuery, error)
	GetConnectionStats(ctx context.Context) (*ConnectionSummary, error)
	GetSlowQueries(ctx context.Context, limit int) ([]SlowQuery, error)
	GetDatabaseSizes(ctx context.Context) ([]DatabaseSize, error)
	Version() *Version
}

type Client struct {
	db      *sql.DB
	config  *Config
	version *Version
	mu      sync.RWMutex
}

func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	db, err := sql.Open("postgres", config.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	client := &Client{
		db:     db,
		config: config,
	}

	return client, nil
}

func (c *Client) Connect(ctx context.Context) error {
	// Проверка соединения
	if err := c.db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Определение версии PostgreSQL
	version, err := DetectVersion(ctx, c.db)
	if err != nil {
		return fmt.Errorf("failed to detect PostgreSQL version: %w", err)
	}

	c.mu.Lock()
	c.version = version
	c.mu.Unlock()

	return nil
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *Client) DB() *sql.DB {
	return c.db
}

func (c *Client) Version() *Version {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.version
}

func (c *Client) Config() *Config {
	return c.config
}

func (c *Client) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c *Client) Stats() sql.DBStats {
	return c.db.Stats()
}
