package factory

import (
	"psql-mcp-registry/internal/model"
	"psql-mcp-registry/internal/pg"
)

//go:generate mockery --case snake --name ClientFactory
type ClientFactory interface {
	CreateClient(instance model.Instance) (pg.ClientInterface, error)
}

//go:generate mockery --case snake --name ConfigLoader
type ConfigLoader interface {
	Load(instanceName string) (*pg.Config, error)
}

type PGClientFactory struct {
	configLoader ConfigLoader
}

func NewPGClientFactory(configLoader ConfigLoader) ClientFactory {
	return &PGClientFactory{
		configLoader: configLoader,
	}
}

func (f *PGClientFactory) CreateClient(instance model.Instance) (pg.ClientInterface, error) {
	config, err := f.configLoader.Load(instance.Name)
	if err != nil {
		return nil, err
	}
	return pg.NewClient(config)
}
