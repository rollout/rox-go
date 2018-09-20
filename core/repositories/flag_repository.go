package repositories

import (
	"github.com/rollout/rox-go/core/model"
	"sync"
)

type flagRepository struct {
	variants map[string]model.Variant
	mutex    sync.RWMutex

	flagAddedHandlers []model.FlagAddedHandler
	handlersMutex     sync.RWMutex
}

func NewFlagRepository() model.FlagRepository {
	return &flagRepository{
		variants: make(map[string]model.Variant),
	}
}

func (r *flagRepository) AddFlag(variant model.Variant, name string) {
	if variant.Name() == "" {
		variant.SetName(name)
	}

	r.mutex.Lock()
	r.variants[name] = variant
	r.mutex.Unlock()

	r.raiseFlagAddedEvent(variant)
}

func (r *flagRepository) GetFlag(name string) model.Variant {
	r.mutex.RLock()
	variant := r.variants[name]
	r.mutex.RUnlock()
	return variant
}

func (r *flagRepository) GetAllFlags() []model.Variant {
	r.mutex.RLock()
	result := make([]model.Variant, 0, len(r.variants))
	for _, p := range r.variants {
		result = append(result, p)
	}
	r.mutex.RUnlock()
	return result
}

func (r *flagRepository) RegisterFlagAddedHandler(handler model.FlagAddedHandler) {
	r.handlersMutex.Lock()
	r.flagAddedHandlers = append(r.flagAddedHandlers, handler)
	r.handlersMutex.Unlock()
}

func (r *flagRepository) raiseFlagAddedEvent(flag model.Variant) {
	r.handlersMutex.RLock()
	handlers := make([]model.FlagAddedHandler, len(r.flagAddedHandlers))
	copy(handlers, r.flagAddedHandlers)
	r.handlersMutex.RUnlock()

	for _, handler := range handlers {
		handler(flag)
	}
}
