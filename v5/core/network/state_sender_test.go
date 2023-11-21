package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/rollout/rox-go/v5/core/client"
	"github.com/rollout/rox-go/v5/core/consts"
	"github.com/rollout/rox-go/v5/core/entities"
	"github.com/rollout/rox-go/v5/core/mocks"
	"github.com/rollout/rox-go/v5/core/model"
	"github.com/rollout/rox-go/v5/core/properties"
	"github.com/rollout/rox-go/v5/core/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var appKey string = "123"

func createNewDeviceProp() map[string]string {
	return map[string]string{
		"platform":      "Go",
		"devModeSecret": "shh...",
		"app_key":       appKey,
		"api_version":   "4.0.0",
	}
}

func TestWillSerializeFlags(t *testing.T) {
	request := &mocks.Request{}
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	flag := entities.NewFlag(false)
	flag.(model.InternalVariant).SetName("flag1")

	flagRepo := &mocks.FlagRepository{}
	flagRepo.On("GetAllFlags").Return([]model.Variant{flag})
	flagRepo.On("RegisterFlagAddedHandler", mock.Anything).Return()
	cpRepo := repositories.NewCustomPropertyRepository()
	environment := client.NewSaasEnvironment()

	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, false)

	serializedFlags, featureFlags := stateSender.serializeFeatureFlags()
	var flags []map[string]interface{}
	err := json.Unmarshal([]byte(serializedFlags), &flags)

	assert.Nil(t, err)

	obj := flags[0]
	options := obj["options"].([]interface{})

	assert.Equal(t, "flag1", obj["name"])
	assert.Equal(t, "false", obj["defaultValue"])
	assert.Equal(t, "false", options[0])
	assert.Equal(t, "true", options[1])

	assert.Equal(t, obj["name"], featureFlags[0].Name)
	assert.Equal(t, obj["defaultValue"], featureFlags[0].DefaultValue)
	assert.Equal(t, options[0], featureFlags[0].Options[0])
	assert.Equal(t, options[1], featureFlags[0].Options[1])
}

func TestWillSerializeFlagsNewPlatformFormat(t *testing.T) {
	request := &mocks.Request{}
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	flagBln := entities.NewFlag(false)
	flagBln.(model.InternalVariant).SetName("flagBln")

	flagStr := entities.NewRoxString("o1", []string {"o2"})
	flagStr.(model.InternalVariant).SetName("flagStr")

	flagInt := entities.NewRoxInt(0, []int {1, 2})
	flagInt.(model.InternalVariant).SetName("flagInt")

	flagFlt := entities.NewRoxDouble(3.14, []float64 {2.71})
	flagFlt.(model.InternalVariant).SetName("flagFlt")

	flagRepo := &mocks.FlagRepository{}
	flagRepo.On("GetAllFlags").Return([]model.Variant{flagStr, flagBln, flagInt, flagFlt})
	flagRepo.On("RegisterFlagAddedHandler", mock.Anything).Return()
	cpRepo := repositories.NewCustomPropertyRepository()
	environment := client.NewSaasEnvironment()

	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, true)

	serializedFlags, featureFlags := stateSender.serializeFeatureFlags()
	var flags []map[string]interface{}
	err := json.Unmarshal([]byte(serializedFlags), &flags)

	assert.Nil(t, err)

	assert.Equal(t, len(featureFlags), 4)
	assert.Equal(t, len(flags), 4)

	obj := flags[0]
	assert.Equal(t, "flagBln", obj["name"])
	assert.Equal(t, "false", obj["defaultValue"])
	assert.Equal(t, "Boolean", obj["externalType"])
	options := obj["options"].([]interface{})
	assert.Equal(t, "false", options[0])
	assert.Equal(t, "true", options[1])

	assert.Equal(t, "flagBln", featureFlags[0].Name)
	assert.Equal(t, "false", featureFlags[0].DefaultValue)
	assert.Equal(t, "Boolean", featureFlags[0].ExternalType)
	assert.Equal(t, "false", featureFlags[0].Options[0])
	assert.Equal(t, "true", featureFlags[0].Options[1])

	// number options should probably be strings, like in other SDKs
	obj2 := flags[1]
	assert.Equal(t, "flagFlt", obj2["name"])
	assert.Equal(t, 3.14, obj2["defaultValue"])
	assert.Equal(t, "Number", obj2["externalType"])
	options2 := obj2["options"].([]interface{})
	assert.Equal(t, 2.71, options2[0])
	assert.Equal(t, 3.14, options2[1])

	assert.Equal(t, "flagFlt", featureFlags[1].Name)
	assert.Equal(t, 3.14, featureFlags[1].DefaultValue)
	assert.Equal(t, "Number", featureFlags[1].ExternalType)
	assert.Equal(t, 2.71, featureFlags[1].Options[0])
	assert.Equal(t, 3.14, featureFlags[1].Options[1])

	obj3 := flags[2]
	assert.Equal(t, "flagInt", obj3["name"])
	assert.Equal(t, float64(0), obj3["defaultValue"])
	assert.Equal(t, "Number", obj3["externalType"])
	options3 := obj3["options"].([]interface{})
	assert.Equal(t, float64(1), options3[0])
	assert.Equal(t, float64(2), options3[1])
	assert.Equal(t, float64(0), options3[2])

	assert.Equal(t, "flagInt", featureFlags[2].Name)
	assert.Equal(t, int(0), featureFlags[2].DefaultValue)
	assert.Equal(t, "Number", featureFlags[2].ExternalType)
	assert.Equal(t, int64(1), featureFlags[2].Options[0])
	assert.Equal(t, int64(2), featureFlags[2].Options[1])
	assert.Equal(t, int64(0), featureFlags[2].Options[2])

	obj4 := flags[3]
	assert.Equal(t, "flagStr", obj4["name"])
	assert.Equal(t, "o1", obj4["defaultValue"])
	assert.Equal(t, "String", obj4["externalType"])
	options4 := obj4["options"].([]interface{})
	assert.Equal(t, "o2", options4[0])
	assert.Equal(t, "o1", options4[1])

	assert.Equal(t, "flagStr", featureFlags[3].Name)
	assert.Equal(t, "o1", featureFlags[3].DefaultValue)
	assert.Equal(t, "String", featureFlags[3].ExternalType)
	assert.Equal(t, "o2", featureFlags[3].Options[0])
	assert.Equal(t, "o1", featureFlags[3].Options[1])
}

