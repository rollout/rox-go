package configuration_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/rollout/rox-go/v6/core/configuration"
	"github.com/rollout/rox-go/v6/core/mocks"
	"github.com/rollout/rox-go/v6/core/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfigurationParserWillReturnNullWhenUnexpectedException(t *testing.T) {
	nestedJson := `
	{
		"application":"12345",
		"targetGroups": [{"condition":"eq(true,true)","_id":"12345"},{"_id":"123456","condition":"eq(true,true)"}],
		"experiments": [
		{"deploymentConfiguration":{"condition":"ifThen(and(true, true)"},"featureFlags":[{"name":"FeatureFlags.isFeatureFlagsEnabled"}],"archived":false,"name":"Feature Flags Drawer Item","_id":"1"},
		{"deploymentConfiguration":{"condition":"ifThen(and(true, true)"},"featureFlags":[{"name":"Invitations.isInvitationsEnabled"}],"archived":false,"name":"Enable Modern Invitations","_id":"2"}]
	}`
	json := mergeNestedAndMasterJson(nestedJson, `
	{
		"nodata": "%s",
		"signature_v0":"K/bEQCkRXa6+uFr5H2jCRCaVgmtsTwbgfrFGVJ9NebfMH8CgOhCDIvF4TM1Vyyl0bGS9a4r4Qgi/g63NDBWk0ZbRrKAUkVG56V3/bI2GDHxFvRNrNbiPmFv/wmLLuwgh1mdzU0EwLG4M7yXoNXtMr6Jli8t4xfBOaWW1g0QpASkiWa7kdTamVip/1QygyUuhX5hOyUMpy4Ny9Hi/QPvVBn6GDMxQtxpLfTavU9cBly2D7Ex8Z7sUUOKeoEJcdsoF1QzH14XvA2HQSICESz7D/uld0PNdG0tMj9NlAZfki8eY2KuUe/53Z0Og5WrqQUxiAdPuJoZr6+kSqlASZrrkYw==",
		"signed_date":"2018-01-09T19:02:00.720Z"
	}`)
	configFetchResult := configuration.NewFetchResult(json, configuration.SourceCDN)

	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("APIKey").Return("12345")

	sf := &mocks.SignatureVerifier{}
	sf.On("Verify", mock.Anything, mock.Anything).Return(true)

	cfi := configuration.NewFetchedInvoker()
	var cfiEvent *model.ConfigurationFetchedArgs
	cfi.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		cfiEvent = e
	})

	cp := configuration.NewParser(sf, nil, cfi)
	conf := cp.Parse(configFetchResult, sdkSettings)

	assert.Nil(t, conf)
	assert.NotNil(t, cfiEvent)
	assert.Equal(t, model.FetcherErrorUnknown, cfiEvent.ErrorDetails)
}

func TestConfigurationParserWillReturnNullWhenWrongSignature(t *testing.T) {
	nestedJson := `
	{
		"application": "12345",
		"targetGroups": [{"condition":"eq(true,true)","_id":"12345"},{"_id":"123456","condition":"eq(true,true)"}],
		"experiments": [
		{"deploymentConfiguration":{"condition":"ifThen(and(true, true)"},"featureFlags":[{"name":"FeatureFlags.isFeatureFlagsEnabled"}],"archived":false,"name":"Feature Flags Drawer Item","_id":"1"},
		{"deploymentConfiguration":{"condition":"ifThen(and(true, true)"},"featureFlags":[{"name":"Invitations.isInvitationsEnabled"}],"archived":false,"name":"Enable Modern Invitations","_id":"2"}]
	}`
	json := mergeNestedAndMasterJson(nestedJson, `
	{
		"data": "%s",
		"signature_v0": "wrongK/bEQCkRXa6+uFr5H2jCRCaVgmtsTwbgfrFGVJ9NebfMH8CgOhCDIvF4TM1Vyyl0bGS9a4r4Qgi/g63NDBWk0ZbRrKAUkVG56V3/bI2GDHxFvRNrNbiPmFv/wmLLuwgh1mdzU0EwLG4M7yXoNXtMr6Jli8t4xfBOaWW1g0QpASkiWa7kdTamVip/1QygyUuhX5hOyUMpy4Ny9Hi/QPvVBn6GDMxQtxpLfTavU9cBly2D7Ex8Z7sUUOKeoEJcdsoF1QzH14XvA2HQSICESz7D/uld0PNdG0tMj9NlAZfki8eY2KuUe/53Z0Og5WrqQUxiAdPuJoZr6+kSqlASZrrkYw==",
		"signed_date":"2018-01-09T19:02:00.720Z"
	}`)
	configFetchResult := configuration.NewFetchResult(json, configuration.SourceAPI)

	errRe := &mocks.ErrorReporter{}

	sf := &mocks.SignatureVerifier{}
	sf.On("Verify", mock.Anything, mock.Anything).Return(false)

	cfi := configuration.NewFetchedInvoker()
	var cfiEvent *model.ConfigurationFetchedArgs
	cfi.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		cfiEvent = e
	})

	cp := configuration.NewParser(sf, errRe, cfi)
	conf := cp.Parse(configFetchResult, nil)

	assert.Nil(t, conf)
	assert.NotNil(t, cfiEvent)
	assert.Equal(t, model.FetcherErrorSignatureVerification, cfiEvent.ErrorDetails)
}

