package mocks

import (
	"github.com/rollout/rox-go/v4/core/model"
	"github.com/stretchr/testify/mock"
)

type FlagRepository struct {
	mock.Mock
}

func (m *FlagRepository) AddFlag(variant model.Variant, name string) {
	m.Called(variant, name)
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