func TestWillSerializeProps(t *testing.T) {
	request := &mocks.Request{}
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	cp := properties.NewStringProperty("prop1", "123")
	flagRepo := repositories.NewFlagRepository()
	cpRepo := &mocks.CustomPropertyRepository{}
	cpRepo.On("GetAllCustomProperties").Return([]*properties.CustomProperty{cp})
	cpRepo.On("RegisterPropertyAddedHandler", mock.Anything).Return()
	environment := client.NewSaasEnvironment()

	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, false)

	var props []map[string]interface{}
	serializedCustomProperties, customProperties := stateSender.serializeCustomProperties()
	err := json.Unmarshal([]byte(serializedCustomProperties), &props)

	assert.Nil(t, err)

	obj := props[0]

	assert.Equal(t, "prop1", obj["name"])
	assert.Equal(t, properties.CustomPropertyTypeString.Type, obj["type"])
	assert.Equal(t, properties.CustomPropertyTypeString.ExternalType, obj["externalType"])

	assert.Equal(t, obj["name"], customProperties[0].Name)
	assert.Equal(t, obj["type"], customProperties[0].Type)
	assert.Equal(t, obj["externalType"], customProperties[0].ExternalType)
}

func TestWillSerializePropsInNewPlatformFormat(t *testing.T) {
	request := &mocks.Request{}
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	cpStr := properties.NewStringProperty("propStr", "123")
	cpBln := properties.NewBooleanProperty("propBln", true)
	cpFlt := properties.NewFloatProperty("propFlt", 123.4)
	cpInt := properties.NewIntegerProperty("propInt", 123)
	cpSmv := properties.NewSemverProperty("propSmv", "1.2.3")
	cpDt := properties.NewTimeProperty("propDt", time.Now())
	flagRepo := repositories.NewFlagRepository()
	cpRepo := &mocks.CustomPropertyRepository{}
	cpRepo.On("GetAllCustomProperties").Return([]*properties.CustomProperty{cpStr, cpBln, cpFlt, cpInt, cpSmv, cpDt})
	cpRepo.On("RegisterPropertyAddedHandler", mock.Anything).Return()
	environment := client.NewSaasEnvironment()

	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, true)

	var props []map[string]interface{}
	serializedCustomProperties, customProperties := stateSender.serializeCustomProperties()
	err := json.Unmarshal([]byte(serializedCustomProperties), &props)

	assert.Nil(t, err)

	assert.Equal(t, len(props), 6)

	assert.Equal(t, "propBln", props[0]["name"])
	assert.Equal(t, "Boolean", props[0]["externalType"])
	assert.Equal(t, "propBln", customProperties[0].Name)
	assert.Equal(t, "Boolean", customProperties[0].ExternalType)

	assert.Equal(t, "propDt", props[1]["name"])
	assert.Equal(t, "DateTime", props[1]["externalType"])
	assert.Equal(t, "propDt", customProperties[1].Name)
	assert.Equal(t, "DateTime", customProperties[1].ExternalType)

	assert.Equal(t, "propFlt", props[2]["name"])
	assert.Equal(t, "Number", props[2]["externalType"])
	assert.Equal(t, "propFlt", customProperties[2].Name)
	assert.Equal(t, "Number", customProperties[2].ExternalType)

	assert.Equal(t, "propInt", props[3]["name"])
	assert.Equal(t, "Number", props[3]["externalType"])
	assert.Equal(t, "propInt", customProperties[3].Name)
	assert.Equal(t, "Number", customProperties[3].ExternalType)

	assert.Equal(t, "propSmv", props[4]["name"])
	assert.Equal(t, "Semver", props[4]["externalType"])
	assert.Equal(t, "propSmv", customProperties[4].Name)
	assert.Equal(t, "Semver", customProperties[4].ExternalType)

	assert.Equal(t, "propStr", props[5]["name"])
	assert.Equal(t, "String", props[5]["externalType"])
	assert.Equal(t, "propStr", customProperties[5].Name)
	assert.Equal(t, "String", customProperties[5].ExternalType)
}

