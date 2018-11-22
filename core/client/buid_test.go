package client

import (
	"encoding/json"
	"github.com/rollout/rox-go/core/entities"
	"github.com/rollout/rox-go/core/mocks"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/properties"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBUIDWillGenerateCorrectMD5Value(t *testing.T) {
	flag := entities.NewFlag(false)
	flag.(model.InternalVariant).SetName("flag1")

	flagRepo := &mocks.FlagRepository{}
	flagRepo.On("GetAllFlags").Return([]model.Variant{flag})

	cp := properties.NewStringProperty("prop1", "123")

	cpRepo := &mocks.CustomPropertyRepository{}
	cpRepo.On("GetAllCustomProperties").Return([]*properties.CustomProperty{cp})

	sdkSettings := &mocks.SdkSettings{}

	deviceProperties := &mocks.DeviceProperties{}
	deviceProperties.On("GetAllProperties").Return(map[string]string{
		"app_key":        "123",
		"api_version":    "4.0.0",
		"cache_miss_url": "harta",
	})

	buid := NewBUID(sdkSettings, deviceProperties, flagRepo, cpRepo)

	assert.Equal(t, "5512E154362F3127B817C913A3B286CF", buid.GetValue())
}

func TestBUIDWillSerializeFlags(t *testing.T) {
	flag := entities.NewFlag(false)
	flag.(model.InternalVariant).SetName("flag1")

	flagRepo := &mocks.FlagRepository{}
	flagRepo.On("GetAllFlags").Return([]model.Variant{flag})

	b := NewBUID(nil, nil, flagRepo, nil)
	var flags []map[string]interface{}
	err := json.Unmarshal([]byte(b.(*buid).serializeFeatureFlags()), &flags)

	assert.Nil(t, err)

	obj := flags[0]
	options := obj["options"].([]interface{})

	assert.Equal(t, "flag1", obj["name"])
	assert.Equal(t, "false", obj["defaultValue"])
	assert.Equal(t, "false", options[0])
	assert.Equal(t, "true", options[1])
}

func TestBUIDWillSerializeProps(t *testing.T) {
	cp := properties.NewStringProperty("prop1", "123")

	cpRepo := &mocks.CustomPropertyRepository{}
	cpRepo.On("GetAllCustomProperties").Return([]*properties.CustomProperty{cp})

	b := NewBUID(nil, nil, nil, cpRepo)
	var props []map[string]interface{}
	err := json.Unmarshal([]byte(b.(*buid).serializeCustomProperties()), &props)

	assert.Nil(t, err)

	obj := props[0]

	assert.Equal(t, "prop1", obj["name"])
	assert.Equal(t, properties.CustomPropertyTypeString.Type, obj["type"])
	assert.Equal(t, properties.CustomPropertyTypeString.ExternalType, obj["externalType"])
}
