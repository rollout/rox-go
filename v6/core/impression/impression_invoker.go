package impression

import (
	"sync"
	"time"

	"github.com/rollout/rox-go/v6/core/context"
	"github.com/rollout/rox-go/v6/core/logging"
	"github.com/rollout/rox-go/v6/core/model"
)

const impressionEventType = "IMPRESSION"

type impressionInvoker struct {
	internalFlags            model.InternalFlags
	customPropertyRepository model.CustomPropertyRepository
	deviceProperties         model.DeviceProperties
	analytics                model.Analytics
	isRoxy                   bool

	impressionHandlers []model.ImpressionHandler
	handlersMutex      sync.RWMutex
}

type ImpressionsDeps struct {
	InternalFlags            model.InternalFlags
	CustomPropertyRepository model.CustomPropertyRepository
	DeviceProperties         model.DeviceProperties
	Analytics                model.Analytics
	IsRoxy                   bool
}

func NewImpressionInvoker(deps *ImpressionsDeps) model.ImpressionInvoker {
	return &impressionInvoker{
		internalFlags:            deps.InternalFlags,
		customPropertyRepository: deps.CustomPropertyRepository,
		deviceProperties:         deps.DeviceProperties,
		analytics:                deps.Analytics,
		isRoxy:                   deps.IsRoxy,
	}
}

func (ii *impressionInvoker) Invoke(value *model.ReportingValue, context context.Context) {
	if value == nil {
		return
	}

	if ii.analytics != nil && !ii.isRoxy && ii.internalFlags.IsEnabled("rox.internal.analytics") {
		ii.analytics.CaptureImpressions([]model.Impression{{
			DistinctId: ii.deviceProperties.DistinctID(),
			Timestamp:  float64(time.Now().UnixMilli()),
			FlagName:   value.Name,
			Value:      value.Value,
			Type:       impressionEventType,
			Count:      1,
		}})
	}

	ii.raiseImpressionEvent(model.ImpressionArgs{ReportingValue: value, Context: context})
}

func (ii *impressionInvoker) RegisterImpressionHandler(handler model.ImpressionHandler) {
	ii.handlersMutex.Lock()
	ii.impressionHandlers = append(ii.impressionHandlers, handler)
	ii.handlersMutex.Unlock()
}

func (ii *impressionInvoker) raiseImpressionEvent(args model.ImpressionArgs) {
	ii.handlersMutex.RLock()
	handlers := make([]model.ImpressionHandler, len(ii.impressionHandlers))
	copy(handlers, ii.impressionHandlers)
	ii.handlersMutex.RUnlock()

	defer func() {
		if r := recover(); r != nil {
			logging.GetLogger().Error("Failed to execute impression handler, panic", r)
		}
	}()
	for _, handler := range handlers {
		handler(args)
	}
}
