package roxx

import (
	"fmt"
	"github.com/rollout/rox-go/v4/core/context"
	"regexp"
)

type RegularExpressionExtensions struct {
	parser Parser
}

func NewRegularExpressionExtensions(parser Parser) *RegularExpressionExtensions {
	return &RegularExpressionExtensions{parser: parser}
}

func (e *RegularExpressionExtensions) Extend() {
	e.parser.AddOperator("match", func(p Parser, stack *CoreStack, context context.Context) {
		str, ok1 := stack.Pop().(string)
		pattern, ok2 := stack.Pop().(string)
		flags, ok3 := stack.Pop().(string)

		if !ok1 || !ok2 || !ok3 {
			stack.Push(false)
			return
		}

		if flags != "" {
			pattern = fmt.Sprintf("(?%s)%s", flags, pattern)
		}

		matched, _ := regexp.MatchString(pattern, str)
		stack.Push(matched)
	})
}
