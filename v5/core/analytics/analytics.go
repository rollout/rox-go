package analytics

import (
	"fmt"
	"sync"
	"time"

	"github.com/rollout/rox-go/v5/core/consts"
	"github.com/rollout/rox-go/v5/core/logging"
	"github.com/rollout/rox-go/v5/core/model"
)

type AnalyticsHandler struct {
	uriPath          string
	request          model.Request
	deviceProperties model.DeviceProperties
	impressionsQueue ImpressionsStore
	logger           logging.Logger
	isDisabled       bool
}

type ImpressionsStore struct {
	mu          sync.Mutex
	impressions []model.Impression
}

type AnalyticsDeps struct {
	UriPath           string
	Request           model.Request
	DeviceProperities model.DeviceProperties
	Logger            logging.Logger
	IsDisabled        bool
}

func NewAnalyticsHandler(deps *AnalyticsDeps) model.Analytics {
	if deps.Logger == nil {
		deps.Logger = logging.GetLogger()
	}
	return &AnalyticsHandler{
		uriPath:          deps.UriPath,
		request:          deps.Request,
		deviceProperties: deps.DeviceProperities,
		logger:           deps.Logger,
		impressionsQueue: ImpressionsStore{
			impressions: make([]model.Impression, 0),
		},
		isDisabled: deps.IsDisabled,
	}
}

// InitiateReporting starts the analytics reporting process all
// impressions accumulated over 'interval' time will be sent to the analytics server
func (ah *AnalyticsHandler) InitiateReporting(interval time.Duration) {
	if interval == 0 {
		interval = time.Minute
	}
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			err := ah.postImpressions()
			if err != nil {
				ah.logger.Error("Error posting impressions: %v", err)
			}
		}
	}()
}

func (ah *AnalyticsHandler) postImpressions() error {
	ah.impressionsQueue.mu.Lock()
	defer ah.impressionsQueue.mu.Unlock()

	if len(ah.impressionsQueue.impressions) == 0 {
		return nil
	}

	properties := ah.deviceProperties.GetAllProperties()
	bodyContent := &model.SDKEventBatch{
		AnalyticsVersion: "1.0.0",
		SdkKeyId:         ah.deviceProperties.RolloutKey(),
		Timestamp:        float64(time.Now().Second()),
		Platform:         properties[consts.PropertyTypePlatform.Name],
		SDKVersion:       properties[consts.PropertyTypeLibVersion.Name],
		Events:           ah.impressionsQueue.impressions,
	}

	uri := fmt.Sprintf("%s/impressions/%s", ah.uriPath, bodyContent.SdkKeyId)
	res, err := ah.request.SendPost(uri, bodyContent)

	// accumulated impressions are flushed regardless of reporting success or
	// failure to avoid stack overflow if analytics server is unreachable
	ah.flushImpressions()

	if err != nil {
		return err
	}
	if !res.IsSuccessStatusCode() {
		return fmt.Errorf("Impression reporting failed. Status code: %d. Request: %+v", res.StatusCode, res)
	}

	return nil
}

func (ah *AnalyticsHandler) Enqueue(time float64, flagName string, value interface{}) {
	ah.impressionsQueue.mu.Lock()
	defer ah.impressionsQueue.mu.Unlock()

	ah.impressionsQueue.impressions = append(ah.impressionsQueue.impressions, model.Impression{
		Timestamp: time,
		FlagName:  flagName,
		Value:     value,
	})
}

func (ah *AnalyticsHandler) flushImpressions() {
	ah.impressionsQueue.impressions = make([]model.Impression, 0)
}

func (ah *AnalyticsHandler) IsAnalyticsReportingDisabled() bool {
	return ah.isDisabled
}
