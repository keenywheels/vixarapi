package vk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/keenywheels/backend/pkg/httpclient"
)

var (
	ErrGotEmptyResponse      = errors.New("got empty response")
	ErrGotNoOKStatusResponse = errors.New("got not ok status response")
)

// endpoints VK API endpoints
type endpoints struct {
	getTokenByGrantType string
	getUserInfo         string
}

// Client represents a VK API client
type Client struct {
	client    *http.Client
	endpoints endpoints
	cfg       *Config
}

// New creates a new VK API client
func New(cfg *Config) *Client {
	client := httpclient.DefaultClient(cfg.HTTP.Timeout)

	// trim trailing slash from base URL if exists
	if cfg.Auth.BaseURL[len(cfg.Auth.BaseURL)-1] == '/' {
		cfg.Auth.BaseURL = cfg.Auth.BaseURL[:len(cfg.Auth.BaseURL)-1]
	}

	endpoints := endpoints{
		getTokenByGrantType: fmt.Sprintf("%s%s", cfg.Auth.BaseURL, "/oauth2/auth"),
		getUserInfo:         fmt.Sprintf("%s%s", cfg.Auth.BaseURL, "/oauth2/user_info"),
	}

	return &Client{
		client:    client,
		endpoints: endpoints,
		cfg:       cfg,
	}
}

// makeRequest make HTTP request
func (c *Client) makeRequest(
	ctx context.Context,
	method string,
	url string,
	headers map[string]string,
	data []byte,
) ([]byte, error) {
	// create request
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	// parse response
	if resp.ContentLength <= 0 {
		return nil, fmt.Errorf("request to %s failed, code=%d: %w",
			url, resp.StatusCode, ErrGotEmptyResponse)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("request to %s failed, code=%d body=%s: %w",
			url, resp.StatusCode, string(bytes), ErrGotNoOKStatusResponse)
	}

	return bytes, nil
}
