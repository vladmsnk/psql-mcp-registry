package factory

import (
	"errors"
	"testing"

	"psql-mcp-registry/internal/factory/mocks"
	"psql-mcp-registry/internal/model"
	"psql-mcp-registry/internal/pg"

	"github.com/stretchr/testify/assert"
)

func TestPGClientFactory_CreateClient_Success(t *testing.T) {
	// Arrange
	instance := model.Instance{
		Name:         "test-instance",
		DatabaseName: "testdb",
		Description:  "Test PostgreSQL instance",
		Status:       "active",
	}

	expectedConfig := &pg.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "testuser",
		Password: "testpass",
		Database: "testdb",
		SSLMode:  "disable",
	}

	mockLoader := mocks.NewConfigLoader(t)
	mockLoader.On("Load", instance.Name).Return(expectedConfig, nil)

	factory := NewPGClientFactory(mockLoader)

	// Act
	client, err := factory.CreateClient(instance)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestPGClientFactory_CreateClient_ConfigLoaderError(t *testing.T) {
	// Arrange
	instance := model.Instance{
		Name:         "test-instance",
		DatabaseName: "testdb",
		Description:  "Test PostgreSQL instance",
		Status:       "active",
	}

	expectedError := errors.New("configuration not found")

	mockLoader := mocks.NewConfigLoader(t)
	mockLoader.On("Load", instance.Name).Return((*pg.Config)(nil), expectedError)

	factory := NewPGClientFactory(mockLoader)

	// Act
	client, err := factory.CreateClient(instance)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Equal(t, expectedError, err)
}

func TestPGClientFactory_CreateClient_InvalidConfig(t *testing.T) {
	// Arrange
	instance := model.Instance{
		Name:         "test-instance",
		DatabaseName: "testdb",
		Description:  "Test PostgreSQL instance",
		Status:       "active",
	}

	// Return nil config to trigger NewClient error
	mockLoader := mocks.NewConfigLoader(t)
	mockLoader.On("Load", instance.Name).Return((*pg.Config)(nil), nil)

	factory := NewPGClientFactory(mockLoader)

	// Act
	client, err := factory.CreateClient(instance)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, client)
}
