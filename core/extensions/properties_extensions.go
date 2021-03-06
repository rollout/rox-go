package extensions

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/roxx"
)

type PropertiesExtensions struct {
	parser               roxx.Parser
	propertiesRepository model.CustomPropertyRepository
}

func NewPropertiesExtensions(parser roxx.Parser, propertiesRepository model.CustomPropertyRepository) *PropertiesExtensions {
	return &PropertiesExtensions{
		parser:               parser,
		propertiesRepository: propertiesRepository,
	}
}

func (e *PropertiesExtensions) Extend() {
	e.parser.AddOperator("property", func(p roxx.Parser, stack *roxx.CoreStack, context context.Context) {
		propName := stack.Pop().(string)
		property := e.propertiesRepository.GetCustomProperty(propName)

		if property == nil {
			stack.Push(roxx.TokenTypeUndefined)
		} else {
			value := property.Value(context)
			if value == nil {
				stack.Push(roxx.TokenTypeUndefined)
			} else {
				stack.Push(value)
			}
		}
	})
}
