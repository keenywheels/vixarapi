package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/keenywheels/backend/pkg/httpclient"
)

// endpoints LLM service endpoints
type endpoints struct {
	SentimentAnalysis string
}

// Client LLM client
type Client struct {
	client    *http.Client
	endpoints endpoints
	cfg       *Config
}

// NewClient create new LLM client
func NewClient(cfg *Config) *Client {
	client := httpclient.DefaultClient(cfg.Timeout)

	endpoints := endpoints{
		SentimentAnalysis: "/v1/sentiment", // TODO: согласовать
	}

	return &Client{
		client:    client,
		cfg:       cfg,
		endpoints: endpoints,
	}
}

// makeRequestJSON make HTTP request with JSON body and return response body
func (c *Client) makeRequestJSON(ctx context.Context, method string, path string, data interface{}) ([]byte, error) {
	if data == nil {
		return nil, errors.New("empty data for request")
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("failed to encode request body")
	}

	req, err := http.NewRequestWithContext(ctx, method, c.cfg.URL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.ContentLength <= 0 {
		return nil, fmt.Errorf("got empty response from %s: code=%d", path, resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("got not ok status code from %s: code=%d body=%s",
			path, resp.StatusCode, string(bytes))
	}

	return bytes, nil
}
