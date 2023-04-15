package internal

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type IHttpClient interface {
	Do(req *http.Request) (*http.Response, error)
	SetNoTLSVerify()
	SetTimeout(timeout time.Duration)
}

type HttpClientConfig struct {
	Logger zerolog.Logger
}

func NewHttpClient(conf HttpClientConfig) *HttpClient {
	return &HttpClient{
		log: conf.Logger,
		client: &http.Client{},
	}
}

var _ IHttpClient = (*HttpClient)(nil)

type HttpClient struct {
	log    zerolog.Logger
	client *http.Client
}

func (h *HttpClient) Do(req *http.Request) (*http.Response, error) {
	return h.client.Do(req)
}

func (h *HttpClient) SetNoTLSVerify() {
	h.client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func (h *HttpClient) SetTimeout(timeout time.Duration) {
	h.client.Timeout = timeout
}
