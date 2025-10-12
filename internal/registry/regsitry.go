package registry

import (
	"context"
	"fmt"
	"sync"

	"psql-mcp-registry/internal/factory"
	"psql-mcp-registry/internal/model"
	"psql-mcp-registry/internal/pg"
)

type Implementation struct {
	registryMap     map[string]*pg.Client
	instanceStorage InstanceStorage
	clientFactory   factory.ClientFactory
	mu              sync.RWMutex
}

//go:generate mockery --case snake --name InstanceStorage
type InstanceStorage interface {
	ListInstances(ctx context.Context) ([]model.Instance, error)
}

//go:generate mockery --case snake --name Registry
type Registry interface {
	AddInstanceToRegistry(instance model.Instance) error
	GetInstanceClient(instance model.Instance) pg.ClientInterface
}

func NewRegistry(ctx context.Context, instanceStorage InstanceStorage, clientFactory factory.ClientFactory) (Registry, error) {
	instances, err := instanceStorage.ListInstances(ctx)
	if err != nil {
		return nil, fmt.Errorf("instanceStorage.ListInstances: %w", err)
	}

	r := &Implementation{
		registryMap:     make(map[string]*pg.Client),
		instanceStorage: instanceStorage,
		clientFactory:   clientFactory,
	}

	for _, instance := range instances {
		err = r.AddInstanceToRegistry(instance)
		if err != nil {
			continue
		}
	}

	return r, nil
}

func (r *Implementation) GetInstanceClient(instance model.Instance) pg.ClientInterface {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.registryMap[instance.Name]
}

func (r *Implementation) AddInstanceToRegistry(instance model.Instance) error {
	client, err := r.clientFactory.CreateClient(instance)
	if err != nil {
		return fmt.Errorf("failed to create client for instance %s: %w", instance.Name, err)
	}

	// Need to assert back to *pg.Client for internal storage
	// This is safe because our factory returns *pg.Client which implements ClientInterface
	concreteClient, ok := client.(*pg.Client)
	if !ok {
		return fmt.Errorf("client factory returned unexpected type for instance %s", instance.Name)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.registryMap[instance.Name] = concreteClient

	return nil
}
