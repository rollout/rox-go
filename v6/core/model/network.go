package model

import "net/url"

type Request interface {
	SendGet(requestData RequestData) (response *Response, err error)
	SendPost(uri string, content interface{}) (response *Response, err error)
}

type RequestData struct {
	URL         string
	QueryParams map[string]string
}

type Response struct {
	StatusCode int
	Content    []byte
}

func (r Response) IsSuccessStatusCode() bool {
	return 200 <= r.StatusCode && r.StatusCode < 300
}

func (requestData RequestData) URLWithQuery() (*url.URL, error) {
	uri, err := url.Parse(requestData.URL)
	if err != nil {
		return nil, err
	}

	q := uri.Query()
	for k, v := range requestData.QueryParams {
		q.Set(k, v)
	}
	uri.RawQuery = q.Encode()
	return uri, nil
}
