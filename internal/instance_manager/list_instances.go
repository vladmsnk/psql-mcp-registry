package instance_manager

import (
	"context"

	"psql-mcp-registry/internal/model"
)

func (i *Implementation) ListInstances(ctx context.Context) ([]model.Instance, error) {
	instances, err := i.storage.ListInstances(ctx)
	if err != nil {
		return nil, err
	}

	return instances, nil
}
