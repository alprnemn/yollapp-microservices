package http

import (
	"github.com/alprnemn/yollapp-microservices/services/apigateway/internal/config"
	"net/http"
)

// Client is a wrapper around http.Client used by the API Gateway
// to communicate with downstream microservices.
type Client struct {
	HTTPClient *http.Client
}

// NewClient creates and returns a new Client instance based on the
// provided configuration.
//
// The function configures a custom http.Transport which controls
// connection pooling, timeouts, and other low-level networking behavior.
func NewClient(cfg config.ClientConfig) *Client {

	// Transport defines how HTTP connections are created and managed.
	// It controls connection reuse, pooling, timeouts, and protocol details.
	transport := &http.Transport{
		MaxIdleConns:          cfg.TransportConfig.MaxIdleConns,
		MaxIdleConnsPerHost:   cfg.TransportConfig.MaxIdleConnsPerHost,
		MaxConnsPerHost:       cfg.TransportConfig.MaxConnsPerHost,
		IdleConnTimeout:       cfg.TransportConfig.IdleConnTimeout,
		DisableCompression:    cfg.TransportConfig.DisableCompression,
		DisableKeepAlives:     cfg.TransportConfig.DisableKeepAlives,
		TLSHandshakeTimeout:   cfg.TransportConfig.TLSHandshakeTimeout,
		ExpectContinueTimeout: cfg.TransportConfig.ExpectContinueTimeout,
	}

	httpClient := &http.Client{
		Timeout:   cfg.Timeout,
		Transport: transport,
	}

	return &Client{
		HTTPClient: httpClient,
	}
}
