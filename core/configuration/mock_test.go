package configuration_test

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"strings"
)

type mockedErrorReporter struct {
	mock.Mock
}

func (m *mockedErrorReporter) Report(message string, err error) {
}

type mockedSignatureVerifier struct {
	mock.Mock
}

func (m *mockedSignatureVerifier) Verify(data, signatureBase64 string) bool {
	args := m.Called(data, signatureBase64)
	return args.Bool(0)
}

type mockedSdkSettings struct {
	mock.Mock
}

func (m mockedSdkSettings) ApiKey() string {
	args := m.Called()
	return args.String(0)
}

func (m mockedSdkSettings) DevModeSecret() string {
	args := m.Called()
	return args.String(0)
}

func mergeNestedAndMasterJson(nestedJson, masterJson string) string {
	nestedJson = strings.Replace(nestedJson, "\n", `\n`, -1)
	nestedJson = strings.Replace(nestedJson, "\t", ` `, -1)
	nestedJson = strings.Replace(nestedJson, `"`, `\"`, -1)
	return fmt.Sprintf(masterJson, nestedJson)
}
