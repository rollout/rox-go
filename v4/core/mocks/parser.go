package mocks

import (
	"github.com/rollout/rox-go/v4/core/context"
	"github.com/rollout/rox-go/v4/core/roxx"
	"github.com/stretchr/testify/mock"
)

type Parser struct {
	mock.Mock
}

func (m *Parser) EvaluateExpression(expression string, context context.Context) roxx.EvaluationResult {
	args := m.Called(expression, context)
	return args.Get(0).(roxx.EvaluationResult)
}

func (m *Parser) AddOperator(name string, operation roxx.Operation) {
	m.Called(name, operation)
}
