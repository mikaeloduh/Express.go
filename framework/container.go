package framework

import (
	"context"
)

type Factory func() any

type ServiceDefinition struct {
	name     string
	factory  Factory
	strategy ScopeStrategy
}

type Container struct {
	services map[string]*ServiceDefinition
}

func NewContainer() *Container {
	return &Container{
		services: make(map[string]*ServiceDefinition),
	}
}

func (c *Container) Register(name string, factory Factory, strategy ScopeStrategy) {
	def := &ServiceDefinition{
		name:     name,
		factory:  factory,
		strategy: strategy,
	}

	def.strategy.Init(def)
	c.services[name] = def
}

func (c *Container) Get(name string) any {
	return c.GetWithContext(context.Background(), name)
}

func (c *Container) GetWithContext(ctx context.Context, name string) any {
	def, exists := c.services[name]
	if !exists {
		return nil
	}
	return def.strategy.Resolve(c, ctx, def)
}

type InstaceKey string

const REQUESTID = InstaceKey("request_id")
