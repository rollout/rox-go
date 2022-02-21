package configuration

import (
	"encoding/json"
	"strconv"
)

type Source int

const (
	SourceCDN Source = iota
	SourceAPI
	SourceRoxy
)

func (cs Source) String() string {
	switch cs {
	case SourceCDN:
		return "CDN"
	case SourceAPI:
		return "API"
	case SourceRoxy:
		return "Roxy"
	}
	return strconv.Itoa(int(cs))
}

type FetchResult struct {
	ParsedData jsonConfiguration
	Source     Source
}

func NewFetchResult(data string, source Source) *FetchResult {
	if data == "" {
		return nil
	}

	var parsedData jsonConfiguration
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		return nil
	}

	return &FetchResult{ParsedData: parsedData, Source: source}
}