func TestConfigurationParserWillReturnNullWhenWrongApiKey(t *testing.T) {
	nestedJson := `
	{
		"application": "12345",
		"targetGroups": [{"condition":"eq(true,true)","_id":"12345"},{"_id":"123456","condition":"eq(true,true)"}],
		"experiments": [
		{"deploymentConfiguration":{"condition":"ifThen(and(true, true)"},"featureFlags":[{"name":"FeatureFlags.isFeatureFlagsEnabled"}],"archived":false,"name":"Feature Flags Drawer Item","_id":"1"},
		{"deploymentConfiguration":{"condition":"ifThen(and(true, true)"},"featureFlags":[{"name":"Invitations.isInvitationsEnabled"}],"archived":false,"name":"Enable Modern Invitations","_id":"2"}]
	}`
	json := mergeNestedAndMasterJson(nestedJson, `
	{
		"data": "%s",
		"signature_v0":"K/bEQCkRXa6+uFr5H2jCRCaVgmtsTwbgfrFGVJ9NebfMH8CgOhCDIvF4TM1Vyyl0bGS9a4r4Qgi/g63NDBWk0ZbRrKAUkVG56V3/bI2GDHxFvRNrNbiPmFv/wmLLuwgh1mdzU0EwLG4M7yXoNXtMr6Jli8t4xfBOaWW1g0QpASkiWa7kdTamVip/1QygyUuhX5hOyUMpy4Ny9Hi/QPvVBn6GDMxQtxpLfTavU9cBly2D7Ex8Z7sUUOKeoEJcdsoF1QzH14XvA2HQSICESz7D/uld0PNdG0tMj9NlAZfki8eY2KuUe/53Z0Og5WrqQUxiAdPuJoZr6+kSqlASZrrkYw==",
		"signed_date":"2018-01-09T19:02:00.720Z"
	}`)

	configFetchResult := configuration.NewFetchResult(json, configuration.SourceAPI)

	errRe := &mocks.ErrorReporter{}

	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("APIKey").Return("123")

	sf := &mocks.SignatureVerifier{}
	sf.On("Verify", mock.Anything, mock.Anything).Return(true)

	cfi := configuration.NewFetchedInvoker()
	var cfiEvent *model.ConfigurationFetchedArgs
	cfi.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		cfiEvent = e
	})

	cp := configuration.NewParser(sf, errRe, cfi)
	conf := cp.Parse(configFetchResult, sdkSettings)

	assert.Nil(t, conf)
	assert.NotNil(t, cfiEvent)
	assert.Equal(t, model.FetcherErrorMismatchAppKey, cfiEvent.ErrorDetails)
}