func TestWillCallOnlyCDNStateMD5ChangesForFlag(t *testing.T) {
	request := &mocks.Request{}
	flagRepo := repositories.NewFlagRepository()
	cpRepo := repositories.NewCustomPropertyRepository()
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	environment := client.NewSaasEnvironment()

	var requestData model.RequestData
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("{\"result\": 200}")}
	request.On("SendGet", mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		requestData = args.Get(0).(model.RequestData)
	})

	flagRepo.AddFlag(entities.NewFlag(false), "flag1")
	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, false)
	stateSender.Send()

	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "C1C65A5AC8A732EAB7FCD81017BF5A87"), requestData.URL)
	request.AssertNumberOfCalls(t, "SendGet", 1)

	flagRepo.AddFlag(entities.NewFlag(false), "flag2")
	stateSender.Send()
	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "F367809AB0CCA5A05EA9DFB3C4E9E15C"), requestData.URL)
}

func TestWillCallOnlyCDNStateMD5ChangesForCustomProperty(t *testing.T) {
	request := &mocks.Request{}
	flagRepo := repositories.NewFlagRepository()
	cpRepo := repositories.NewCustomPropertyRepository()
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	environment := client.NewSaasEnvironment()

	var requestData model.RequestData
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("{\"result\": 200}")}
	request.On("SendGet", mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		requestData = args.Get(0).(model.RequestData)
	})

	cpRepo.AddCustomProperty(properties.NewStringProperty("cp1", "true"))
	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, false)
	stateSender.Send()

	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "02338C470874941BEB8335F76A0F0FBB"), requestData.URL)
	request.AssertNumberOfCalls(t, "SendGet", 1)

	cpRepo.AddCustomProperty(properties.NewFloatProperty("cp2", 20))
	stateSender.Send()
	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "AE3A6DCB39C8306E854CB682548020F1"), requestData.URL)
}

func TestWillCallOnlyCDNStateMD5FlagOrderNotImportant(t *testing.T) {
	request := &mocks.Request{}
	flagRepo := repositories.NewFlagRepository()
	cpRepo := repositories.NewCustomPropertyRepository()
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	environment := client.NewSaasEnvironment()

	var requestData model.RequestData
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("{\"result\": 200}")}
	request.On("SendGet", mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		requestData = args.Get(0).(model.RequestData)
	})

	flagRepo.AddFlag(entities.NewFlag(false), "flag1")
	flagRepo.AddFlag(entities.NewFlag(false), "flag2")
	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, false)
	stateSender.Send()

	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "F367809AB0CCA5A05EA9DFB3C4E9E15C"), requestData.URL)
	request.AssertNumberOfCalls(t, "SendGet", 1)

	flagRepo2 := repositories.NewFlagRepository()
	flagRepo2.AddFlag(entities.NewFlag(false), "flag2")
	flagRepo2.AddFlag(entities.NewFlag(false), "flag1")
	stateSender = NewStateSender(request, dp, flagRepo2, cpRepo, environment, false)
	stateSender.Send()
	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "F367809AB0CCA5A05EA9DFB3C4E9E15C"), requestData.URL)
}

