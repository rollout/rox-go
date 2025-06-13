package analytics

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rollout/rox-go/v6/core/consts"
	"github.com/rollout/rox-go/v6/core/mocks"
	"github.com/rollout/rox-go/v6/core/model"
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
			expectedUri := fmt.Sprintf("%s/impression/%s", path, sdkKey)
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

			impressions := make([]model.Impression, 0)
			for i := 0; i < tc.msgCount; i++ {
				impressions = append(impressions, model.Impression{
					Timestamp: float64(time.Now().Unix()),
					FlagName:  "name",
					Value:     "value",
				})
			}
			err := analytics.postImpressions(impressions)

			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAnalyticsHandler_CaptureImpressions(t *testing.T) {
	testCases := []struct {
		testName          string
		existingQueueSize int
		numImpressions    int
		maxQueueSize      int
		expectedRequest   int
		responseErr       error
		errorsExpected    int
		expectedQueueSize int
	}{
		{
			testName:          "Queue size met",
			existingQueueSize: 400,
			numImpressions:    100,
			maxQueueSize:      500,
			expectedRequest:   1,
			errorsExpected:    0,
			responseErr:       nil,
			expectedQueueSize: 0,
		},
		{
			testName:          "Queue size exceeded",
			existingQueueSize: 400,
			numImpressions:    500,
			maxQueueSize:      500,
			expectedRequest:   1,
			errorsExpected:    0,
			responseErr:       nil,
			expectedQueueSize: 0,
		},
		{
			testName:          "Queue size not exceeded",
			existingQueueSize: 400,
			numImpressions:    99,
			maxQueueSize:      500,
			expectedRequest:   0,
			errorsExpected:    0,
			responseErr:       nil,
			expectedQueueSize: 499,
		},
		{
			testName:          "Empty queue, size met",
			existingQueueSize: 0,
			numImpressions:    500,
			maxQueueSize:      500,
			expectedRequest:   1,
			errorsExpected:    0,
			responseErr:       nil,
			expectedQueueSize: 0,
		},
		{
			testName:          "Error sending request",
			existingQueueSize: 400,
			numImpressions:    100,
			maxQueueSize:      500,
			expectedRequest:   1,
			errorsExpected:    1,
			responseErr:       fmt.Errorf("error"),
			expectedQueueSize: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			sdkKey := "sdkKey1"
			path := "hostPath1"
			deviceProperties := commonDevicePropertiesMock(sdkKey, "platform", "libVersion")
			request := &mocks.Request{}
			expectedUri := fmt.Sprintf("%s/%s", path, sdkKey)
			request.
				On("SendPost", expectedUri, mock.Anything).
				Return(&model.Response{
					StatusCode: http.StatusOK,
				}, tc.responseErr).
				Times(tc.expectedRequest)

			logger := &mocks.Logger{}
			logger.On("Error", mock.Anything).Times(tc.errorsExpected)

			analytics := &AnalyticsHandler{
				uriPath:          path,
				request:          request,
				deviceProperties: deviceProperties,
				logger:           logger,
				impressionsQueue: ImpressionsStore{
					impressions: make([]model.Impression, tc.existingQueueSize),
				},
				flushAtSize: tc.maxQueueSize,
			}

			analytics.CaptureImpressions(make([]model.Impression, tc.numImpressions))
			if tc.errorsExpected > 0 {
				assert.Error(t, tc.responseErr)
			} else {
				assert.NoError(t, tc.responseErr)
			}
			assert.Equal(t, tc.expectedQueueSize, len(analytics.impressionsQueue.impressions))
		})
	}
}
