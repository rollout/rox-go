package repositories

import (
	"github.com/rollout/rox-go/core/model"
	"sync"
)

type flagRepository struct {
	variants map[string]interface{}
	//variants map[string]model.Variant
	stringFlags map[string]model.RoxString
	intFlags map[string]model.RoxInt
	doubleFlags map[string]model.RoxDouble
	boolFlags map[string]model.Flag
	mutex    sync.RWMutex

	flagAddedHandlers []model.FlagAddedHandler
	handlersMutex     sync.RWMutex
}

func NewFlagRepository() model.FlagRepository {
	return &flagRepository{
		variants: make(map[string]interface{}),
		stringFlags: make(map[string]model.RoxString),
		intFlags: make(map[string]model.RoxInt),
		doubleFlags: make(map[string]model.RoxDouble),
		boolFlags:make(map[string]model.Flag),
	}
}

func (r *flagRepository) AddFlag(variant interface{}, name string) {

	_, isVariant := variant.(model.Variant)
	if !isVariant {
		return
	}

	if variant.(model.Variant).Name() == "" {
		variant.(model.InternalVariant).SetName(name)
	}

	//switch roxFlag.(type) {
	//case model.RoxString:
	//	r.mutex.Lock()
	//	r.stringFlags[name] = variant.(model.RoxString)
	//	r.mutex.Unlock()
	//case model.RoxInt:
	//	r.mutex.Lock()
	//	r.intFlags[name] = variant.(model.RoxInt)
	//	r.mutex.Unlock()
	//case model.RoxDouble:
	//	r.mutex.Lock()
	//	r.doubleFlags[name] = variant.(model.RoxDouble)
	//	r.mutex.Unlock()
	//case model.Flag:
	//	r.mutex.Lock()
	//	r.variants[name] = variant.(model.Flag)
	//	r.mutex.Unlock()
	//}

	switch variant.(type) {
	case model.RoxString:
		r.mutex.Lock()
		r.variants[name] = variant.(model.RoxString)
		r.mutex.Unlock()
	case model.RoxInt:
		r.mutex.Lock()
		r.variants[name] = variant.(model.RoxInt)
		r.mutex.Unlock()
	case model.RoxDouble:
		r.mutex.Lock()
		r.variants[name] = variant.(model.RoxDouble)
		r.mutex.Unlock()
	case model.Flag:
		r.mutex.Lock()
		r.variants[name] = variant.(model.Flag)
		r.mutex.Unlock()
	}

	r.raiseFlagAddedEvent(variant.(model.Variant))
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
	result := make([]model.Variant, 0, len(r.variants)+len(r.stringFlags)+len(r.intFlags)+len(r.boolFlags)+len(r.doubleFlags))
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

	for _, handler := range handlers {
		handler(flag)
	}
}
