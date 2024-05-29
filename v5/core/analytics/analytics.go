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
	flushAtSize      int
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
	FlushAtSize       int
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
		isDisabled:  deps.IsDisabled,
		flushAtSize: deps.FlushAtSize | 500,
	}
}

// InitiateReporting starts the analytics reporting process all
// impressions accumulated over 'interval' time will be sent to the analytics server
func (ah *AnalyticsHandler) InitiateIntervalReporting(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			// extract current impressions and flush the queue
			ah.impressionsQueue.mu.Lock()
			extractedImpressions := ah.impressionsQueue.impressions
			ah.impressionsQueue.impressions = make([]model.Impression, 0)
			ah.impressionsQueue.mu.Unlock()

			if len(extractedImpressions) > 0 {
				err := ah.postImpressions(extractedImpressions)
				if err != nil {
					// don't requeue to avoid stack overflow if analytics server is unreachable
					ah.logger.Error("Error posting impressions: %v", err)
				}
			}
		}
	}()
}

// CaptureImpressions adds a new impression to the queue and will report the
// impressions if max queues size will be exceeded
func (ah *AnalyticsHandler) CaptureImpressions(newImpressions []model.Impression) {
	ah.impressionsQueue.mu.Lock()
	totalImpressions := append(ah.impressionsQueue.impressions, newImpressions...)
	metFlushSize := len(totalImpressions) >= ah.flushAtSize
	if metFlushSize {
		ah.impressionsQueue.impressions = make([]model.Impression, 0)
	} else {
		ah.impressionsQueue.impressions = totalImpressions
	}
	ah.impressionsQueue.mu.Unlock()

	if metFlushSize {
		go func() {
			err := ah.postImpressions(totalImpressions)
			if err != nil {
				// don't requeue to avoid stack overflow if analytics server is unreachable
				ah.logger.Error("Error posting full queue of impressions due to http error, impressions data lost", err)
			}
		}()
	}
}

func (ah *AnalyticsHandler) postImpressions(impressions []model.Impression) error {
	properties := ah.deviceProperties.GetAllProperties()
	bodyContent := &model.SDKEventBatch{
		AnalyticsVersion: "1.0.0",
		SdkKeyId:         ah.deviceProperties.RolloutKey(),
		Timestamp:        float64(time.Now().Unix()),
		Platform:         properties[consts.PropertyTypePlatform.Name],
		SDKVersion:       properties[consts.PropertyTypeLibVersion.Name],
		Events:           impressions,
	}

	uri := fmt.Sprintf("%s/%s", ah.uriPath, bodyContent.SdkKeyId)
	res, err := ah.request.SendPost(uri, bodyContent)

	if err != nil {
		return err
	}
	if !res.IsSuccessStatusCode() {
		return fmt.Errorf("Impression reporting failed. Status code: %d. Request: %+v", res.StatusCode, res)
	}

	return nil
}

func (ah *AnalyticsHandler) IsAnalyticsReportingDisabled() bool {
	return ah.isDisabled
}
