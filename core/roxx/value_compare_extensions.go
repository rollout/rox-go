package roxx

import (
	"github.com/hashicorp/go-version"
	"github.com/rollout/rox-go/core/context"
	"strings"
)

type ValueCompareExtensions struct {
	parser Parser
}

func NewValueCompareExtensions(parser Parser) *ValueCompareExtensions {
	return &ValueCompareExtensions{parser: parser}
}

func (e *ValueCompareExtensions) Extend() {
	e.parser.AddOperator("lt", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop()
		op2 := stack.Pop()

		number1, ok1 := e.toFloat(op1)
		number2, ok2 := e.toFloat(op2)

		if !ok1 || !ok2 {
			stack.Push(false)
		} else {
			stack.Push(number1 < number2)
		}
	})

	e.parser.AddOperator("lte", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop()
		op2 := stack.Pop()

		number1, ok1 := e.toFloat(op1)
		number2, ok2 := e.toFloat(op2)

		if !ok1 || !ok2 {
			stack.Push(false)
		} else {
			stack.Push(number1 <= number2)
		}
	})

	e.parser.AddOperator("gt", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop()
		op2 := stack.Pop()

		number1, ok1 := e.toFloat(op1)
		number2, ok2 := e.toFloat(op2)

		if !ok1 || !ok2 {
			stack.Push(false)
		} else {
			stack.Push(number1 > number2)
		}
	})

	e.parser.AddOperator("gte", func(p Parser, stack *CoreStack, context context.Context) {
		op1 := stack.Pop()
		op2 := stack.Pop()

		number1, ok1 := e.toFloat(op1)
		number2, ok2 := e.toFloat(op2)

		if !ok1 || !ok2 {
			stack.Push(false)
		} else {
			stack.Push(number1 >= number2)
		}
	})

	e.parser.AddOperator("semverNe", func(p Parser, stack *CoreStack, context context.Context) {
		op1, ok1 := stack.Pop().(string)
		op2, ok2 := stack.Pop().(string)

		if !ok1 || !ok2 {
			stack.Push(false)
		} else {
			v1, v2 := e.normalizeVersions(op1, op2)
			version1, err1 := version.NewVersion(v1)
			version2, err2 := version.NewVersion(v2)
			if err1 != nil || err2 != nil {
				stack.Push(false)
			} else {
				stack.Push(!version1.Equal(version2))
			}
		}
	})

	e.parser.AddOperator("semverEq", func(p Parser, stack *CoreStack, context context.Context) {
		op1, ok1 := stack.Pop().(string)
		op2, ok2 := stack.Pop().(string)

		if !ok1 || !ok2 {
			stack.Push(false)
		} else {
			v1, v2 := e.normalizeVersions(op1, op2)
			version1, err1 := version.NewVersion(v1)
			version2, err2 := version.NewVersion(v2)
			if err1 != nil || err2 != nil {
				stack.Push(false)
			} else {
				stack.Push(version1.Equal(version2))
			}
		}
	})

	e.parser.AddOperator("semverLt", func(p Parser, stack *CoreStack, context context.Context) {
		op1, ok1 := stack.Pop().(string)
		op2, ok2 := stack.Pop().(string)

		if !ok1 || !ok2 {
			stack.Push(false)
		} else {
			v1, v2 := e.normalizeVersions(op1, op2)
			version1, err1 := version.NewVersion(v1)
			version2, err2 := version.NewVersion(v2)
			if err1 != nil || err2 != nil {
				stack.Push(false)
			} else {
				stack.Push(version1.LessThan(version2))
			}
		}
	})

	e.parser.AddOperator("semverLte", func(p Parser, stack *CoreStack, context context.Context) {
		op1, ok1 := stack.Pop().(string)
		op2, ok2 := stack.Pop().(string)

		if !ok1 || !ok2 {
			stack.Push(false)
		} else {
			v1, v2 := e.normalizeVersions(op1, op2)
			version1, err1 := version.NewVersion(v1)
			version2, err2 := version.NewVersion(v2)
			if err1 != nil || err2 != nil {
				stack.Push(false)
			} else {
				stack.Push(version1.LessThan(version2) || version1.Equal(version2))
			}
		}
	})

	e.parser.AddOperator("semverGt", func(p Parser, stack *CoreStack, context context.Context) {
		op1, ok1 := stack.Pop().(string)
		op2, ok2 := stack.Pop().(string)

		if !ok1 || !ok2 {
			stack.Push(false)
		} else {
			v1, v2 := e.normalizeVersions(op1, op2)
			version1, err1 := version.NewVersion(v1)
			version2, err2 := version.NewVersion(v2)
			if err1 != nil || err2 != nil {
				stack.Push(false)
			} else {
				stack.Push(version1.GreaterThan(version2))
			}
		}
	})

	e.parser.AddOperator("semverGte", func(p Parser, stack *CoreStack, context context.Context) {
		op1, ok1 := stack.Pop().(string)
		op2, ok2 := stack.Pop().(string)

		if !ok1 || !ok2 {
			stack.Push(false)
		} else {
			v1, v2 := e.normalizeVersions(op1, op2)
			version1, err1 := version.NewVersion(v1)
			version2, err2 := version.NewVersion(v2)
			if err1 != nil || err2 != nil {
				stack.Push(false)
			} else {
				stack.Push(version1.GreaterThan(version2) || version1.Equal(version2))
			}
		}
	})
}

func (e *ValueCompareExtensions) toFloat(value interface{}) (float64, bool) {
	if value, ok := value.(float64); ok {
		return value, true
	}
	if value, ok := value.(int); ok {
		return float64(value), true
	}
	return 0, false
}

func (e *ValueCompareExtensions) normalizeVersions(version1, version2 string) (string, string) {
	// Package github.com/hashicorp/go-version treats versions "1.1", "1.1.0", "1.1.0.0", etc as equal.
	// Fix this behavior.

	segmentCount1 := len(strings.Split(version1, "."))
	segmentCount2 := len(strings.Split(version2, "."))

	if segmentCount1 == segmentCount2 {
		return version1, version2
	}

	if segmentCount1 < segmentCount2 {
		i := segmentCount2 - segmentCount1
		for i > 0 {
			version1 += ".0"
			i--
		}
		version1 += ".0"
		version2 += ".1"
		return version1, version2
	} else {
		i := segmentCount1 - segmentCount2
		for i > 0 {
			version2 += ".0"
			i--
		}
		version1 += ".1"
		version2 += ".0"
		return version1, version2
	}
}
