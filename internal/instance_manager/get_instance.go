package instance_manager

import (
	"context"

	"psql-mcp-registry/internal/model"
)

func (i *Implementation) GetInstance(ctx context.Context, instanceName string) (*model.Instance, error) {
	instance, err := i.storage.GetInstanceByName(ctx, instanceName)
	if err != nil {
		return nil, err
	}

	return instance, nil
}
