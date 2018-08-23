package impression

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/model"
	"sync"
)

type impressionInvoker struct {
	internalFlags            model.InternalFlags
	customPropertyRepository model.CustomPropertyRepository
	deviceProperties         model.DeviceProperties
	isRoxy                   bool

	impressionHandlers []model.ImpressionHandler
	handlersMutex      sync.RWMutex
}

func NewImpressionInvoker(internalFlags model.InternalFlags, customPropertyRepository model.CustomPropertyRepository, deviceProperties model.DeviceProperties, isRoxy bool) model.ImpressionInvoker {
	return &impressionInvoker{
		internalFlags:            internalFlags,
		customPropertyRepository: customPropertyRepository,
		deviceProperties:         deviceProperties,
		isRoxy:                   isRoxy,
	}
}

func (ii *impressionInvoker) Invoke(value *model.ReportingValue, experiment *model.Experiment, context context.Context) {
	// TODO Implement analytics logic

	ii.raiseImpressionEvent(model.ImpressionArgs{value, experiment, context})
}

func (ii *impressionInvoker) RegisterImpressionHandler(handler model.ImpressionHandler) {
	ii.handlersMutex.Lock()
	ii.impressionHandlers = append(ii.impressionHandlers, handler)
	ii.handlersMutex.Unlock()
}

func (ii *impressionInvoker) raiseImpressionEvent(args model.ImpressionArgs) {
	ii.handlersMutex.RLock()
	defer ii.handlersMutex.RUnlock()

	for _, handler := range ii.impressionHandlers {
		handler(args)
	}
}
