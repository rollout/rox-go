package repositories

import (
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/properties"
	"sync"
)

type customPropertyRepository struct {
	customProperties map[string]*properties.CustomProperty
	mutex            sync.RWMutex

	propertyAddedHandlers []model.CustomPropertyAddedHandler
	handlersMutex         sync.RWMutex
}

func NewCustomPropertyRepository() model.CustomPropertyRepository {
	return &customPropertyRepository{
		customProperties: make(map[string]*properties.CustomProperty),
	}
}

func (r *customPropertyRepository) AddCustomProperty(customProperty *properties.CustomProperty) {
	if customProperty.Name == "" {
		return
	}

	r.mutex.Lock()
	r.customProperties[customProperty.Name] = customProperty
	r.mutex.Unlock()

	r.raisePropertyAddedEvent(customProperty)
}

func (r *customPropertyRepository) AddCustomPropertyIfNotExists(customProperty *properties.CustomProperty) {
	if customProperty.Name == "" {
		return
	}

	r.mutex.Lock()
	var ok bool
	if _, ok = r.customProperties[customProperty.Name]; !ok {
		r.customProperties[customProperty.Name] = customProperty
	}
	r.mutex.Unlock()

	if !ok {
		r.raisePropertyAddedEvent(customProperty)
	}
}

func (r *customPropertyRepository) GetCustomProperty(name string) *properties.CustomProperty {
	r.mutex.RLock()
	property := r.customProperties[name]
	r.mutex.RUnlock()
	return property
}

func (r *customPropertyRepository) GetAllCustomProperties() []*properties.CustomProperty {
	r.mutex.RLock()
	result := make([]*properties.CustomProperty, 0, len(r.customProperties))
	for _, p := range r.customProperties {
		result = append(result, p)
	}
	r.mutex.RUnlock()
	return result
}

func (r *customPropertyRepository) RegisterPropertyAddedHandler(handler model.CustomPropertyAddedHandler) {
	r.handlersMutex.Lock()
	r.propertyAddedHandlers = append(r.propertyAddedHandlers, handler)
	r.handlersMutex.Unlock()
}

func (r *customPropertyRepository) raisePropertyAddedEvent(property *properties.CustomProperty) {
	r.handlersMutex.RLock()
	handlers := make([]model.CustomPropertyAddedHandler, len(r.propertyAddedHandlers))
	copy(handlers, r.propertyAddedHandlers)
	r.handlersMutex.RUnlock()

	for _, handler := range handlers {
		handler(property)
	}
}
