package network

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
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
	request, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept-Encoding", "gzip")

	resp, err := r.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return r.readBody(resp)
}

func (r *request) SendPost(uri string, content interface{}) (*Response, error) {
	data, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)

	request, err := http.NewRequest("POST", uri, buf)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept-Encoding", "gzip")

	resp, err := r.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return r.readBody(resp)
}

func (r *request) readBody(resp *http.Response) (*Response, error) {
	defer resp.Body.Close()

	var reader io.ReadCloser
	var err error
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	respContent, err := ioutil.ReadAll(reader)
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
