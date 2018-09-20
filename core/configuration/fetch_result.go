package configuration

import "strconv"

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
	Data   string
	Source Source
}

func NewFetchResult(data string, source Source) *FetchResult {
	return &FetchResult{Data: data, Source: source}
}
