package repositories

import (
	"github.com/rollout/rox-go/v5/core/logging"
	"github.com/rollout/rox-go/v5/core/model"
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

func (r *flagRepository) AddFlag(variant model.Variant, name string, tag string) {
	if variant.(model.Variant).Name() == "" {
		variant.(model.InternalVariant).SetName(name)
	}
	variant.(model.InternalVariant).SetTag(tag)

	r.mutex.Lock()
	r.variants[name] = variant
	r.mutex.Unlock()

	r.raiseFlagAddedEvent(variant)
}

func (r *flagRepository) GetFlag(name string) model.Variant {

	r.mutex.RLock()
	variant, ok := r.variants[name]
	r.mutex.RUnlock()
	if !ok {
		return nil
	}
	return variant.(model.Variant)
}

func (r *flagRepository) GetAllFlags() []model.Variant {
	r.mutex.RLock()
	result := make([]model.Variant, 0, len(r.variants))
	for _, p := range r.variants {
		result = append(result, p.(model.Variant))
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

	defer func() {
		if r := recover(); r != nil {
			logging.GetLogger().Error("Failed to execute flag added handler, panic", r)
		}
	}()

	for _, handler := range handlers {
		handler(flag)
	}
}
