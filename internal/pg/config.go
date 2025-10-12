package pg

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config содержит настройки подключения к PostgreSQL
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnTimeout     time.Duration
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "",
		Database:        "postgres",
		SSLMode:         "disable",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * time.Minute,
		ConnTimeout:     10 * time.Second,
	}
}

// LoadConfigFromEnv загружает конфигурацию из переменных окружения
func LoadConfigFromEnv() *Config {
	cfg := DefaultConfig()

	if host := os.Getenv("PGHOST"); host != "" {
		cfg.Host = host
	}
	if port := os.Getenv("PGPORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}
	if user := os.Getenv("PGUSER"); user != "" {
		cfg.User = user
	}
	if password := os.Getenv("PGPASSWORD"); password != "" {
		cfg.Password = password
	}
	if database := os.Getenv("PGDATABASE"); database != "" {
		cfg.Database = database
	}
	if sslmode := os.Getenv("PGSSLMODE"); sslmode != "" {
		cfg.SSLMode = sslmode
	}

	return cfg
}

// ConnectionString возвращает строку подключения для lib/pq
func (c *Config) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Database,
		c.SSLMode,
		int(c.ConnTimeout.Seconds()),
	)
}
