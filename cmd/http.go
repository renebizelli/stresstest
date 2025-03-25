package cmd

import (
	"fmt"
)

type HttpClient struct {
	url string
}

func NewHttpClient(url string) *HttpClient {
	return &HttpClient{
		url: url,
	}
}

func (h *HttpClient) Get() {
	fmt.Printf("\n %s", h.url)
}
