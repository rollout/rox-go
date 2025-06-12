package model

import "time"

type Analytics interface {
	CaptureImpressions([]Impression)
	InitiateIntervalReporting(interval time.Duration)
	StopIntervalReporting()
}

type Impression struct {
	Timestamp float64     `json:"time"`
	FlagName  string      `json:"flag"`
	Value     interface{} `json:"value"`
	// Optional: telegraf-input-impressions will default to 'IMPRESSION' if not set / empty.
	Type string `json:"type,omitempty"`
	// Optional: telegraf-input-impressions will default to 1 if not set / zero.
	Count float64 `json:"count,omitempty"`
}

type SDKEventBatch struct {
	AnalyticsVersion string       `json:"analyticsVersion"`
	SDKVersion       string       `json:"sdkVersion"`
	Timestamp        float64      `json:"time"`
	Platform         string       `json:"platform"`
	SdkKeyId         string       `json:"rolloutKey"`
	Events           []Impression `json:"events"`
	CountFieldName   string       `json:"countField,omitempty"`
	Origin           string       `json:"origin,omitempty"`
}
