package instance_manager

import (
	"context"
	"errors"
	"testing"

	"psql-mcp-registry/internal/instance_manager/mocks"
	"psql-mcp-registry/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRegisterInstance_Success(t *testing.T) {
	ctx := context.Background()
	mockStorage := mocks.NewStorage(t)
	mockRegistry := mocks.NewInstanceRegistry(t)

	impl := &Implementation{
		storage:  mockStorage,
		registry: mockRegistry,
	}

	instance := model.Instance{
		Name:            "test-instance",
		DatabaseName:    "test_db",
		Description:     "Test database instance",
		CreatorUsername: "testuser",
		Status:          "active",
	}

	mockStorage.On("GetInstanceByName", ctx, instance.Name).Return(nil, errors.New("not found"))

	mockStorage.On("CreateInstance", ctx, &instance).Return(nil)

	mockRegistry.On("AddInstanceToRegistry", instance).Return(nil)

	err := impl.RegisterInstance(ctx, instance)

	assert.NoError(t, err)
}
