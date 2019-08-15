package client

import (
	"testing"

	"github.com/rollout/rox-go/core/entities"
	"github.com/rollout/rox-go/core/mocks"
	"github.com/rollout/rox-go/core/model"
	"github.com/rollout/rox-go/core/properties"
	"github.com/stretchr/testify/assert"
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
		"app_key":     "123",
		"api_version": "4.0.0",
		"platform":    "plat",
		"lib_version": "1.5.0",
	})

	buid := NewBUID(sdkSettings, deviceProperties, flagRepo, cpRepo)
	assert.Equal(t, "234A32BB4341EAFD91FC8D0395F4E66F", buid.GetValue())

	deviceProperties2 := &mocks.DeviceProperties{}
	deviceProperties2.On("GetAllProperties").Return(map[string]string{
		"app_key":     "122",
		"api_version": "4.0.0",
		"platform":    "plat",
		"lib_version": "1.5.0",
	})
	buid2 := NewBUID(sdkSettings, deviceProperties2, flagRepo, cpRepo)
	assert.Equal(t, "F5F30C84B8A806E0004043864724A56E", buid2.GetValue())
}

/* TODO move this in StateSender?
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
*/
