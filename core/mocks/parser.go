package mocks

import (
	"github.com/rollout/rox-go/core/context"
	"github.com/rollout/rox-go/core/roxx"
	"github.com/stretchr/testify/mock"
)

type Parser struct {
	mock.Mock
}

func (m *Parser) GetGlobalContext() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

func (m *Parser) SetGlobalContext(context context.Context) {
	m.Called(context)
}

func (m *Parser) EvaluateExpression(expression string, context context.Context) roxx.EvaluationResult {
	args := m.Called(expression, context)
	return args.Get(0).(roxx.EvaluationResult)
}

func (m *Parser) AddOperator(name string, operation roxx.Operation) {
	m.Called(name, operation)
}
