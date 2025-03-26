package cmd

import (
	"net/http"
)

type HttpClient struct {
	url string
}

func NewHttpClient(url string) *HttpClient {
	return &HttpClient{
		url: url,
	}
}

func (h *HttpClient) Get() int {

	client := http.Client{}
	req, e := http.NewRequest("GET", h.url, nil)

	if e != nil {
		return 500
	}

	resp, e := client.Do(req)

	if e != nil {
		return 500
	}

	defer resp.Body.Close()

	return resp.StatusCode
}