func TestWillCallOnlyCDNStateMD5CustomPropertyOrderNotImportant(t *testing.T) {
	request := &mocks.Request{}
	flagRepo := repositories.NewFlagRepository()
	cpRepo := repositories.NewCustomPropertyRepository()
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	environment := client.NewSaasEnvironment()

	var requestData model.RequestData
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("{\"result\": 200}")}
	request.On("SendGet", mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		requestData = args.Get(0).(model.RequestData)
	})

	cpRepo.AddCustomProperty(properties.NewStringProperty("cp1", "1111"))
	cpRepo.AddCustomProperty(properties.NewStringProperty("cp2", "2222"))
	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, false)
	stateSender.Send()

	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "8BB417F48703DDBD07EC0C2F2160B4B2"), requestData.URL)
	request.AssertNumberOfCalls(t, "SendGet", 1)

	cpRepo2 := repositories.NewCustomPropertyRepository()
	cpRepo2.AddCustomProperty(properties.NewStringProperty("cp2", "2222"))
	cpRepo2.AddCustomProperty(properties.NewStringProperty("cp1", "1111"))
	stateSender.Send()
	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "8BB417F48703DDBD07EC0C2F2160B4B2"), requestData.URL)
}

func TestWillReturnNullWhenCDNFailsWithException(t *testing.T) {
	request := &mocks.Request{}
	flagRepo := repositories.NewFlagRepository()
	cpRepo := repositories.NewCustomPropertyRepository()
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	environment := client.NewSaasEnvironment()

	var reqCDNData model.RequestData
	var reqAPIData model.RequestData
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("{\"result\": 200}")}
	request.On("SendGet", mock.Anything).Return(response, errors.New("not found")).Run(func(args mock.Arguments) {
		reqCDNData = args.Get(0).(model.RequestData)
	})
	request.On("SendPost", mock.Anything, mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		reqAPIData = args.Get(0).(model.RequestData)
	})

	flagRepo.AddFlag(entities.NewFlag(false), "flag")
	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, false)
	stateSender.Send()

	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "00C4910E8BA69D08C65D05849C9E6DFB"), reqCDNData.URL)
	request.AssertNumberOfCalls(t, "SendGet", 1)
	assert.Equal(t, "", reqAPIData.URL)
	request.AssertNumberOfCalls(t, "SendPost", 0)
}

func TestWillReturnNullWhenCDNSucceedWithEmptyResponse(t *testing.T) {
	request := &mocks.Request{}
	flagRepo := repositories.NewFlagRepository()
	cpRepo := repositories.NewCustomPropertyRepository()
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	environment := client.NewSaasEnvironment()

	var reqCDNData model.RequestData
	var reqAPIData string
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("")}
	request.On("SendGet", mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		reqCDNData = args.Get(0).(model.RequestData)
	})
	request.On("SendPost", mock.Anything, mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		reqAPIData = args.Get(0).(string)
	})

	flagRepo.AddFlag(entities.NewFlag(false), "flag")
	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, false)
	stateSender.Send()

	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "00C4910E8BA69D08C65D05849C9E6DFB"), reqCDNData.URL)
	request.AssertNumberOfCalls(t, "SendGet", 1)
	assert.Equal(t, "", reqAPIData)
	request.AssertNumberOfCalls(t, "SendPost", 0)
}

func TestWillReturnNullWhenCDNSucceedWithNotJsonResponse(t *testing.T) {
	request := &mocks.Request{}
	flagRepo := repositories.NewFlagRepository()
	cpRepo := repositories.NewCustomPropertyRepository()
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	environment := client.NewSaasEnvironment()

	var reqCDNData model.RequestData
	var reqAPIData string
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("fdsafdas{")}
	request.On("SendGet", mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		reqCDNData = args.Get(0).(model.RequestData)
	})
	request.On("SendPost", mock.Anything, mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		reqAPIData = args.Get(0).(string)
	})

	flagRepo.AddFlag(entities.NewFlag(false), "flag")
	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, false)
	stateSender.Send()

	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "00C4910E8BA69D08C65D05849C9E6DFB"), reqCDNData.URL)
	request.AssertNumberOfCalls(t, "SendGet", 1)
	assert.Equal(t, "", reqAPIData)
	request.AssertNumberOfCalls(t, "SendPost", 0)
}

