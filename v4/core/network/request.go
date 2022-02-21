package network

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rollout/rox-go/v4/core/model"
)

type request struct {
	httpClient *http.Client
}

func NewRequest(httpClient *http.Client) model.Request {
	return &request{httpClient: httpClient}
}

func (r *request) SendGet(requestData model.RequestData) (*model.Response, error) {
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

func (r *request) SendPost(uri string, content interface{}) (*model.Response, error) {
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
	request.Header.Add("Content-Type", "application/json")

	resp, err := r.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return r.readBody(resp)
}

func (r *request) readBody(resp *http.Response) (*model.Response, error) {
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
	return &model.Response{resp.StatusCode, respContent}, err
}
