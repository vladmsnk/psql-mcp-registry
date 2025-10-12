package factory

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"psql-mcp-registry/internal/pg"
)

type EnvConfigLoader struct{}

func NewEnvConfigLoader() ConfigLoader {
	return &EnvConfigLoader{}
}

func (e *EnvConfigLoader) Load(instanceName string) (*pg.Config, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instance name cannot be empty")
	}

	cfg := pg.DefaultConfig()

	prefix := fmt.Sprintf("PSQL_INSTANCE_%s_", strings.ToUpper(instanceName))

	hostKey := prefix + "HOST"
	host := os.Getenv(hostKey)
	if host == "" {
		return nil, fmt.Errorf("configuration for instance '%s' not found (missing %s)", instanceName, hostKey)
	}
	cfg.Host = host

	if port := os.Getenv(prefix + "PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}

	if user := os.Getenv(prefix + "USER"); user != "" {
		cfg.User = user
	}

	if password := os.Getenv(prefix + "PASSWORD"); password != "" {
		cfg.Password = password
	}

	if database := os.Getenv(prefix + "DATABASE"); database != "" {
		cfg.Database = database
	}

	if sslmode := os.Getenv(prefix + "SSLMODE"); sslmode != "" {
		cfg.SSLMode = sslmode
	}

	if maxOpenConns := os.Getenv(prefix + "MAX_OPEN_CONNS"); maxOpenConns != "" {
		if m, err := strconv.Atoi(maxOpenConns); err == nil {
			cfg.MaxOpenConns = m
		}
	}

	if maxIdleConns := os.Getenv(prefix + "MAX_IDLE_CONNS"); maxIdleConns != "" {
		if m, err := strconv.Atoi(maxIdleConns); err == nil {
			cfg.MaxIdleConns = m
		}
	}

	return cfg, nil
}
