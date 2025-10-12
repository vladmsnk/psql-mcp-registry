package instance_manager

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"psql-mcp-registry/internal/instance_manager/mocks"
	"psql-mcp-registry/internal/model"
)

func TestGetInstance_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockStorage := mocks.NewStorage(t)

	impl := &Implementation{
		storage: mockStorage,
	}

	instanceName := "test-instance"
	expectedInstance := &model.Instance{
		ID:              1,
		Name:            instanceName,
		DatabaseName:    "test_db",
		Description:     "Test database instance",
		CreatorUsername: "testuser",
		Status:          "active",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockStorage.On("GetInstanceByName", ctx, instanceName).Return(expectedInstance, nil)

	result, err := impl.GetInstance(ctx, instanceName)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedInstance, result)
}
