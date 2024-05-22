package analytics

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rollout/rox-go/v5/core/consts"
	"github.com/rollout/rox-go/v5/core/mocks"
	"github.com/rollout/rox-go/v5/core/model"
)

func commonDevicePropertiesMock(sdkKey, platform, libVersion string) *mocks.DeviceProperties {
	deviceProperties := &mocks.DeviceProperties{}

	properties := map[string]string{
		consts.PropertyTypePlatform.Name:   platform,
		consts.PropertyTypeLibVersion.Name: libVersion,
	}
	deviceProperties.On("GetAllProperties").Return(properties)
	deviceProperties.On("RolloutKey").Return(sdkKey)

	return deviceProperties
}

func Test_postImpressions(t *testing.T) {
	testCases := []struct {
		name          string
		msgCount      int
		sendPostCount int
		sendPostErr   error
		resStatusCode int
		errLogCount   int
		expectedErr   bool
	}{
		{
			name:          "Happy path",
			msgCount:      1,
			sendPostCount: 1,
			sendPostErr:   nil,
			resStatusCode: 200,
			errLogCount:   0,
			expectedErr:   false,
		},
		{
			name:          "No messages in queue",
			msgCount:      0,
			sendPostCount: 0,
			sendPostErr:   nil,
			resStatusCode: 200,
			errLogCount:   0,
			expectedErr:   false,
		},
		{
			name:          "Error sending POST request",
			msgCount:      1,
			sendPostCount: 1,
			sendPostErr:   fmt.Errorf("error"),
			resStatusCode: 0,
			errLogCount:   1,
			expectedErr:   true,
		},
		{
			name:          "Non-successs response status code",
			msgCount:      1,
			sendPostCount: 1,
			sendPostErr:   nil,
			resStatusCode: 500,
			errLogCount:   1,
			expectedErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sdkKey := "sdkKey"
			path := "hostPath"
			deviceProperties := commonDevicePropertiesMock(sdkKey, "platform", "libVersion")
			request := &mocks.Request{}
			expectedUri := fmt.Sprintf("%s/impressions/%s", path, sdkKey)
			request.
				On("SendPost", expectedUri, mock.Anything).
				Return(&model.Response{
					StatusCode: tc.resStatusCode,
				}, tc.sendPostErr).
				Times(tc.sendPostCount)
			logger := &mocks.Logger{}
			logger.On("Error", mock.Anything).Times(tc.errLogCount)

			analytics := &AnalyticsHandler{
				uriPath:          path,
				request:          request,
				deviceProperties: deviceProperties,
				logger:           logger,
				impressionsQueue: ImpressionsStore{
					impressions: make([]model.Impression, 0),
				},
			}

			for i := 0; i < tc.msgCount; i++ {
				analytics.Enqueue(float64(time.Now().Second()), "name", "value")
			}
			err := analytics.postImpressions()
			// ensure message queue is flushed after POST request
			analytics.impressionsQueue.mu.Lock()
			assert.Equal(t, 0, len(analytics.impressionsQueue.impressions))
			analytics.impressionsQueue.mu.Unlock()

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
