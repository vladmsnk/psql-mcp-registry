package instance_manager

import (
	"context"
	"testing"
	"time"

	"psql-mcp-registry/internal/instance_manager/mocks"
	"psql-mcp-registry/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestListInstances_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockStorage := mocks.NewStorage(t)

	impl := &Implementation{
		storage: mockStorage,
	}

	expectedInstances := []model.Instance{
		{
			ID:              1,
			Name:            "instance-1",
			DatabaseName:    "test_db_1",
			Description:     "Test database instance 1",
			CreatorUsername: "testuser",
			Status:          "active",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              2,
			Name:            "instance-2",
			DatabaseName:    "test_db_2",
			Description:     "Test database instance 2",
			CreatorUsername: "testuser",
			Status:          "active",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	mockStorage.On("ListInstances", ctx).Return(expectedInstances, nil)

	result, err := impl.ListInstances(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedInstances, result)
	assert.Len(t, result, 2)
}

func TestListInstances_EmptyList(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockStorage := mocks.NewStorage(t)

	impl := &Implementation{
		storage: mockStorage,
	}

	expectedInstances := []model.Instance{}

	mockStorage.On("ListInstances", ctx).Return(expectedInstances, nil)

	result, err := impl.ListInstances(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}
