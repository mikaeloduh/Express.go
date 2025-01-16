package framework

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

type ScopeStrategy interface {
	Init(def *ServiceDefinition)
	Resolve(c *Container, ctx context.Context, def *ServiceDefinition) any
	Cleanup(ctx context.Context)
}

// SingletonScope
type SingletonScopeStrategy struct {
	instance any
}

func (s *SingletonScopeStrategy) Init(def *ServiceDefinition) {
	s.instance = def.factory()
}

func (s *SingletonScopeStrategy) Resolve(_ *Container, _ context.Context, _ *ServiceDefinition) any {
	return s.instance
}

func (s *SingletonScopeStrategy) Cleanup(ctx context.Context) {
	// Singleton scope doesn't need cleanup
}

// PrototypeScope
type PrototypeScopeStrategy struct{}

func (p *PrototypeScopeStrategy) Init(def *ServiceDefinition) {
}

func (p *PrototypeScopeStrategy) Resolve(_ *Container, _ context.Context, def *ServiceDefinition) any {
	return def.factory()
}

func (p *PrototypeScopeStrategy) Cleanup(ctx context.Context) {
	// Prototype scope doesn't need cleanup
}

// HttpRequestScope
type HttpRequestScopeStrategy struct {
	instances sync.Map // key: requestID-serviceName, value: instance
}

func (h *HttpRequestScopeStrategy) Init(def *ServiceDefinition) {
}

func (h *HttpRequestScopeStrategy) Resolve(c *Container, ctx context.Context, def *ServiceDefinition) any {
	requestID := ctx.Value(REQUESTID)
	if requestID == nil {
		return def.factory()
	}

	instanceKey := fmt.Sprintf("%v-%s", requestID, def.name)
	if instance, ok := h.instances.Load(instanceKey); ok {
		return instance
	}

	instance := def.factory()
	h.instances.Store(instanceKey, instance)
	return instance
}

func (h *HttpRequestScopeStrategy) Cleanup(ctx context.Context) {
	requestID, ok := ctx.Value(REQUESTID).(string)
	if ok {
		h.clearRequestInstances(requestID)
	}
}

func (h *HttpRequestScopeStrategy) clearRequestInstances(requestID string) {
	prefix := fmt.Sprintf("%s-", requestID)
	h.instances.Range(func(key, value any) bool {
		if k, ok := key.(string); ok && strings.HasPrefix(k, prefix) {
			h.instances.Delete(k)
		}
		return true
	})
}

// HttpRequestScopeMiddleware is a Middleware that manages request scoped services
func HttpRequestScopeMiddleware(container *Container) Middleware {
	var requestCounter uint64

	return func(w *ResponseWriter, r *Request, next func()) error {
		requestID := atomic.AddUint64(&requestCounter, 1)

		ctx := context.WithValue(r.Context(), REQUESTID, requestID)
		r.Request = r.Request.WithContext(ctx)

		next()

		// clean up all instances associated with this request
		for _, def := range container.services {
			def.strategy.Cleanup(ctx)
		}

		return nil
	}
}