func TestWillReturnNullWhenCDNFails404APIWithException(t *testing.T) {
	request := &mocks.Request{}
	flagRepo := repositories.NewFlagRepository()
	cpRepo := repositories.NewCustomPropertyRepository()
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	environment := client.NewSaasEnvironment()

	var reqCDNData model.RequestData
	var reqAPIData string
	response := &model.Response{StatusCode: http.StatusNotFound, Content: []byte("")}
	request.On("SendGet", mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		reqCDNData = args.Get(0).(model.RequestData)
	})
	request.On("SendPost", mock.Anything, mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		reqAPIData = args.Get(0).(string)
	})

	flagRepo.AddFlag(entities.NewFlag(false), "flag")
	cpRepo.AddCustomProperty(properties.NewStringProperty("id", "1111"))
	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, false)
	stateSender.Send()

	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "996ABD4ED5D9D4DF02E56C39ED1F701C"), reqCDNData.URL)
	request.AssertNumberOfCalls(t, "SendGet", 1)
	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateAPIPath(), appKey, "996ABD4ED5D9D4DF02E56C39ED1F701C"), reqAPIData)
	request.AssertNumberOfCalls(t, "SendPost", 1)
}

func TestWillReturnAPIDataWhenCDNSucceedWithResult404APIOK(t *testing.T) {
	request := &mocks.Request{}
	flagRepo := repositories.NewFlagRepository()
	cpRepo := repositories.NewCustomPropertyRepository()
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	environment := client.NewSaasEnvironment()

	var reqCDNData model.RequestData
	var reqAPIData string
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("{\"result\": 404}")}
	request.On("SendGet", mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		reqCDNData = args.Get(0).(model.RequestData)
	})
	request.On("SendPost", mock.Anything, mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		reqAPIData = args.Get(0).(string)
	})

	flagRepo.AddFlag(entities.NewFlag(false), "flag")
	cpRepo.AddCustomProperty(properties.NewStringProperty("id", "1111"))
	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, false)
	stateSender.Send()

	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateCDNPath(), appKey, "996ABD4ED5D9D4DF02E56C39ED1F701C"), reqCDNData.URL)
	request.AssertNumberOfCalls(t, "SendGet", 1)
	assert.Equal(t, fmt.Sprintf("%s/%s/%s", consts.EnvironmentStateAPIPath(), appKey, "996ABD4ED5D9D4DF02E56C39ED1F701C"), reqAPIData)
	request.AssertNumberOfCalls(t, "SendPost", 1)
}

func TestWillReturnAPIDataWhenWhenSelfManaged(t *testing.T) {
	request := &mocks.Request{}
	flagRepo := repositories.NewFlagRepository()
	cpRepo := repositories.NewCustomPropertyRepository()
	dp := &mocks.DeviceProperties{}
	dp.On("GetAllProperties").Return(createNewDeviceProp())
	environment := client.NewSelfManagedEnvironment(client.NewSelfManagedOptions(
		client.SelfManagedOptionsBuilder{
			ServerURL:    "http://harta2.com",
			AnalyticsURL: "http://harta2.com",
		}))

	var reqAPIData string
	response := &model.Response{StatusCode: http.StatusOK, Content: []byte("{\"result\": 404}")}
	request.On("SendPost", mock.Anything, mock.Anything).Return(response, nil).Run(func(args mock.Arguments) {
		reqAPIData = args.Get(0).(string)
	})

	flagRepo.AddFlag(entities.NewFlag(false), "flag")
	cpRepo.AddCustomProperty(properties.NewStringProperty("id", "1111"))
	stateSender := NewStateSender(request, dp, flagRepo, cpRepo, environment, false)
	stateSender.Send()

	assert.Equal(t, fmt.Sprintf("%s/%s/%s", "http://harta2.com/device/update_state_store", appKey, "996ABD4ED5D9D4DF02E56C39ED1F701C"), reqAPIData)
	request.AssertNumberOfCalls(t, "SendPost", 1)
}

