package instance_manager

import (
	"context"
	"errors"

	"psql-mcp-registry/internal/model"
)

var (
	ErrInstanceAlreadyExists = errors.New("instance already exists")
)

func (i *Implementation) RegisterInstance(ctx context.Context, instance model.Instance) error {
	_, err := i.storage.GetInstanceByName(ctx, instance.Name)
	if err == nil {
		return ErrInstanceAlreadyExists
	}

	err = i.storage.CreateInstance(ctx, &instance)
	if err != nil {
		return err
	}

	err = i.registry.AddInstanceToRegistry(instance)
	if err != nil {
		return err
	}

	return nil
}