func TestConfigurationParserWillParseExperimentsAndTargetGroups(t *testing.T) {
	nestedJson := `
	{
		"application": "12345",
		"targetGroups": [{"condition":"eq(true,true)","_id":"12345"},{"_id":"123456","condition":"eq(true,true)"}],
		"experiments": [
		{"deploymentConfiguration":{"condition":"ifThen(and(true, true)"},"featureFlags":[{"name":"FeatureFlags.isFeatureFlagsEnabled"}],"archived":false,"name":"Feature Flags Drawer Item","_id":"1","labels":["label1"]},
		{"deploymentConfiguration":{"condition":"ifThen(and(true, true)"},"featureFlags":[{"name":"Invitations.isInvitationsEnabled"}],"archived":false,"name":"Enable Modern Invitations","_id":"2"}]
	}
	`
	json := mergeNestedAndMasterJson(nestedJson, `
	{
		"data": "%s",
		"signature_v0":"K/bEQCkRXa6+uFr5H2jCRCaVgmtsTwbgfrFGVJ9NebfMH8CgOhCDIvF4TM1Vyyl0bGS9a4r4Qgi/g63NDBWk0ZbRrKAUkVG56V3/bI2GDHxFvRNrNbiPmFv/wmLLuwgh1mdzU0EwLG4M7yXoNXtMr6Jli8t4xfBOaWW1g0QpASkiWa7kdTamVip/1QygyUuhX5hOyUMpy4Ny9Hi/QPvVBn6GDMxQtxpLfTavU9cBly2D7Ex8Z7sUUOKeoEJcdsoF1QzH14XvA2HQSICESz7D/uld0PNdG0tMj9NlAZfki8eY2KuUe/53Z0Og5WrqQUxiAdPuJoZr6+kSqlASZrrkYw==",
		"signed_date":"2018-01-09T19:02:00.720Z"
	}`)

	configFetchResult := configuration.NewFetchResult(json, configuration.SourceAPI)

	errRe := &mocks.ErrorReporter{}

	sdkSettings := &mocks.SdkSettings{}
	sdkSettings.On("APIKey").Return("12345")

	sf := &mocks.SignatureVerifier{}
	sf.On("Verify", mock.Anything, mock.Anything).Return(true)

	cfi := configuration.NewFetchedInvoker()
	var cfiEvent *model.ConfigurationFetchedArgs
	cfi.RegisterFetchedHandler(func(e *model.ConfigurationFetchedArgs) {
		cfiEvent = e
	})

	cp := configuration.NewParser(sf, errRe, cfi)
	conf := cp.Parse(configFetchResult, sdkSettings)

	assert.NotNil(t, conf)
	assert.Equal(t, 2, len(conf.TargetGroups))
	assert.Equal(t, "12345", conf.TargetGroups[0].ID)
	assert.Equal(t, "eq(true,true)", conf.TargetGroups[0].Condition)
	assert.Equal(t, "123456", conf.TargetGroups[1].ID)
	assert.Equal(t, "eq(true,true)", conf.TargetGroups[1].Condition)

	assert.Equal(t, 2, len(conf.Experiments))
	assert.Equal(t, "ifThen(and(true, true)", conf.Experiments[0].Condition)
	assert.Equal(t, "Feature Flags Drawer Item", conf.Experiments[0].Name)
	assert.Equal(t, "1", conf.Experiments[0].ID)
	assert.False(t, conf.Experiments[0].IsArchived)
	assert.Equal(t, 1, len(conf.Experiments[0].Flags))
	assert.Equal(t, "FeatureFlags.isFeatureFlagsEnabled", conf.Experiments[0].Flags[0])
	assert.Equal(t, 1, len(conf.Experiments[0].Labels))
	assert.Contains(t, conf.Experiments[0].Labels, "label1")
	assert.Equal(t, "ifThen(and(true, true)", conf.Experiments[1].Condition)
	assert.Equal(t, "Enable Modern Invitations", conf.Experiments[1].Name)
	assert.Equal(t, "2", conf.Experiments[1].ID)
	assert.False(t, conf.Experiments[1].IsArchived)
	assert.Equal(t, 1, len(conf.Experiments[1].Flags))
	assert.Equal(t, "Invitations.isInvitationsEnabled", conf.Experiments[1].Flags[0])
	assert.Equal(t, 0, len(conf.Experiments[1].Labels))

	assert.Nil(t, cfiEvent)
}

func mergeNestedAndMasterJson(nestedJson, masterJson string) string {
	nestedJson = strings.Replace(nestedJson, "\n", `\n`, -1)
	nestedJson = strings.Replace(nestedJson, "\t", ` `, -1)
	nestedJson = strings.Replace(nestedJson, `"`, `\"`, -1)
	return fmt.Sprintf(masterJson, nestedJson)
}
