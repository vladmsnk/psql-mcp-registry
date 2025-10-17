package instance_manager

import (
	"context"

	"psql-mcp-registry/internal/model"
)

type Manager interface {
	RegisterInstance(ctx context.Context, instance model.Instance) error
	GetInstance(ctx context.Context, instanceName string) (*model.Instance, error)
	ListInstances(ctx context.Context) ([]model.Instance, error)
}

//go:generate mockery --case snake --name Storage
type Storage interface {
	CreateInstance(ctx context.Context, instance *model.Instance) error
	GetInstanceByName(ctx context.Context, name string) (*model.Instance, error)
	ListInstances(ctx context.Context) ([]model.Instance, error)
}

//go:generate mockery --case snake --name InstanceRegistry
type InstanceRegistry interface {
	AddInstanceToRegistry(instance model.Instance) error
}

type Implementation struct {
	storage  Storage
	registry InstanceRegistry
}

func NewManager(storage Storage, registry InstanceRegistry) Manager {
	return &Implementation{
		storage:  storage,
		registry: registry,
	}
}
