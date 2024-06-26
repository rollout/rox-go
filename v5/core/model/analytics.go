package model

import "time"

type Analytics interface {
	CaptureImpressions([]Impression)
	InitiateIntervalReporting(interval time.Duration)
	StopIntervalReporting()
}

type Impression struct {
	Timestamp         float64     `json:"time"`
	FlagName          string      `json:"flag"`
	Value             interface{} `json:"value"`
	Type              string      `json:"type"`
	Count             float64     `json:"count,omitempty"`
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
