package plugin

import "sync"

type Registry struct {
	steps     map[string]StepExecutor
	resources map[string]ResourceHandler
	mu        sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		steps:     make(map[string]StepExecutor),
		resources: make(map[string]ResourceHandler),
	}
}

func (r *Registry) RegisterStep(name string, executor StepExecutor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.steps[name] = executor
}

func (r *Registry) GetStep(name string) (StepExecutor, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	executor, exists := r.steps[name]
	return executor, exists
}
