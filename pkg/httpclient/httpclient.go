package httpclient

import (
	"context"
	"net"
	"net/http"
	"time"
)

// DefaultClient return default HTTP client
func DefaultClient(timeout time.Duration) *http.Client {
	transport := &http.Transport{
		DialContext: defaultTransportDialContext(&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}),
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}

// defaultTransportDialContext return default DialContext function
func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
