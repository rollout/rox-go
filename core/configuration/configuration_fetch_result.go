package configuration

import "strconv"

type ConfigurationSource int

const (
	SourceCDN ConfigurationSource = iota
	SourceAPI
	SourceRoxy
)

func (cs ConfigurationSource) String() string {
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

type ConfigurationFetchResult struct {
	Data   string
	Source ConfigurationSource
}

func NewConfigurationFetchResult(data string, source ConfigurationSource) *ConfigurationFetchResult {
	return &ConfigurationFetchResult{Data: data, Source: source}
}
