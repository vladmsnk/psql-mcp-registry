package factory

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvConfigLoader_Load_Success(t *testing.T) {
	// Arrange
	instanceName := "TEST_INSTANCE"

	// Set environment variables
	os.Setenv("PSQL_INSTANCE_TEST_INSTANCE_HOST", "testhost")
	os.Setenv("PSQL_INSTANCE_TEST_INSTANCE_PORT", "5433")
	os.Setenv("PSQL_INSTANCE_TEST_INSTANCE_USER", "testuser")
	os.Setenv("PSQL_INSTANCE_TEST_INSTANCE_PASSWORD", "testpass")
	os.Setenv("PSQL_INSTANCE_TEST_INSTANCE_DATABASE", "testdb")
	os.Setenv("PSQL_INSTANCE_TEST_INSTANCE_SSLMODE", "require")
	os.Setenv("PSQL_INSTANCE_TEST_INSTANCE_MAX_OPEN_CONNS", "50")
	os.Setenv("PSQL_INSTANCE_TEST_INSTANCE_MAX_IDLE_CONNS", "25")

	// Clean up after test
	defer func() {
		os.Unsetenv("PSQL_INSTANCE_TEST_INSTANCE_HOST")
		os.Unsetenv("PSQL_INSTANCE_TEST_INSTANCE_PORT")
		os.Unsetenv("PSQL_INSTANCE_TEST_INSTANCE_USER")
		os.Unsetenv("PSQL_INSTANCE_TEST_INSTANCE_PASSWORD")
		os.Unsetenv("PSQL_INSTANCE_TEST_INSTANCE_DATABASE")
		os.Unsetenv("PSQL_INSTANCE_TEST_INSTANCE_SSLMODE")
		os.Unsetenv("PSQL_INSTANCE_TEST_INSTANCE_MAX_OPEN_CONNS")
		os.Unsetenv("PSQL_INSTANCE_TEST_INSTANCE_MAX_IDLE_CONNS")
	}()

	loader := NewEnvConfigLoader()

	// Act
	config, err := loader.Load(instanceName)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "testhost", config.Host)
	assert.Equal(t, 5433, config.Port)
	assert.Equal(t, "testuser", config.User)
	assert.Equal(t, "testpass", config.Password)
	assert.Equal(t, "testdb", config.Database)
	assert.Equal(t, "require", config.SSLMode)
	assert.Equal(t, 50, config.MaxOpenConns)
	assert.Equal(t, 25, config.MaxIdleConns)
}

func TestEnvConfigLoader_Load_MinimalConfig(t *testing.T) {
	// Arrange
	instanceName := "MINIMAL"

	// Set only required environment variable
	os.Setenv("PSQL_INSTANCE_MINIMAL_HOST", "minimalhost")

	// Clean up after test
	defer os.Unsetenv("PSQL_INSTANCE_MINIMAL_HOST")

	loader := NewEnvConfigLoader()

	// Act
	config, err := loader.Load(instanceName)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "minimalhost", config.Host)
	// Defaults should be applied for other fields
	assert.Equal(t, 5432, config.Port)
	assert.Equal(t, "postgres", config.User)
}

func TestEnvConfigLoader_Load_MissingHost(t *testing.T) {
	// Arrange
	instanceName := "MISSING_HOST"

	loader := NewEnvConfigLoader()

	// Act
	config, err := loader.Load(instanceName)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "configuration for instance")
	assert.Contains(t, err.Error(), "missing PSQL_INSTANCE_MISSING_HOST_HOST")
}

func TestEnvConfigLoader_Load_EmptyInstanceName(t *testing.T) {
	// Arrange
	loader := NewEnvConfigLoader()

	// Act
	config, err := loader.Load("")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "instance name cannot be empty")
}

func TestEnvConfigLoader_Load_LowercaseInstanceName(t *testing.T) {
	// Arrange - lowercase instance name should be converted to uppercase
	instanceName := "lowercase_test"

	os.Setenv("PSQL_INSTANCE_LOWERCASE_TEST_HOST", "lowercasehost")
	defer os.Unsetenv("PSQL_INSTANCE_LOWERCASE_TEST_HOST")

	loader := NewEnvConfigLoader()

	// Act
	config, err := loader.Load(instanceName)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "lowercasehost", config.Host)
}
