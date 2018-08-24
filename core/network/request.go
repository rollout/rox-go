package network

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Request interface {
	SendGet(requestData RequestData) (response *Response, err error)
	SendPost(uri string, content interface{}) (response *Response, err error)
}

type request struct {
	httpClient *http.Client
}

type RequestData struct {
	URL         string
	QueryParams map[string]string
}

type Response struct {
	StatusCode int
	Content    []byte
}

func NewRequest(httpClient *http.Client) Request {
	return &request{httpClient: httpClient}
}

func (r *request) SendGet(requestData RequestData) (*Response, error) {
	uri, err := requestData.URLWithQuery()
	if err != nil {
		return nil, err
	}
	resp, err := r.httpClient.Get(uri.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respContent, err := ioutil.ReadAll(resp.Body)
	return &Response{resp.StatusCode, respContent}, err
}

func (r *request) SendPost(uri string, content interface{}) (*Response, error) {
	data, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)

	resp, err := r.httpClient.Post(uri, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respContent, err := ioutil.ReadAll(resp.Body)
	return &Response{resp.StatusCode, respContent}, err
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

func (r Response) IsSuccessStatusCode() bool {
	return 200 <= r.StatusCode && r.StatusCode < 300
}
