package mocks

import (
	"github.com/rollout/rox-go/core/model"
	"github.com/stretchr/testify/mock"
)

type FlagRepository struct {
	mock.Mock
}

func (m *FlagRepository) AddFlag(flag interface{}, name string) {
	m.Called(flag, name)
}

func (m *FlagRepository) GetFlag(name string) model.Variant {
	args := m.Called(name)
	return args.Get(0).(model.Variant)
}

func (m *FlagRepository) GetAllFlags() []model.Variant {
	args := m.Called()
	return args.Get(0).([]model.Variant)
}

func (m *FlagRepository) RegisterFlagAddedHandler(handler model.FlagAddedHandler) {
	m.Called(handler)
}
